package main

import "encoding/json"
import "fmt"
import "net/http"
import log "github.com/Sirupsen/logrus"

type ElasticsearchHook struct{}

func (hook *ElasticsearchHook) Fire(entry *log.Entry) error {
	data := MarshalData(entry)

	serialized, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		return nil
	}

	entry_index := Getenv("ELASTICSEARCH_INDEX", "logstash-data")
	entry_type := Getenv("METRIC_TYPE", "metricsd")
	status := ElasticsearchPost(fmt.Sprintf("/%s/%s", entry_index, entry_type), serialized)
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
