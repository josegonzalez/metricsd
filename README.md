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

create an ini file with the following contents:

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
```
