package main

import "encoding/json"
import "fmt"
import "net/http"
import log "github.com/Sirupsen/logrus"

type ElasticsearchShipper struct{}

func (hook *ElasticsearchShipper) Ship(logs []log.Fields) error {
	log.Debug(fmt.Sprintf("Shipping %d logs", len(logs)))
	index := Getenv("ELASTICSEARCH_INDEX", "logstash-data")
	metric_type := Getenv("METRIC_TYPE", "metricsd")

	action := ActionMap{
		"index": IndexMap{
			"_index": index,
			"_type": metric_type,
		},
	}
	serializedAction, err := json.Marshal(action)
	if err != nil {
		fmt.Errorf("Failed to marshal action to JSON, %v", err)
		return nil
	}

	var slice []byte
	newline := []byte("\n")

	for _, item := range logs {
		serialized, err := json.Marshal(item)
		if err != nil {
			fmt.Errorf("Failed to marshal fields to JSON, %v", err)
			return nil
		}
		slice = extend(slice, serializedAction)
		slice = extend(slice, newline)
		slice = extend(slice, serialized)
		slice = extend(slice, newline)
	}

	status, err := ElasticsearchPost("/_bulk", slice)
	if err != nil {
		log.Warning("Indexing serialized data failed with err: ", err)
	}

	if status != http.StatusOK {
		log.Warning("Indexing serialized data failed with status: ", status)
	}
	return nil
}

func extend(slice []byte, sliceTwo []byte) []byte {
  for i := range sliceTwo {
  	slice = append(slice, sliceTwo[i])
   }

   return slice
}
