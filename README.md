# New Relic Infrastructure Integration for MongoDB 

Reports status and metrics for MongoDB service

See our [documentation web site](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/mongodb-monitoring-integration) for more details.

## Requirements

None

## Installation

* Download an archive file for the `MongoDB` Integration
* Extract `mongodb-definition.yml` and the `bin` directory into `/var/db/newrelic-infra/newrelic-integrations`
* Add execute permissions for the binary file `nr-mongodb` (if required)
* Extract `mongodb-config.yml.sample` into `/etc/newrelic-infra/integrations.d`

## Usage

To run the MongoDB integration, you must have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

To use the MongoDB integration, first rename `mongodb-config.yml.sample` to `mongodb-config.yml`, then configure the integration
by editing the fields in the file. 

You can view your data in Insights by creating your own NRQL queries. To do so, use the **MongoSample**, **MongodSample**, **MongosSample**, **MongoConfigServerSample**, and **MongoCollectionSample** event types. 

## Compatibility

* Supported OS: No limitations
* MongoDB versions: 3.0+

## Integration Development usage

Assuming you have the source code, you can build and run the MongoDB integration locally

* Go to the directory of the MongoDB Integration and build it
```
$ make
```

* The command above will execute tests for the MongoDB integration and build an executable file called `nr-mongodb` in the `bin` directory
```
$ ./bin/nr-mongodb --help
```

For managing external dependencies, the [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to a specific version (if possible) in the vendor directory.
