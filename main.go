package main

import "fmt"
import "sync"
import "github.com/josegonzalez/metricsd/collectors"
import "github.com/josegonzalez/metricsd/config"
import "github.com/josegonzalez/metricsd/shippers"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

func main() {
	conf := config.Setup()
	initializeLogging(conf)
	shippers := getShippers(conf)
	collectorList := getCollectors(conf)

	var c chan *structs.Metric = make(chan *structs.Metric)
	var collector_wg sync.WaitGroup
	var reporter_wg sync.WaitGroup
	collector_wg.Add(len(collectorList))
	reporter_wg.Add(1)

	for _, collector := range collectorList {
		go func(collector collectors.CollectorInterface) {
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
	if config.LogLevel == "panic" {
		logrus.SetLevel(logrus.PanicLevel)
	} else if config.LogLevel == "fatal" {
		logrus.SetLevel(logrus.FatalLevel)
	} else if config.LogLevel == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
	} else if config.LogLevel == "warning" {
		logrus.SetLevel(logrus.WarnLevel)
	} else if config.LogLevel == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if config.LogLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func collect(c chan *structs.Metric, collector collectors.CollectorInterface) {
	data, err := collector.Report()
	if err != nil {
		close(c)
		return
	}

	for _, element := range data {
		c <- element
	}
}

func report(c chan *structs.Metric, shippers []shippers.ShipperInterface) {
	var list structs.MetricSlice

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

func getShippers(conf ini.File) []shippers.ShipperInterface {
	var shipperList []shippers.ShipperInterface
	var enabled string

	enabled, _ = conf.Get("GraphiteShipper", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling GraphiteShipper")
		graphiteShipper := &shippers.GraphiteShipper{}
		graphiteShipper.Setup(conf)
		shipperList = append(shipperList, graphiteShipper)
	}

	enabled, _ = conf.Get("LogstashElasticsearchShipper", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling LogstashElasticsearchShipper")
		elasticsearchShipper := &shippers.LogstashElasticsearchShipper{}
		elasticsearchShipper.Setup(conf)
		shipperList = append(shipperList, elasticsearchShipper)
	}

	enabled, _ = conf.Get("StdoutShipper", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling StdoutShipper")
		stdoutShipper := &shippers.StdoutShipper{}
		stdoutShipper.Setup(conf)
		shipperList = append(shipperList, stdoutShipper)
	}

	enabled, _ = conf.Get("LogstashRedisShipper", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling LogstashRedisShipper")
		redisShipper := &shippers.LogstashRedisShipper{}
		redisShipper.Setup(conf)
		shipperList = append(shipperList, redisShipper)
	}

	return shipperList
}

func getCollectors(conf ini.File) []collectors.CollectorInterface {
	var collectorList []collectors.CollectorInterface
	var enabled string

	// iostat: (diskstat.go + mangling) /proc/diskstats

	enabled, _ = conf.Get("CpuCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling CpuCollector")
		collectorList = append(collectorList, &collectors.CpuCollector{})
	}

	enabled, _ = conf.Get("DiskspaceCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling DiskspaceCollector")
		collectorList = append(collectorList, &collectors.DiskspaceCollector{})
	}

	enabled, _ = conf.Get("IostatCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling IostatCollector")
		collectorList = append(collectorList, &collectors.IostatCollector{})
	}

	enabled, _ = conf.Get("LoadAvgCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling LoadAvgCollector")
		collectorList = append(collectorList, &collectors.LoadAvgCollector{})
	}

	enabled, _ = conf.Get("MemoryCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling MemoryCollector")
		collectorList = append(collectorList, &collectors.MemoryCollector{})
	}

	enabled, _ = conf.Get("SocketsCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling SocketsCollector")
		collectorList = append(collectorList, &collectors.SocketsCollector{})
	}

	enabled, _ = conf.Get("VmstatCollector", "enabled")
	if enabled == "true" {
		logrus.Debug("enabling VmstatCollector")
		collectorList = append(collectorList, &collectors.VmstatCollector{})
	}

	return collectorList
}
