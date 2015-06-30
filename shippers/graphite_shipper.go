package shippers

import "fmt"
import "net"
import "net/url"
import "strings"
import "time"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type GraphiteShipper struct{
	enabled bool
}

var debug bool
var graphiteHost string
var graphitePort string
var prefix string

func (shipper *GraphiteShipper) Enabled() (bool) {
	return shipper.enabled
}

func (shipper *GraphiteShipper) State(state bool) {
	shipper.enabled = state
}

func (shipper *GraphiteShipper) Setup(conf ini.File) {
	shipper.State(true)

	useDebug, ok := conf.Get("GraphiteShipper", "debug")
	if ok && useDebug == "true" {
		debug = true
	} else {
		debug = false
	}

	graphiteHost = "127.0.0.1"
	graphitePort = "2003"
	useGraphiteUrl, ok := conf.Get("GraphiteShipper", "url")
	if ok {
		graphiteUrl, err := url.Parse(useGraphiteUrl)
		if err == nil {
			splitted := strings.Split(graphiteUrl.Host, ":")
			graphiteHost, graphitePort = splitted[0], "2003"
			switch {
			case len(splitted) > 2:
				logrus.Warning("error parsing graphite url")
				logrus.Warning("using default 127.0.0.1:2003 for graphite url")
			case len(splitted) > 1:
				graphiteHost, graphitePort = splitted[0], splitted[1]
			default:
				graphiteHost, graphitePort = splitted[0], "2003"
			}
		} else {
			logrus.Warning("error parsing graphite url: %s", err)
			logrus.Warning("using default 127.0.0.1:2003 for graphite url")
		}
	}

	usePrefix, ok := conf.Get("GraphiteShipper", "prefix")
	if ok {
		prefix = fmt.Sprintf("%s.", usePrefix)
	}
}

func (shipper *GraphiteShipper) Ship(logs structs.MetricSlice) error {
	con, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", graphiteHost, graphitePort), 1*time.Second)
	if err != nil {
		logrus.Warning("connecting to graphite failed with err: ", err)
		return err
	}
	defer con.Close()

	for _, item := range logs {
		serialized := item.ToGraphite(prefix)
		if debug {
			fmt.Printf("%s\n", serialized)
		}
		fmt.Fprintln(con, serialized)
	}

	return nil
}
