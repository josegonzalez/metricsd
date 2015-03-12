# go-metricsd
a metrics collecting agent written in go

# installation

You need the following go packages

- github.com/Sirupsen/logrus
- github.com/fzzy/radix
- github.com/c9s/goprocinfo

and then run `go run *.go` in this directory

# configuration

- `ELASTICSEARCH_INDEX`: defaults to `metricsd-data`
- `ELASTICSEARCH_URL`: defaults to `http://127.0.0.1:9200`
- `METRIC_TYPE`: defaults to `metricsd`
- `REDIS_HOST`: defaults to `127.0.0.1`
- `REDIS_LIST`: defaults to `metricsd`
- `REDIS_PORT`: defaults to `6379`
