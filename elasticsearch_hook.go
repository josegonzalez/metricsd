package main

import "encoding/json"
import "net/http"
import "github.com/Sirupsen/logrus"

type ElasticsearchHook struct{}

func (hook *ElasticsearchHook) Fire(entry *logrus.Entry) error {
    data := MarshalData(entry)

    serialized, err := json.Marshal(data)
    if err != nil {
        logrus.Error("Failed to marshal fields to JSON, %v", err)
        return nil
    }

    status := ElasticsearchPost("/logstash-data/metricsd", serialized)
    if status != http.StatusCreated {
        logrus.Warning("Indexing serialized data failed, %s", err)
    }

    return nil
}

func (hook *ElasticsearchHook) Levels() []logrus.Level {
  return []logrus.Level{
    logrus.InfoLevel,
  }
}
