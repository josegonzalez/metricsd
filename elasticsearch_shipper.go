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
			"_type":  metric_type,
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
		slice = Extend(slice, serializedAction)
		slice = Extend(slice, newline)
		slice = Extend(slice, serialized)
		slice = Extend(slice, newline)
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

func ElasticsearchPost(url string, data []byte) (int, error) {
	elasticsearchUrl := Getenv("ELASTICSEARCH_URL", "http://127.0.0.1:9200")
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", elasticsearchUrl, url), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Failed to make request, %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func SetupTemplate() {
	template := `
{
	"order": 0,
	"template": "logstash-*",
	"settings": {
		"index.refresh_interval": "5s"
	},
	"mappings": {
		"_default_": {
			"dynamic_templates": [
				{
					"string_fields": {
						"mapping": {
							"index": "analyzed",
							"omit_norms": true,
							"type": "string",
							"fields": {
								"raw": {
									"index": "not_analyzed",
									"ignore_above": 256,
									"type": "string"
								}
							}
						},
						"match_mapping_type": "string",
						"match": "*"
					}
				}
			],
			"properties": {
				"geoip": {
					"dynamic": true,
					"path": "full",
					"properties": {
						"location": {
							"type": "geo_point"
						}
					},
					"type": "object"
				},
				"@version": {
					"index": "not_analyzed",
					"type": "string"
				}
			},
			"_all": {
				"enabled": true
			}
		}
	},
	"aliases": {}
}
`
	var data = []byte(template)

	status, err := ElasticsearchPost("/_template/logstash", data)
	if err != nil {
		log.Error("Indexing serialized data failed: ", err)
	}

	if status != http.StatusOK {
		log.Error("Creating index failed")
	}
}
