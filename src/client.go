package main

import (
	"github.com/memcachier/mc"
)

// Client is an interface that represents a memcached client
type Client interface {
	StatsWithKey(string) (map[string]mc.McStats, error)
}
