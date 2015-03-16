package shippers

import "encoding/json"
import "fmt"
import "os"
import "time"
import "github.com/josegonzalez/metricsd/mappings"

func MarshalData(data mappings.MetricMap) []byte {
	data["@version"] = "1"
	data["@timestamp"] = time.Now().Format("2006-01-02T15:04:05.000Z")

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	data["host"] = hostname

	serialized, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		return nil
	}
	return serialized
}
