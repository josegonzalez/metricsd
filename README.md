# metricsd [![Build Status](https://travis-ci.org/josegonzalez/metricsd.svg?branch=master)](https://travis-ci.org/josegonzalez/metricsd)

a metrics collecting agent for linux written in go

## why

I was looking for a simple replacement for diamond. It's codebase is a bit weird to navigate, and it has built up lots of cruft over the years.

I also wanted to learn a bit of golang, and this seemed like a perfect excuse to write my own metrics collector.

## features

- written in golang, so it is shippable as a single binary
- simple configuration via an ini file
- assumes linux as operating system
- supports [metrics 2.0](http://metrics20.org/) standard
- easy generation of new collectors

## installation

You can setup metricsd via the following short script:

```shell
make dependencies
go get github.com/josegonzalez/metricsd
go build github.com/josegonzalez/metricsd
```

Then you can execute the script using an ini-file of your choice:

```shell
./metricsd --config="path/to/config.ini"
```

## configuration

Most configuration is handled via an `ini` file. You can specify the ini file file via the `--config` flag when calling `metricsd`:

```shell
metricsd --config="path/to/config.ini"
```

By default, the ini file is set to the path `/etc/metricsd/metricsd.ini`.

You can also specify a debug level for metricsd via the `--loglevel` flag as follows:

```shell
metricsd --loglevel=debug
metricsd --loglevel=info
metricsd --loglevel=warning
metricsd --loglevel=error
metricsd --loglevel=fatal
metricsd --loglevel=panic
```

The default `loglevel` is `warning`.

### metricsd configuration

`metricsd` has a few configuration fields that can be set via the ini file:

```ini
[metricsd]
interval = 30
loop = false
```

- `interval`: Default `30`. Time in seconds to query for metrics.
- `loop`: Default `false`. If set to `true`, then `metricsd` will continue running, collecting metrics at the configured `interval`.

### collectors and shippers

Collectors and Shippers are configured in an `ini` file. You *must* specify `enabled = true` under the stanza for that collector/shipper in order to enable it. Other configuration for the respective collector/shipper can also be place in those sections.

Below is a sample `config.ini` that enables every collector and shipper:

```ini
[GraphiteShipper]
debug = true
enabled = true
url = tcp://127.0.0.1:2003
prefix = servers

[LogstashElasticsearchShipper]
enabled = true
index = metricsd-data
type = metricsd
url = http://127.0.0.1:9200

[LogstashRedisShipper]
enabled = true
url = redis://127.0.0.1:6379/0
list = metricsd

[StdoutShipper]
enabled = true

[CpuCollector]
enabled = true

[DiskspaceCollector]
enabled = true

[ElasticsearchCollector]
enabled = true

[LoadAvgCollector]
enabled = true

[MemoryCollector]
enabled = true

[RedisCollector]
enabled = true

[SocketsCollector]
enabled = true

[VmstatCollector]
enabled = true
```
