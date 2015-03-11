package main

import "encoding/json"
import "net/http"
import log "github.com/Sirupsen/logrus"

type ElasticsearchHook struct{}

func (hook *ElasticsearchHook) Fire(entry *log.Entry) error {
    data := MarshalData(entry)

    serialized, err := json.Marshal(data)
    if err != nil {
        log.Error("Failed to marshal fields to JSON, %v", err)
        return nil
    }

    status := ElasticsearchPost("/logstash-data/metricsd", serialized)
    if status != http.StatusCreated {
        log.Warning("Indexing serialized data failed, %s", err)
    }

    return nil
}

func (hook *ElasticsearchHook) Levels() []log.Level {
  return []log.Level{
    log.InfoLevel,
  }
}
