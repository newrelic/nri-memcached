# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#render-markdown-and-update-markdown)

## Unreleased

## v2.7.1 - 2025-06-30

### ⛓️ Dependencies
- Updated golang version to v1.24.4

## v2.7.0 - 2025-03-10

### 🚀 Enhancements
- Add FIPS compliant packages

### ⛓️ Dependencies
- Updated golang patch version to v1.23.6

## v2.6.2 - 2025-01-20

### ⛓️ Dependencies
- Updated golang patch version to v1.23.5

## v2.6.1 - 2024-12-09

### ⛓️ Dependencies
- Updated golang patch version to v1.23.4

## v2.6.0 - 2024-10-14

### dependency
- Upgrade go to 1.23.2

### 🚀 Enhancements
- Upgrade integrations SDK so the interval is variable and allows intervals up to 5 minutes

## v2.5.8 - 2024-09-16

### ⛓️ Dependencies
- Updated golang version to v1.23.1

## v2.5.7 - 2024-07-08

### ⛓️ Dependencies
- Updated golang version to v1.22.5

## v2.5.6 - 2024-05-13

### ⛓️ Dependencies
- Updated golang version to v1.22.3

## v2.5.5 - 2024-04-15

### ⛓️ Dependencies
- Updated golang version to v1.22.2

## v2.5.4 - 2024-03-11

### 🐞 Bug fixes
- Updated golang to version v1.21.7 to fix a vulnerability

## v2.5.3 - 2024-02-26

### ⛓️ Dependencies
- Updated github.com/newrelic/infra-integrations-sdk to v3.8.2+incompatible

## v2.5.2 - 2024-02-12

### ⛓️ Dependencies
- Updated github.com/newrelic/infra-integrations-sdk to v3.8.0+incompatible

## v2.5.1 - 2023-10-30

### ⛓️ Dependencies
- Updated golang version to 1.21

## 2.5.0 (2023-06-06)
### Changed
- Upgrade Go version to 1.20

## 2.4.1  (2022-06-28)

### Changed
- Bump dependencies
### Added
Added support for more distributions:
- RHEL(EL) 9
- Ubuntu 22.04

## 2.4.0 (2022-05-09)
### Fixed
- Use `prate` for metrics that report `*.PerSecond` stats. This prevents that metrics have negative values which is unexpected for this kind of metric.
- Use `prate` for metrics that are reported as counters (accumulators) by the service. This prevent the metric has negative values if the counter resets.
### Changed
- bump dependencies:
    `github.com/mitchellh/mapstructure v1.5.0`
	`github.com/newrelic/infra-integrations-sdk v3.7.2+incompatible`
	`github.com/stretchr/testify v1.7.1`
- change pipeline to compile with Go 1.18

## 2.3.2 (2021-10-20)
### Added
Added support for more distributions:
- Debian 11
- Ubuntu 20.10
- Ubuntu 21.04
- SUSE 12.15
- SUSE 15.1
- SUSE 15.2
- SUSE 15.3
- Oracle Linux 7
- Oracle Linux 8

## 2.3.1 (2021-08-27)
### Added
Moved default config.sample to [V4](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-newer-configuration-format/), added a dependency for infra-agent version 1.20.0

Please notice that old [V3](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) configuration format is deprecated, but still supported.

## 2.2.1 (2021-06-10)
### Changed
- Update Go to v1.16.

## 2.2.0 (2021-05-10)
### Changed
- Update Go to v1.16.
- Migrate to Go Modules
- Update Infrastracture SDK to v3.6.7.
- Update other dependecies.

## 2.1.3 (2021-03-24)
### Changed
- Added arm packages and binaries

## 2.1.2 (2020-07-15)
### Fixed
- Issue with calculating deltas for slabs because of entity uniqueness

## 2.1.0 (2019-11-18)
### Changed
- Renamed the integration executable from nr-memcached to nri-memcached in order to be consistent with the package naming. **Important Note:** if you have any security module rules (eg. SELinux), alerts or automation that depends on the name of this binary, these will have to be updated.

## 2.0.1 - 2018-10-23
### Added
- Add windows install packaging

## 2.0.0 - 2018-05-06
### Changed
- Update SDK
- Make entity keys more unique

## 1.0.1 - 2018-02-25
### Fixed
- Added prefix for all_data

## 1.0.0 - 2018-02-05
### Changed
- Bumped version to 1.0.0

## 0.1.1 - 2018-11-15
### Added
- Added metadata with hostname for easier filtering

## 0.1.0 - 2018-11-08
### Added
- Initial version: Includes Metrics and Inventory data
