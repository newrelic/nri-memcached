package main

import (
	"regexp"

	"github.com/mitchellh/mapstructure"
	"github.com/newrelic/infra-integrations-sdk/v3/data/attribute"
	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
)

type statsSet map[string]string

// CollectGeneralStats collects general stats from the client and populates them into the integration
func CollectGeneralStats(client Client, i *integration.Integration) {
	generalStats, err := client.StatsWithKey("")
	if err != nil {
		log.Error("Failed to retrieve general stats: %s", err.Error())
		return
	}

	// We will only get one host back
	for host, hostStats := range generalStats {
		e, err := i.EntityReportedVia(host, host, "mc-instance")
		if err != nil {
			log.Error("Failed to retrieve entity for instance %s: %s", host, err.Error())
		}

		processGeneralStats(hostStats, e)
	}
}

// CollectSlabStats collects slab stats from the client and populates them into the integration
func CollectSlabStats(client Client, i *integration.Integration) {
	slabStats, err := client.StatsWithKey("slabs")
	if err != nil {
		log.Error("Failed to retrieve slabs stats: %s", err.Error())
		return
	}

	// We will only get one host back
	for host, serverStats := range slabStats {
		processSlabStats(serverStats, i, host)
	}
}

// CollectItemStats collects item stats from the client and populates them into the integration
func CollectItemStats(client Client, i *integration.Integration) {
	itemStats, err := client.StatsWithKey("items")
	if err != nil {
		log.Error("Failed to retrieve items stats: %s", err.Error())
		return
	}
	// Usually only one
	for host, serverStats := range itemStats {
		processItemStats(serverStats, i, host)
	}
}

// CollectSettings collects the list of settings and from the client and populates them into the integration
func CollectSettings(client Client, i *integration.Integration) {
	settingsResult, err := client.StatsWithKey("settings")
	if err != nil {
		log.Error("Failed to retrieve settings: %s", err.Error())
		return
	}

	// We will only get one host back
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
		attribute.Attribute{Key: "displayName", Value: e.Metadata.Name},
		attribute.Attribute{Key: "entityName", Value: "instance:" + e.Metadata.Name},
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

func processItemStats(stats map[string]string, i *integration.Integration, host string) {
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

		slabIDAttr := integration.NewIDAttribute("mc-slab", slabID)
		hostIDAttr := integration.NewIDAttribute("host", host)
		e, _ := i.EntityReportedVia(host, slabID, "mc-slab", slabIDAttr, hostIDAttr)
		ms := e.NewMetricSet("MemcachedSlabSample",
			attribute.Attribute{Key: "displayName", Value: e.Metadata.Name},
			attribute.Attribute{Key: "entityName", Value: "slab:" + e.Metadata.Name},
			attribute.Attribute{Key: "host", Value: host},
		)
		err = ms.MarshalMetrics(s)
		if err != nil {
			log.Error("Failed to marshal item statistics: %s", err.Error())
		}

	}
}

func partitionItemsBySlabID(items map[string]string) map[string]statsSet {
	pattern := regexp.MustCompile(`items:(\d+):([a-z_]+)`)

	partitionedMetrics := make(map[string]statsSet)
	for key, val := range items {
		matches := pattern.FindStringSubmatch(key)
		if len(matches) != 3 {
			log.Error("Failed to parse key %s", key)
			continue
		}

		slabID := matches[1]
		metricName := matches[2]

		// Retrieve the slab metrics. Create it if it doesn't exist
		slab, ok := partitionedMetrics[slabID]
		if !ok {
			slab = make(map[string]string)
			partitionedMetrics[slabID] = slab
		}

		slab[metricName] = val
	}

	return partitionedMetrics
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

		slabIDAttr := integration.NewIDAttribute("mc-slab", slabID)
		hostIDAttr := integration.NewIDAttribute("host", host)
		e, _ := i.EntityReportedVia(host, slabID, "mc-slab", slabIDAttr, hostIDAttr)
		ms := e.NewMetricSet("MemcachedSlabSample",
			attribute.Attribute{Key: "displayName", Value: e.Metadata.Name},
			attribute.Attribute{Key: "entityName", Value: "slab:" + e.Metadata.Name},
			attribute.Attribute{Key: "host", Value: host},
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

	instanceEntity, _ := i.EntityReportedVia(host, host, "mc-instance")
	ms := instanceEntity.NewMetricSet("MemcachedSample",
		attribute.Attribute{Key: "displayName", Value: instanceEntity.Metadata.Name},
		attribute.Attribute{Key: "entityName", Value: "instance:" + instanceEntity.Metadata.Name},
		attribute.Attribute{Key: "host", Value: host},
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
		if len(matches) != 3 {
			log.Error("Failed to parse key %s", key)
			continue
		}
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
	e, err := i.EntityReportedVia(host, host, "mc-instance")
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
