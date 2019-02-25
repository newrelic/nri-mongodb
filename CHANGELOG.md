# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
