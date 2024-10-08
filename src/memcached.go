//go:generate goversioninfo
package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/memcachier/mc"
	sdkArgs "github.com/newrelic/infra-integrations-sdk/v3/args"
	"github.com/newrelic/infra-integrations-sdk/v3/integration"
	"github.com/newrelic/infra-integrations-sdk/v3/log"
)

const (
	integrationName = "com.newrelic.memcached"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Host        string `default:"localhost" help:"The memcached host."`
	Port        string `default:"11211"     help:"The memcached port."`
	Username    string `default:""          help:"The memcached SASL username."`
	Password    string `default:""          help:"The memcached SASL password."`
	ShowVersion bool   `default:"false" help:"Print build information and exit"`
}

var (
	args               argumentList
	integrationVersion = "0.0.0"
	gitCommit          = ""
	buildDate          = ""
)

func main() {
	memcachedIntegration, err := integration.New(integrationName, integrationVersion, integration.Args(&args))
	if err != nil {
		log.Error("Failed to create integration: %s", err)
		return
	}

	if args.ShowVersion {
		fmt.Printf(
			"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
			strings.Title(strings.Replace(integrationName, "com.newrelic.", "", 1)),
			integrationVersion,
			fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			runtime.Version(),
			gitCommit,
			buildDate)
		os.Exit(0)
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
