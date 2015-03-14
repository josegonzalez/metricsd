package main

import "bytes"
import "encoding/json"
import "fmt"
import "net/http"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type ElasticsearchShipper struct{}

var elasticsearchUrl string
var index string
var metricType string

func (shipper *ElasticsearchShipper) Setup(conf ini.File) {
	elasticsearchUrl = "http://127.0.0.1:9200"
	useElasticsearchUrl, ok := conf.Get("ElasticsearchShipper", "url")
	if ok {
		elasticsearchUrl = useElasticsearchUrl
	}

	index = "metricsd-data"
	useIndex, ok := conf.Get("ElasticsearchShipper", "enabled")
	if ok {
		index = useIndex
	}

	metricType = "metricsd"
	useMetricType, ok := conf.Get("ElasticsearchShipper", "type")
	if ok {
		metricType = useMetricType
	}

	SetupTemplate()
}

func (shipper *ElasticsearchShipper) Ship(logs MetricMapSlice) error {
	action := ActionMap{
		"index": IndexMap{
			"_index": index,
			"_type":  metricType,
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
		serialized := MarshalData(item)
		slice = Extend(slice, serializedAction)
		slice = Extend(slice, newline)
		slice = Extend(slice, serialized)
		slice = Extend(slice, newline)
	}

	status, err := ElasticsearchPost("/_bulk", slice)
	if err != nil {
		logrus.Warning("Indexing serialized data failed with err: ", err)
	}

	if status != http.StatusOK {
		logrus.Warning("Indexing serialized data failed with status: ", status)
	}
	return nil
}

func ElasticsearchPost(url string, data []byte) (int, error) {
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
	"template": "metricsd-*",
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

	status, err := ElasticsearchPost("/_template/metricsd", data)
	if err != nil {
		logrus.Fatal("Creating index failed: ", err)
	}

	if status != http.StatusOK {
		logrus.Fatal("Creating index failed")
	}
}
