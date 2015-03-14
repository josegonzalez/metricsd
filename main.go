package main

import "fmt"
import "sync"
import "github.com/vaughan0/go-ini"
import "github.com/Sirupsen/logrus"

func main() {
	conf := Setup()
	initializeLogging(conf)
	shippers := shippers(conf)
	collectors := collectors(conf)

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
	if LogLevel == "panic" {
		logrus.SetLevel(logrus.PanicLevel)
	} else if LogLevel == "fatal" {
		logrus.SetLevel(logrus.FatalLevel)
	} else if LogLevel == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
	} else if LogLevel == "warning" {
		logrus.SetLevel(logrus.WarnLevel)
	} else if LogLevel == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if LogLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
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
		logrus.Debug("enabling ElasticsearchShipper")
		elasticsearchShipper := &ElasticsearchShipper{}
		elasticsearchShipper.Setup(conf)
		shippers = append(shippers, elasticsearchShipper)
	}

	enabled, _ = conf.Get("StdoutShipper", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling StdoutShipper")
		stdoutShipper := &StdoutShipper{}
		stdoutShipper.Setup(conf)
		shippers = append(shippers, stdoutShipper)
	}

	enabled, _ = conf.Get("RedisShipper", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling RedisShipper")
		redisShipper := &RedisShipper{}
		redisShipper.Setup(conf)
		shippers = append(shippers, redisShipper)
	}

	return shippers
}

func collectors(conf ini.File) []CollectorInterface {
	var collectors []CollectorInterface
	var enabled string

	// iostat: (diskstat.go + mangling) /proc/diskstats
	// sockets: (sockstat.go in a pr) /proc/net/sockstat

	enabled, _ = conf.Get("CpuCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling CpuCollector")
		collectors = append(collectors, &CpuCollector{})
	}

	enabled, _ = conf.Get("DiskspaceCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling DiskspaceCollector")
		collectors = append(collectors, &DiskspaceCollector{})
	}

	enabled, _ = conf.Get("LoadAvgCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling LoadAvgCollector")
		collectors = append(collectors, &LoadAvgCollector{})
	}

	enabled, _ = conf.Get("MemoryCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling MemoryCollector")
		collectors = append(collectors, &MemoryCollector{})
	}

	enabled, _ = conf.Get("VmstatCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling VmstatCollector")
		collectors = append(collectors, &VmstatCollector{})
	}

	return collectors
}
