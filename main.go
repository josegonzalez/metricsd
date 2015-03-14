package main

import "fmt"
import "sync"
import "github.com/vaughan0/go-ini"
import "github.com/Sirupsen/logrus"

func main() {
	conf := Config()
	initializeLogging(conf)
	SetupTemplate()
	shippers := shippers(conf)

	// iostat: (diskstat.go + mangling) /proc/diskstats
	// sockets: (sockstat.go in a pr) /proc/net/sockstat

	collectors := []CollectorInterface{
		&CpuCollector{},
		&DiskspaceCollector{},
		&LoadAvgCollector{},
		&MemoryCollector{},
		&VmstatCollector{},
	}

	var c chan MetricMap = make(chan MetricMap)
	var collector_wg sync.WaitGroup
	var reporter_wg sync.WaitGroup
	collector_wg.Add(len(collectors))
	reporter_wg.Add(1)

	for _, collector := range collectors {
		go func(collector CollectorInterface) {
			defer collector_wg.Done()
			collect(c, collector)
		}(collector)
	}

	go func() {
		defer reporter_wg.Done()
		report(c, shippers)
	}()

	collector_wg.Wait()
	close(c)
	reporter_wg.Wait()
}

func initializeLogging(conf ini.File) {
	logrus.SetLevel(logrus.DebugLevel)
}

func collect(c chan MetricMap, collector CollectorInterface) {
	data, err := collector.Report()
	if err != nil {
		close(c)
		return
	}

	for _, element := range data {
		c <- element
	}
}

func report(c chan MetricMap, shippers []ShipperInterface) {
	var list MetricMapSlice

	for item := range c {
		list = append(list, item)

		if len(list) == 10 {
			logrus.Debug(fmt.Sprintf("Shipping %d messages", len(list)))
			for _, shipper := range shippers {
				shipper.Ship(list)
			}
			list = nil
		}
	}

	if len(list) > 0 {
		logrus.Debug(fmt.Sprintf("Shipping %d messages", len(list)))
		for _, shipper := range shippers {
			shipper.Ship(list)
		}
		list = nil
	}
}

func shippers(conf ini.File) []ShipperInterface {
	var shippers []ShipperInterface
	var enabled string

	enabled, _ = conf.Get("ElasticsearchShipper", "enabled")
	if enabled == "true" {
		logrus.Info("enabling ElasticsearchShipper")
		shippers = append(shippers, &ElasticsearchShipper{})
	}

	enabled, _ = conf.Get("StdoutShipper", "enabled")
	if enabled == "true" {
		logrus.Info("enabling StdoutShipper")
		shippers = append(shippers, &StdoutShipper{})
	}

	enabled, _ = conf.Get("RedisShipper", "enabled")
	if enabled == "true" {
		logrus.Info("enabling RedisShipper")
		shippers = append(shippers, &RedisShipper{})
	}


	return shippers
}
