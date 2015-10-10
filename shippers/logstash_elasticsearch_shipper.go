package shippers

import "bytes"
import "encoding/json"
import "fmt"
import "net/http"
import "github.com/Sirupsen/logrus"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/josegonzalez/metricsd/utils"
import "github.com/vaughan0/go-ini"

type actionMap map[string]indexMap
type indexMap map[string]string

// LogstashElasticsearchShipper is an exported type that
// allows shipping metrics to elasticsearch in logstash format
type LogstashElasticsearchShipper struct {
	enabled    bool
	index      string
	metricType string
	url        string
}

// Enabled allows checking whether the shipper is enabled or not
func (s *LogstashElasticsearchShipper) Enabled() bool {
	return s.enabled
}

// State allows setting the enabled state of the shipper
func (s *LogstashElasticsearchShipper) State(state bool) {
	s.enabled = state
}

// Setup configures the shipper
func (s *LogstashElasticsearchShipper) Setup(conf ini.File) {
	s.State(true)

	if url, ok := conf.Get("LogstashElasticsearchShipper", "url"); ok {
		s.url = url
	} else {
		s.url = "http://127.0.0.1:9200"
	}

	if index, ok := conf.Get("LogstashElasticsearchShipper", "enabled"); ok {
		s.index = index
	} else {
		s.index = "metricsd-data"
	}

	if metricType, ok := conf.Get("LogstashElasticsearchShipper", "type"); ok {
		s.metricType = metricType
	} else {
		s.metricType = "metricsd"
	}

	s.setupTemplate()
}

// Ship sends a list of MetricSlices to elasticsearch
func (s *LogstashElasticsearchShipper) Ship(logs structs.MetricSlice) error {
	action := actionMap{
		"index": indexMap{
			"_index": s.index,
			"_type":  s.metricType,
		},
	}
	serializedAction, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("Failed to marshal action to JSON, %v", err)
	}

	var slice []byte
	newline := []byte("\n")

	for _, item := range logs {
		serialized := item.ToJson()
		slice = utils.Extend(slice, serializedAction)
		slice = utils.Extend(slice, newline)
		slice = utils.Extend(slice, serialized)
		slice = utils.Extend(slice, newline)
	}

	status, err := s.elasticsearchPost("/_bulk", slice)
	if err != nil {
		logrus.Warning("indexing serialized data failed with err: ", err)
	}

	if status != http.StatusOK {
		logrus.Warning("indexing serialized data failed with status: ", status)
	}
	return nil
}

func (s *LogstashElasticsearchShipper) elasticsearchPost(url string, data []byte) (int, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", s.url, url), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("Failed to make request, %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (s *LogstashElasticsearchShipper) setupTemplate() {
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

	status, err := s.elasticsearchPost("/_template/metricsd", data)
	if err != nil {
		logrus.Fatal("creating index failed: ", err)
	}

	if status != http.StatusOK {
		logrus.Fatal("creating index failed")
	}
}
