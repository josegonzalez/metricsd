# go-metricsd
a metrics collecting agent written in go

# installation

You need the following go packages

- github.com/Sirupsen/logrus
- github.com/fzzy/radix
- github.com/c9s/goprocinfo

and then run `go run *.go` in this directory

# configuration

- `ELASTICSEARCH_URL`
- `REDIS_HOST`
- `REDIS_PORT`
