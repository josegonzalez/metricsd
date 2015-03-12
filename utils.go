package main

import "bytes"
import "fmt"
import "net/http"
import "syscall"
import log "github.com/Sirupsen/logrus"

func Extend(slice []byte, sliceTwo []byte) []byte {
	for i := range sliceTwo {
		slice = append(slice, sliceTwo[i])
	}

	return slice
}

func Getenv(key string, def string) string {
	v, err := syscall.Getenv(key)
	if err == true {
		return def
	}
	if v == "" {
		return def
	}
	return v
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
