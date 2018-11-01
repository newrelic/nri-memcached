package main

import (
	"github.com/ccheek21/mc"
	"github.com/stretchr/testify/mock"
)

type Client interface {
	StatsWithKey(string) (map[string]mc.McStats, error)
}

type MockClient struct {
	mock.Mock
}

func (m MockClient) StatsWithKey(key string) (map[string]mc.McStats, error) {
	args := m.Called(key)
	return args.Get(0).(map[string]mc.McStats), args.Error(1)
}
