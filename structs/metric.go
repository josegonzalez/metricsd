package structs

import "encoding/json"
import "fmt"
import "os"
import "time"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type Metric struct {
	Collector  string
	Path       string
	From       string
	Name       string
	Value      interface{}
	RawValue   interface{}
	Timestamp  time.Time
	Precision  int
	Host       string
	MetricType string
	TTL        int
	Data       FieldsMap
}

var Hostname string

func init() {
	var err error
	Hostname, err = os.Hostname()

	if err != nil {
		logrus.Warning("error retrieving hostname, using `localhost` for now")
		Hostname = "localhost"
	}
}

type FieldsMap map[string]interface{}

type MetricSlice []*Metric

func BuildMetric(collector string, from string, metricType string, name string, value interface{}, data FieldsMap) *Metric {
	return &Metric{
		Collector:  collector,
		From:       from,
		Name:       name,
		Value:      value,
		MetricType: metricType,
		Timestamp:  time.Now(),
		Host:       Hostname,
		Precision:  0,
		TTL:        0,
		Data:       data,
	}
}

func (m *Metric) Process(conf ini.File) {
	if hostname, ok := conf.Get("metricsd", "hostname"); ok {
		m.Host = hostname
	}

	if hostname, ok := conf.Get(m.Collector, "hostname"); ok {
		m.Host = hostname
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

func (m *Metric) ToJSON() []byte {
	data := m.ToMap()

	serialized, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		return nil
	}
	return serialized
}

func (m *Metric) ToGraphite(prefix string) (response string) {
	path := m.From
	if m.Path != "" {
		path = m.Path
	}
	key := fmt.Sprintf("%s.%s.%s", m.Host, path, m.Name)
	if prefix != "" {
		key = fmt.Sprintf("%s%s", prefix, key)
	}
	return fmt.Sprintf("%s %v %d", key, m.Value, int32(m.Timestamp.Unix()))
}
