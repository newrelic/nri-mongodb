# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
