# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 2.2.1 (2021-06-10)
## Changed
- Update Go to v1.16.

## 2.2.0 (2021-05-10)
## Changed
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
