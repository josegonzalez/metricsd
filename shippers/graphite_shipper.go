package shippers

import "fmt"
import "net"
import "net/url"
import "strings"
import "time"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

// GraphiteShipper is an exported type that
// allows shipping metrics to graphite
type GraphiteShipper struct {
	debug   bool
	enabled bool
	host    string
	prefix  string
	port    string
}

// Enabled allows checking whether the shipper is enabled or not
func (s *GraphiteShipper) Enabled() bool {
	return s.enabled
}

// State allows setting the enabled state of the shipper
func (s *GraphiteShipper) State(state bool) {
	s.enabled = state
}

// Setup configures the shipper
func (s *GraphiteShipper) Setup(conf ini.File) {
	s.State(true)

	useDebug, ok := conf.Get("GraphiteShipper", "debug")
	if ok && useDebug == "true" {
		s.debug = true
	} else {
		s.debug = false
	}

	s.host = "127.0.0.1"
	s.port = "2003"
	useGraphiteURL, ok := conf.Get("GraphiteShipper", "url")
	if ok {
		graphiteURL, err := url.Parse(useGraphiteURL)
		if err == nil {
			splitted := strings.Split(graphiteURL.Host, ":")
			s.host, s.port = splitted[0], "2003"
			switch {
			case len(splitted) > 2:
				logrus.Warning("error parsing graphite url")
				logrus.Warning("using default 127.0.0.1:2003 for graphite url")
			case len(splitted) > 1:
				s.host, s.port = splitted[0], splitted[1]
			default:
				s.host, s.port = splitted[0], "2003"
			}
		} else {
			logrus.Warning("error parsing graphite url: %s", err)
			logrus.Warning("using default 127.0.0.1:2003 for graphite url")
		}
	}

	usePrefix, ok := conf.Get("GraphiteShipper", "prefix")
	if ok {
		s.prefix = fmt.Sprintf("%s.", usePrefix)
	}
}

// Ship sends a list of MetricSlices to graphite
func (s *GraphiteShipper) Ship(logs structs.MetricSlice) error {
	con, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", s.host, s.port), 1*time.Second)
	if err != nil {
		logrus.Warning("connecting to graphite failed with err: ", err)
		return err
	}
	defer con.Close()

	for _, item := range logs {
		serialized := item.ToGraphite(s.prefix)
		if s.debug {
			fmt.Printf("%s\n", serialized)
		}
		fmt.Fprintln(con, serialized)
	}

	return nil
}
