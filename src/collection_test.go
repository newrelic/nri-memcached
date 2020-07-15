package main

import (
	"testing"

	"github.com/memcachier/mc"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mockable client
type MockClient struct {
	mock.Mock
}

// StatsWithKey is a mocked version of the mc.Client object
func (m MockClient) StatsWithKey(key string) (map[string]mc.McStats, error) {
	args := m.Called(key)
	return args.Get(0).(map[string]mc.McStats), args.Error(1)
}

func Test_processGeneralStats(t *testing.T) {
	stats := map[string]string{
		"pid":            "28442",
		"bytes":          "999",
		"uptime":         "1235",
		"time":           "1235",
		"version":        "1.4.39",
		"libevent":       "2.0.21-stable",
		"get_hits":       "44",
		"get_misses":     "123",
		"incr_misses":    "1",
		"curr_items":     "333",
		"limit_maxbytes": "9990",
		"cmd_get":        "100",
	}

	i, _ := integration.New("test", "test")
	e, _ := i.Entity("test", "test")

	processGeneralStats(stats, e)

	assert.Equal(t, 1, len(e.Metrics))
	assert.Equal(t, float64(999), e.Metrics[0].Metrics["bytesUsedServerInBytes"])
	assert.Equal(t, float64(1235000), e.Metrics[0].Metrics["uptimeInMilliseconds"])
	assert.Equal(t, float64(3.0), e.Metrics[0].Metrics["avgItemSizeInBytes"])
	assert.Equal(t, float64(10), e.Metrics[0].Metrics["storingItemsPercentMemory"])
	assert.Equal(t, float64(44), e.Metrics[0].Metrics["getHitPercent"])
}

func Test_processItemStats(t *testing.T) {
	stats := map[string]string{
		"items:1:evicted_time": "1234",
		"items:2:evicted":      "1235",
	}

	i, _ := integration.New("test", "test")

	processItemStats(stats, i, "testhost")

	e1, _ := i.Entity("1", "mc-slab", integration.NewIDAttribute("mc-slab", "1"), integration.NewIDAttribute("host", "testhost"))
	e2, _ := i.Entity("2", "mc-slab", integration.NewIDAttribute("mc-slab", "2"), integration.NewIDAttribute("host", "testhost"))

	assert.Equal(t, 1, len(e1.Metrics))
	assert.Equal(t, 1, len(e2.Metrics))
	assert.Equal(t, float64(1234), e1.Metrics[0].Metrics["itemsTimeSinceEvictionInMilliseconds"])
	assert.Equal(t, float64(0), e2.Metrics[0].Metrics["evictionsBeforeExpirationPerSecond"])
}

func Test_processSlabStats(t *testing.T) {
	stats := map[string]string{
		"active_slabs":    "5",
		"1:chunk_size":    "4",
		"1:mem_requested": "4",
	}

	i, _ := integration.New("test", "test")

	id1 := integration.NewIDAttribute("mc-slab", "1")
	id2 := integration.NewIDAttribute("host", "testhost")
	e1, _ := i.Entity("1", "mc-slab", id1, id2)

	e2, _ := i.Entity("testHost", "mc-instance")

	processSlabStats(stats, i, "testHost")

	assert.Equal(t, 1, len(e1.Metrics))
	assert.Equal(t, float64(4), e1.Metrics[0].Metrics["chunkSizeInBytes"])
	assert.Equal(t, float64(0), e1.Metrics[0].Metrics["memRequestedSlabInBytesPerSecond"])
	assert.Equal(t, float64(5), e2.Metrics[0].Metrics["activeSlabs"])
}

func Test_processSettings(t *testing.T) {
	settings := map[string]string{
		"test1": "val1",
	}

	i, _ := integration.New("test", "test")
	processSettings(settings, i, "testHost")
	e, _ := i.Entity("testHost", "mc-instance")

	assert.Equal(t, "val1", e.Inventory.Items()["test1"]["value"])

}

func Test_CollectGeneralStats(t *testing.T) {
	client := MockClient{}
	stats := map[string]mc.McStats{
		"testHost": {
			"pid":            "28442",
			"bytes":          "999",
			"uptime":         "1235",
			"time":           "1235",
			"version":        "1.4.39",
			"libevent":       "2.0.21-stable",
			"get_hits":       "44",
			"get_misses":     "123",
			"incr_misses":    "1",
			"curr_items":     "333",
			"limit_maxbytes": "9990",
			"cmd_get":        "100",
		},
	}
	client.On("StatsWithKey", "").Return(stats, nil)

	i, _ := integration.New("test", "test")
	e, _ := i.Entity("testHost", "mc-instance")

	CollectGeneralStats(client, i)

	assert.Equal(t, 1, len(e.Metrics))
	assert.Equal(t, float64(999), e.Metrics[0].Metrics["bytesUsedServerInBytes"])
	assert.Equal(t, float64(1235000), e.Metrics[0].Metrics["uptimeInMilliseconds"])
	assert.Equal(t, float64(3.0), e.Metrics[0].Metrics["avgItemSizeInBytes"])
	assert.Equal(t, float64(10), e.Metrics[0].Metrics["storingItemsPercentMemory"])
	assert.Equal(t, float64(44), e.Metrics[0].Metrics["getHitPercent"])

}

func Test_CollectSlabStats(t *testing.T) {
	client := MockClient{}
	stats := map[string]mc.McStats{
		"testHost": {
			"active_slabs":    "5",
			"1:chunk_size":    "4",
			"1:mem_requested": "4",
		},
	}
	client.On("StatsWithKey", "slabs").Return(stats, nil)

	i, _ := integration.New("test", "test")
	idattr := integration.NewIDAttribute("mc-slab", "1")
	id2 := integration.NewIDAttribute("host", "testHost")
	e1, _ := i.Entity("1", "mc-slab", idattr, id2)

	CollectSlabStats(client, i)

	assert.Equal(t, 1, len(e1.Metrics))
	assert.Equal(t, float64(4), e1.Metrics[0].Metrics["chunkSizeInBytes"])
	assert.Equal(t, float64(0), e1.Metrics[0].Metrics["memRequestedSlabInBytesPerSecond"])
}

func Test_CollectItemStats(t *testing.T) {
	client := MockClient{}
	stats := map[string]mc.McStats{
		"testHost": {
			"items:1:evicted_time": "1234",
			"items:2:evicted":      "1235",
		},
	}
	client.On("StatsWithKey", "items").Return(stats, nil)

	i, _ := integration.New("test", "test")

	CollectItemStats(client, i)

	id1 := integration.NewIDAttribute("mc-slab", "1")
	idHost := integration.NewIDAttribute("host", "testHost")
	e1, _ := i.Entity("1", "mc-slab", id1, idHost)
	id2 := integration.NewIDAttribute("mc-slab", "2")
	e2, _ := i.Entity("2", "mc-slab", id2, idHost)

	assert.Equal(t, 1, len(e1.Metrics))
	assert.Equal(t, 1, len(e2.Metrics))
	assert.Equal(t, float64(1234), e1.Metrics[0].Metrics["itemsTimeSinceEvictionInMilliseconds"])
	assert.Equal(t, float64(0), e2.Metrics[0].Metrics["evictionsBeforeExpirationPerSecond"])
}

func Test_CollectSettings(t *testing.T) {
	client := MockClient{}
	stats := map[string]mc.McStats{
		"testHost": {
			"test1": "val1",
		},
	}
	client.On("StatsWithKey", "settings").Return(stats, nil)

	i, _ := integration.New("test", "test")

	CollectSettings(client, i)

	e, _ := i.Entity("testHost", "mc-instance")

	assert.Equal(t, "val1", e.Inventory.Items()["test1"]["value"])
}
