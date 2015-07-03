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
	Ttl        int
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
		Ttl:        0,
		Data:       data,
	}
}

func (this *Metric) Process(conf ini.File) {
	if hostname, ok := conf.Get("metricsd", "hostname"); ok {
		this.Host = hostname
	}

	if hostname, ok := conf.Get(this.Collector, "hostname"); ok {
		this.Host = hostname
	}
}

func (this *Metric) ToMap() map[string]interface{} {
	data := make(map[string]interface{})
	data["@timestamp"] = this.Timestamp.Format("2006-01-02T15:04:05.000Z")
	data["@version"] = "1"
	data["collector"] = this.From
	data["type"] = this.Name
	data["result"] = this.Value
	data["target_type"] = this.MetricType

	for k, v := range this.Data {
		_, exists := data[k]
		if exists {
			data[fmt.Sprintf("fields.%s", k)] = v
		} else {
			data[k] = v
		}
	}

	if _, ok := data["host"]; !ok {
		data["host"] = this.Host
	}

	return data
}

func (this *Metric) ToJson() []byte {
	data := this.ToMap()

	serialized, err := json.Marshal(data)
	if err != nil {
		fmt.Errorf("Failed to marshal fields to JSON, %v", err)
		return nil
	}
	return serialized
}

func (this *Metric) ToGraphite(prefix string) (response string) {
	path := this.From
	if this.Path != "" {
		path = this.Path
	}
	key := fmt.Sprintf("%s.%s.%s", this.Host, path, this.Name)
	if prefix != "" {
		key = fmt.Sprintf("%s%s", prefix, key)
	}
	return fmt.Sprintf("%s %v %d", key, this.Value, int32(this.Timestamp.Unix()))
}
