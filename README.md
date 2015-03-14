# go-metricsd
a metrics collecting agent written in go

# installation

You need the following go packages

- github.com/Sirupsen/logrus
- github.com/fzzy/radix
- github.com/c9s/goprocinfo
- github.com/vaughan0/go-ini
- github.com/josegonzalez/go-radixurl

and then run `go run *.go -config=path/to/config.ini` in this directory

# configuration

Collectors and Shippers are configured in an `ini` file. You *must* specify `enabled = true` under the stanza for that collector/shipper in order to enable it. Other configuration for the respective collector/shipper can also be place in those sections.

Below is a sample `config.ini`:

```
[ElasticsearchShipper]
enabled = true
index = metricsd-data
type = metricsd
url = http://127.0.0.1:9200

[RedisShipper]
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
