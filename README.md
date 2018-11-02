# New Relic Infrastructure Integration for Memcached 

Reports status and metrics for Memcached service

## Requirements

None

## Installation

* Download an archive file for the `Memcached` Integration
* Extract `memcached-definition.yml` and the `bin` directory into `/var/db/newrelic-infra/newrelic-integrations`
* Add execute permissions for the binary file `nr-memcached` (if required)
* Extract `memcached-config.yml.sample` into `/etc/newrelic-infra/integrations.d`

## Usage

To run the Memcached integration, you must have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

To use the Memcached integration, first rename `memcached-config.yml.sample` to `memcached-config.yml`, then configure the integration
by editing the fields in the file. 

You can view your data in Insights by creating your own NRQL queries. To do so, use the **MemcachedSample** and **MemcachedSlabSample** event types.

## Compatibility

* Supported OS: No limitations
* Memcached 1.4+ 

## Integration Development usage

Assuming you have the source code, you can build and run the Memcached integration locally

* Go to the directory of the Memcached Integration and build it
```
$ make
```

* The command above will execute tests for the Memcached integration and build an executable file called `nr-memcached` in the `bin` directory
```
$ ./bin/nr-memcached --help
```

For managing external dependencies, the [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to a specific version (if possible) in the vendor directory.
