# metricsd

a metrics collecting agent written in go

# installation

You can setup metricsd via the following short script:

```shell
go get github.com/c9s/goprocinfo/linux
go get github.com/fzzy/radix
go get github.com/josegonzalez/go-radixurl
go get github.com/ogier/pflag
go get github.com/Sirupsen/logrus
go get github.com/vaughan0/go-ini

go get github.com/josegonzalez/metricsd
go build github.com/josegonzalez/metricsd
```

Then you can execute the script using an ini-file of your choice:

```shell
./metricsd --config="path/to/config.ini"
```

# configuration

Collectors and Shippers are configured in an `ini` file. You *must* specify `enabled = true` under the stanza for that collector/shipper in order to enable it. Other configuration for the respective collector/shipper can also be place in those sections.

Below is a sample `config.ini` that enables every collector and shipper:

```
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

[LoadAvgCollector]
enabled = true

[MemoryCollector]
enabled = true

[VmstatCollector]
enabled = true
```
