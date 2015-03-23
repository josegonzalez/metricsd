package shippers

import "fmt"
import "net"
import "net/url"
import "strings"
import "time"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type GraphiteShipper struct{}

var debug bool
var graphiteHost string
var graphitePort string
var prefix string

func (shipper *GraphiteShipper) Setup(conf ini.File) {
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
				logrus.Warning("Error parsing graphite url")
				logrus.Warning("Using default 127.0.0.1:2003 for graphite url")
			case len(splitted) > 1:
				graphiteHost, graphitePort = splitted[0], splitted[1]
			default:
				graphiteHost, graphitePort = splitted[0], "2003"
			}
		} else {
			logrus.Warning("Error parsing graphite url: %s", err)
			logrus.Warning("Using default 127.0.0.1:2003 for graphite url")
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
		logrus.Warning("Connecting to graphite failed with err: ", err)
		return err
	}
	defer con.Close()

	for _, item := range logs {
		serialized := item.ToGraphite()
		if debug {
			fmt.Printf("%s%s\n", prefix, serialized)
		}
		fmt.Fprintln(con, serialized)
	}

	return nil
}
