//go:generate goversioninfo
package main

import (
	"fmt"

	"github.com/memcachier/mc"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

const (
	integrationName    = "com.newrelic.memcached"
	integrationVersion = "2.1.0"
)

var (
	args argumentList
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Host     string `default:"localhost" help:"The memcached host."`
	Port     string `default:"11211"     help:"The memcached port."`
	Username string `default:""          help:"The memcached SASL username."`
	Password string `default:""          help:"The memcached SASL password."`
}

func main() {
	memcachedIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	if err != nil {
		log.Error("Failed to create integration: %s", err)
		return
	}

	log.SetupLogging(args.Verbose)

	client := mc.NewMC(fmt.Sprintf("%s:%s", args.Host, args.Port), args.Username, args.Password)

	if args.HasMetrics() {
		CollectGeneralStats(client, memcachedIntegration)
		CollectSlabStats(client, memcachedIntegration)
		CollectItemStats(client, memcachedIntegration)
	}

	if args.HasInventory() {
		CollectSettings(client, memcachedIntegration)
	}

	if err = memcachedIntegration.Publish(); err != nil {
		log.Error("Failed to publish integration: %s", err.Error())
	}
}
