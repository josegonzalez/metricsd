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
	debug   bool
	enabled bool
	host    string
	prefix  string
	port    string
}

func (this *GraphiteShipper) Enabled() (bool) {
	return this.enabled
}

func (this *GraphiteShipper) State(state bool) {
	this.enabled = state
}

func (this *GraphiteShipper) Setup(conf ini.File) {
	this.State(true)

	useDebug, ok := conf.Get("GraphiteShipper", "debug")
	if ok && useDebug == "true" {
		this.debug = true
	} else {
		this.debug = false
	}

	this.host = "127.0.0.1"
	this.port = "2003"
	useGraphiteUrl, ok := conf.Get("GraphiteShipper", "url")
	if ok {
		graphiteUrl, err := url.Parse(useGraphiteUrl)
		if err == nil {
			splitted := strings.Split(graphiteUrl.Host, ":")
			this.host, this.port = splitted[0], "2003"
			switch {
			case len(splitted) > 2:
				logrus.Warning("error parsing graphite url")
				logrus.Warning("using default 127.0.0.1:2003 for graphite url")
			case len(splitted) > 1:
				this.host, this.port = splitted[0], splitted[1]
			default:
				this.host, this.port = splitted[0], "2003"
			}
		} else {
			logrus.Warning("error parsing graphite url: %s", err)
			logrus.Warning("using default 127.0.0.1:2003 for graphite url")
		}
	}

	usePrefix, ok := conf.Get("GraphiteShipper", "prefix")
	if ok {
		this.prefix = fmt.Sprintf("%s.", usePrefix)
	}
}

func (this *GraphiteShipper) Ship(logs structs.MetricSlice) error {
	con, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", this.host, this.port), 1*time.Second)
	if err != nil {
		logrus.Warning("connecting to graphite failed with err: ", err)
		return err
	}
	defer con.Close()

	for _, item := range logs {
		serialized := item.ToGraphite(this.prefix)
		if this.debug {
			fmt.Printf("%s\n", serialized)
		}
		fmt.Fprintln(con, serialized)
	}

	return nil
}
