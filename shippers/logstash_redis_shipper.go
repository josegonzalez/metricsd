package shippers

import "github.com/josegonzalez/go-radixurl"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

// LogstashRedisShipper is an exported type that
// allows shipping metrics to redis in logstash format
type LogstashRedisShipper struct {
	enabled bool
	list    string
	url     string
}

// Enabled allows checking whether the shipper is enabled or not
func (s *LogstashRedisShipper) Enabled() bool {
	return s.enabled
}

// State allows setting the enabled state of the shipper
func (s *LogstashRedisShipper) State(state bool) {
	s.enabled = state
}

// Setup configures the shipper
func (s *LogstashRedisShipper) Setup(conf ini.File) {
	s.State(true)

	if list, ok := conf.Get("LogstashRedisShipper", "list"); ok {
		s.list = list
	} else {
		s.list = "metricsd"
	}

	if url, ok := conf.Get("LogstashRedisShipper", "url"); ok {
		s.url = url
	} else {
		s.url = "redis://127.0.0.1:6379/0"
	}
}

// Ship sends a list of MetricSlices to redis
func (s *LogstashRedisShipper) Ship(logs structs.MetricSlice) error {
	c, err := radixurl.ConnectToURL(s.url)
	errHndlr(err)
	defer c.Close()

	var list []string

	for _, item := range logs {
		serialized := item.ToJson()
		list = append(list, string(serialized))
	}

	length := len(logs)
	if length == 10 {
		r := c.Cmd("rpush", s.list, list[0], list[1], list[2], list[3], list[4], list[5], list[6], list[7], list[8], list[9])
		errHndlr(r.Err)
	} else if length == 9 {
		r := c.Cmd("rpush", s.list, list[0], list[1], list[2], list[3], list[4], list[5], list[6], list[7], list[8])
		errHndlr(r.Err)
	} else if length == 8 {
		r := c.Cmd("rpush", s.list, list[0], list[1], list[2], list[3], list[4], list[5], list[6], list[7])
		errHndlr(r.Err)
	} else if length == 7 {
		r := c.Cmd("rpush", s.list, list[0], list[1], list[2], list[3], list[4], list[5], list[6])
		errHndlr(r.Err)
	} else if length == 6 {
		r := c.Cmd("rpush", s.list, list[0], list[1], list[2], list[3], list[4], list[5])
		errHndlr(r.Err)
	} else if length == 5 {
		r := c.Cmd("rpush", s.list, list[0], list[1], list[2], list[3], list[4])
		errHndlr(r.Err)
	} else if length == 4 {
		r := c.Cmd("rpush", s.list, list[0], list[1], list[2], list[3])
		errHndlr(r.Err)
	} else if length == 3 {
		r := c.Cmd("rpush", s.list, list[0], list[1], list[2])
		errHndlr(r.Err)
	} else if length == 2 {
		r := c.Cmd("rpush", s.list, list[0], list[1])
		errHndlr(r.Err)
	} else if length == 1 {
		r := c.Cmd("rpush", s.list, list[0])
		errHndlr(r.Err)
	}

	return nil
}

func errHndlr(err error) {
	if err != nil {
		logrus.Fatal("redis error: ", err)
	}
}
