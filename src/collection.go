package main

import (
	"regexp"

	"github.com/memcachier/mc"
	"github.com/mitchellh/mapstructure"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

// CollectGeneralStats collects general stats from the client and populates them into the integration
func CollectGeneralStats(client *mc.Client, i *integration.Integration) {
	generalStats, err := client.StatsWithKey("")
	if err != nil {
		log.Error("Failed to retrieve general stats: %s", err.Error())
		return
	}
	// Usually only one
	for host, hostStats := range generalStats {
		e, err := i.Entity(host, "instance")
		if err != nil {
			log.Error("Failed to retrieve entity for instance %s: %s", host, err.Error())
		}

		processGeneralStats(hostStats, e)
	}
}

// CollectSlabStats collects slab stats from the client and populates them into the integration
func CollectSlabStats(client *mc.Client, i *integration.Integration) {
	slabStats, err := client.StatsWithKey("slabs")
	if err != nil {
		log.Error("Failed to retrieve slabs stats: %s", err.Error())
		return
	}
	// Usually only one
	for server, serverStats := range slabStats {
		processSlabStats(serverStats, i, server)
	}
}

// CollectItemStats collects item stats from the client and populates them into the integration
func CollectItemStats(client *mc.Client, i *integration.Integration) {
	itemStats, err := client.StatsWithKey("items")
	if err != nil {
		log.Error("Failed to retrieve items stats: %s", err.Error())
		return
	}
	// Usually only one
	for _, serverStats := range itemStats {
		processItemStats(serverStats, i)
	}
}

// CollectSettings collects the list of settings and from the client and populates them into the integration
func CollectSettings(client *mc.Client, i *integration.Integration) {
	settingsResult, err := client.StatsWithKey("setting")
	if err != nil {
		log.Error("Failed to retrieve settings: %s", err.Error())
		return
	}
	// Usually only one
	for host, settings := range settingsResult {
		processSettings(settings, i, host)
	}
}

func processGeneralStats(stats map[string]string, e *integration.Entity) {
	var s GeneralStats
	config := mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		log.Error("Failed to create map decoder: %s", err.Error())
		return
	}
	err = decoder.Decode(stats)
	if err != nil {
		log.Error("Failed to decode map: %s", err.Error())
	}

	calculateProcessedMetrics(&s)

	// Create metric set
	ms := e.NewMetricSet("MemcachedSample",
		metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: "instance:" + e.Metadata.Name},
	)

	err = ms.MarshalMetrics(s)
	if err != nil {
		log.Error("Failed to marshal general statistics: %s", err.Error())
	}
}

func calculateProcessedMetrics(s *GeneralStats) {
	// Calculate post-processed metrics
	if s.Bytes == nil || s.CurrItems == nil {
		log.Error("Failed to collect metrics needed to calculate averageItemSize")
	} else if *s.CurrItems != 0 {
		averageSize := float64(*s.Bytes) / float64(*s.CurrItems)
		s.AverageItemSize = &averageSize
	}

	if s.LimitMaxbytes == nil {
		log.Error("Failed to collect metrics needed to calculate percentMaxUsed")
	} else if *s.LimitMaxbytes != 0 {
		percentMaxUsed := float64(*s.Bytes) / float64(*s.LimitMaxbytes) * 100.0
		s.PercentMaxUsed = &percentMaxUsed
	}

	if s.CmdGet == nil || s.GetHits == nil {
		log.Error("Failed to collect metrics needed to calculate getHitPercent")
	} else if *s.CmdGet != 0 {
		getHitPercent := float64(*s.GetHits) / float64(*s.CmdGet) * 100.0
		s.GetHitPercent = &getHitPercent
	}

	if s.Uptime == nil {
		log.Error("Failed to collect metrics needed to calculate uptimeMilliseconds")
	} else {
		uptimeMilliseconds := *s.Uptime * 1000
		s.UptimeMilliseconds = &uptimeMilliseconds
	}
}

func processItemStats(stats map[string]string, i *integration.Integration) {
	slabs := partitionItemsBySlabID(stats)
	for slabID, slabMetrics := range slabs {
		var s ItemStats
		config := mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           &s,
		}

		decoder, err := mapstructure.NewDecoder(&config)
		if err != nil {
			log.Error("Failed to create map decoder: %s", err.Error())
			return
		}

		err = decoder.Decode(slabMetrics)
		if err != nil {
			log.Error("Failed to decode map: %s", err.Error())
		}

		e, _ := i.Entity(slabID, "slab")
		ms := e.NewMetricSet("MemcachedSlabSample",
			metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
			metric.Attribute{Key: "slabID", Value: e.Metadata.Name},
			metric.Attribute{Key: "entityName", Value: "slab:" + e.Metadata.Name},
		)
		err = ms.MarshalMetrics(s)
		if err != nil {
			log.Error("Failed to marshal item statistics: %s", err.Error())
		}

	}
}

func partitionItemsBySlabID(items map[string]string) map[string]map[string]string {
	pattern := regexp.MustCompile(`items:(\d+):([a-z_]+)`)

	returnMap := make(map[string]map[string]string)
	for key, val := range items {
		matches := pattern.FindStringSubmatch(key)
		slabID := matches[1]
		metricName := matches[2]

		// Retrieve the slab metrics. Create it if it doesn't exist
		slab, ok := returnMap[slabID]
		if !ok {
			slab = make(map[string]string)
			returnMap[slabID] = slab
		}

		slab[metricName] = val
	}

	return returnMap
}

func processSlabStats(stats map[string]string, i *integration.Integration, host string) {
	slabs, clusterStats := partitionSlabsBySlabID(stats)
	processClusterSlabStats(clusterStats, i, host)

	for slabID, slabStats := range slabs {
		var s SlabStats
		config := mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           &s,
		}

		decoder, err := mapstructure.NewDecoder(&config)
		if err != nil {
			log.Error("Failed to create map decoder: %s", err.Error())
			return
		}
		err = decoder.Decode(slabStats)
		if err != nil {
			log.Error("Failed to decode map: %s", err.Error())
		}

		e, _ := i.Entity(slabID, "slab")
		ms := e.NewMetricSet("MemcachedSlabSample",
			metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
			metric.Attribute{Key: "slabID", Value: e.Metadata.Name},
			metric.Attribute{Key: "entityName", Value: "slab:" + e.Metadata.Name},
		)
		err = ms.MarshalMetrics(s)
		if err != nil {
			log.Error("Failed to marshal slabs metrics: %s", err.Error())
		}

	}
}

func processClusterSlabStats(stats map[string]string, i *integration.Integration, host string) {
	var c ClusterSlabStats
	config := mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &c,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		log.Error("Failed to create map decoder: %s", err.Error())
		return
	}
	err = decoder.Decode(stats)
	if err != nil {
		log.Error("Failed to decode map: %s", err.Error())
	}

	instanceEntity, _ := i.Entity(host, "instance")
	ms := instanceEntity.NewMetricSet("MemcachedSample",
		metric.Attribute{Key: "displayName", Value: instanceEntity.Metadata.Name},
		metric.Attribute{Key: "entityName", Value: "instance:" + instanceEntity.Metadata.Name},
	)
	err = ms.MarshalMetrics(c)
	if err != nil {
		log.Error("Failed to marshal slabs metrics: %s", err.Error())
	}

}

func partitionSlabsBySlabID(slabs map[string]string) (map[string]map[string]string, map[string]string) {
	slabPattern := regexp.MustCompile(`(\d+):([a-z_]+)`)
	generalPattern := regexp.MustCompile(`^[a-z_]+$`)

	statsMap := make(map[string]map[string]string)
	generalMap := make(map[string]string)

	for key, val := range slabs {
		if generalPattern.MatchString(key) {
			generalMap[key] = val
			continue
		}

		matches := slabPattern.FindStringSubmatch(key)
		slabID := matches[1]
		metricName := matches[2]

		// Retrieve the slab metrics. Create it if it doesn't exist
		slab, ok := statsMap[slabID]
		if !ok {
			slab = make(map[string]string)
			statsMap[slabID] = slab
		}

		slab[metricName] = val
	}

	return statsMap, generalMap

}

func processSettings(settings map[string]string, i *integration.Integration, host string) {
	e, err := i.Entity(host, "instance")
	if err != nil {
		log.Error("Failed to get entity for host %s: %s", host, err.Error())
		return
	}
	for key, value := range settings {
		if err := e.SetInventoryItem(key, "value", value); err != nil {
			log.Error("Failed to set inventory item for key %s: %s", key, err.Error())
		}
	}
}
