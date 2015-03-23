package structs

import "encoding/json"
import "fmt"
import "os"
import "time"
import "github.com/Sirupsen/logrus"

type Metric struct {
	Path       string
	From       string
	Name       string
	Value      interface{}
	RawValue   interface{}
	Timestamp  time.Time
	Precision  int
	Host       string
	MetricType string
	Ttl        int
	Data       FieldsMap
}

var Hostname string

func init() {
	var err error
	Hostname, err = os.Hostname()

	if err != nil {
		logrus.Warning("Error retrieving hostname, using `localhost` for now")
		Hostname = "localhost"
	}
}

type FieldsMap map[string]interface{}

type MetricSlice []*Metric

func BuildMetric(from string, metricType string, name string, value interface{}, data FieldsMap) *Metric {
	return &Metric{
		From:       from,
		Name:       name,
		Value:      value,
		MetricType: metricType,
		Timestamp:  time.Now(),
		Host:       Hostname,
		Precision:  0,
		Ttl:        0,
		Data:       data,
	}
}

func (m *Metric) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["@timestamp"] = m.Timestamp.Format("2006-01-02T15:04:05.000Z")
	data["@version"] = "1"
	data["collector"] = m.From
	data["type"] = m.Name
	data["result"] = m.Value
	data["target_type"] = m.MetricType

	for k, v := range m.Data {
		_, exists := data[k]
		if exists {
		    data[fmt.Sprintf("fields.%s", k)] = v
		} else {
			data[k] = v
		}
	}

	if _, ok := data["host"]; !ok {
		data["host"] = m.Host
	}

	return data
}

func (m *Metric) ToJson() []byte {
	data := m.ToMap()

	serialized, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		return nil
	}
	return serialized
}

func (m *Metric) ToGraphite() (response string) {
	path := m.From
	if m.Path != "" {
		path = m.Path
	}
	key := fmt.Sprintf("%s.%s.%s", m.Host, path, m.Name)
	return fmt.Sprintf("%s %v %d", key, m.Value, int32(m.Timestamp.Unix()))
}
