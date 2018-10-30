package main

import (
	//sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	//"github.com/newrelic/infra-integrations-sdk/data/event"
	//"github.com/newrelic/infra-integrations-sdk/data/metric"
	"bufio"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/newrelic/infra-integrations-sdk/data/metric"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
	"net"
	"regexp"
)

const (
	integrationName    = "com.newrelic.memcached"
	integrationVersion = "0.1.0"
)

var (
	args arguments.ArgumentList
)

func main() {

	memcachedIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	if err != nil {
		log.Error("Failed to create integration")
	}

	generalStats := getStats("")
	processGeneralStats(generalStats, memcachedIntegration)

	items := getStats("items")
	processItemStats(items, memcachedIntegration)

	slabs := getStats("slabs")
	processSlabStats(slabs, memcachedIntegration)

	settings := getStats("settings")
	processSettings(settings, memcachedIntegration)

	memcachedIntegration.Publish()
}

func processGeneralStats(stats map[string]string, i *integration.Integration) {
	var s GeneralStats
	config := mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &s,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		panic("unreachable")
	}
	decoder.Decode(stats)

	// Calculate post-processed metrics
	if s.Bytes == nil || s.CurrItems == nil {
		return
	}
	averageSize := float64(*s.Bytes) / float64(*s.CurrItems)
	s.AverageItemSize = &averageSize

	if s.LimitMaxbytes == nil {
		return
	}
	percentMaxUsed := float64(*s.Bytes) / float64(*s.LimitMaxbytes) * 100.0
	s.PercentMaxUsed = &percentMaxUsed

	if s.CmdGet == nil || s.GetHits == nil {
		return
	}
	getHitPercent := float64(*s.GetHits) / float64(*s.CmdGet) * 100.0
	s.GetHitPercent = &getHitPercent

	if s.Uptime == nil {
		return
	}
	uptimeMilliseconds := *s.Uptime * 1000
	s.UptimeMilliseconds = &uptimeMilliseconds

	e, _ := i.Entity("test", "instance")
	ms := e.NewMetricSet("MemcachedSample",
		metric.Attribute{Key: "displayName", Value: "test"},
		metric.Attribute{Key: "entityName", Value: "test"}, // TODO change the naming
	)
	err = ms.MarshalMetrics(s)
	if err != nil {
		println(err.Error())
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
			panic("unreachable")
		}
		decoder.Decode(slabMetrics)

		e, _ := i.Entity("slab"+slabID, "slab")
		ms := e.NewMetricSet("MemcachedSlabSample",
			metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
			metric.Attribute{Key: "entityName", Value: "slab:" + e.Metadata.Name},
		)
		err = ms.MarshalMetrics(s)
		if err != nil {
			println(err.Error())
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

func processSlabStats(stats map[string]string, i *integration.Integration) {
	// TODO do something with general statistics
	slabs, _ := partitionSlabsBySlabID(stats)

	for slabID, slabStats := range slabs {
		var s SlabStats
		config := mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			Result:           &s,
		}

		decoder, err := mapstructure.NewDecoder(&config)
		if err != nil {
			panic("unreachable")
		}
		decoder.Decode(slabStats)

		e, _ := i.Entity("slab"+slabID, "slab")
		ms := e.NewMetricSet("MemcachedSlabSample",
			metric.Attribute{Key: "displayName", Value: e.Metadata.Name},
			metric.Attribute{Key: "entityName", Value: "slab:" + e.Metadata.Name},
		)
		err = ms.MarshalMetrics(s)
		if err != nil {
			println(err.Error())
		}

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

func processSettings(settings map[string]string, i *integration.Integration) {
	e, _ := i.Entity("test", "test") // TODO rename
	for key, value := range settings {
		if err := e.SetInventoryItem(key, "value", value); err != nil {
			log.Error("Failed to set inventory item for key %s: %s", key, err.Error())
		}

	}
}

// TODO authentication
func getStats(key string) map[string]string {

	conn, err := net.Dial("tcp", "mc-14-1:11211")
	if err != nil {
		println(err.Error())
	}

	_, err = conn.Write([]byte(fmt.Sprintf("stats %s\n", key)))
	if err != nil {
		println(err.Error())
	}

	response := bufio.NewReader(conn)
	endRegex := regexp.MustCompile(`^END`)
	statRegex := regexp.MustCompile(`^STAT ([:a-z0-9_]+) ([^\s]+)`)

	returnMap := make(map[string]string)

	for {
		line, err := response.ReadBytes('\n')
		if err != nil {
			break
		}

		if endRegex.Match(line) {
			break
		}

		matches := statRegex.FindSubmatch(line)
		if len(matches) != 3 {
			println("Bad stats form")
		}

		returnMap[string(matches[1])] = string(matches[2])

	}

	return returnMap
}
