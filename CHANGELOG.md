# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#render-markdown-and-update-markdown)

## Unreleased

## v2.8.2 - 2024-02-12

### ⛓️ Dependencies
- Updated github.com/newrelic/infra-integrations-sdk to v3.8.0+incompatible

## v2.8.1 - 2023-10-30

### ⛓️ Dependencies
- Updated golang version to 1.21

## 2.8.0 (2023-06-06)
### Changed
- Update Go version to 1.20

## 2.7.2  (2022-06-28)
### Changed
- Bump dependencies
### Added
Added support for more distributions:
- RHEL(EL) 9
- Ubuntu 22.04

## 2.7.1 (2021-10-20)
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

## 2.7.0 (2021-09-30)
### Changed
- Moved default config.sample to [V4](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-newer-configuration-format/), added a dependency for infra-agent version 1.20.0

Please notice that old [V3](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) configuration format is deprecated, but still supported.

## 2.6.1 (2021-06-09)
### Changed
- Added support for ARM

## 2.6.0 (2021-05-05)
### Changed
- Decorate metrics with `mongodbClusterName`
  `clusterName` (deprecated) attribute might get overwritten by the Agent when running inside kubernetes.
  This change allows Mongodb entities to be uniquely identified when multiple Mongodb clusters are located in the same Kubernetes cluster.
- Removing duplicated win-config file
- Use go 1.16 and go mod
- Adds json schema integration integration tests  
## 2.5.1 (2021-02-09)
### Changed
- `username` and `password` arguments are now optional
  + This allows connecting to passwordless mongoDB instances

## 2.5.0 (2020-02-24)
### Added
- `concurrent_collections` argument to limit resource usage

## 2.4.0 (2019-11-18)
### Changed
- Renamed the integration executable from nr-mongodb to nri-mongodb in order to be consistent with the package naming. **Important Note:** if you have any security module rules (eg. SELinux), alerts or automation that depends on the name of this binary, these will have to be updated.

## 2.3.3 - 2019-11-15
### Fixed
- Remove windows definition file from linux packages

## 2.3.2 - 2019-11-13
### Added
- Windows MSI packaging

## 2.3.1 - 2019-10-28
### Fixed
- Using SSL certificates with unencrypted private key now works

## 2.3.0 - 2019-10-07
### Added
- A number of serverStatus wiredTiger cache statistics
- queryexecutor.scannedObjects

## 2.2.1 - 2019-09-13
### Fixed
- Fix a bug where replication lag was not calculated when the oplog was in timestamp-only format

## 2.2.0 - 2019-07-18
### Changed
- Name mongods by replica set config instead of shard config
### Added
- Support for non-sharded replica set deployments
- deploymentType inventory item (sharded_cluster, replica_set, or standalone)

## 2.1.0 - 2019-06-30
### Added
- PEM key support

## 2.0.0 - 2019-04-26
### Changed
- Updated the SDK
- Made entity keys more unique

## 1.1.3 - 2019-02-25
### Fixed
- Fix prefix for inventory on all_data

## 1.1.2 - 2019-02-12
### Fixed
- Don't panic on invalid mongo top response

## 1.1.1 - 2019-02-04
### Fixed
- Use correct protocol version

## 1.1.0 - 2019-01-09
### Fixed
- Detect if the monitored instance is a standalone instance, and collect the proper subset of metrics

## 1.0.0 - 2018-11-29
### Changes
- Bumped version for GA release

## 0.2.0 - 2018-11-12
### Changed
- Renamed configuration prefix for consistency

## 0.1.1 - 2018-10-25
### Fixed
- Added port to mongod entity name to provide unique identifier for multiple mongods on the same host

## 0.1.0 - 2018-09-25
### Added
- Initial version: Includes Metrics and Inventory data
