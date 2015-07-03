package main

import "fmt"
import "reflect"
import "strconv"
import "strings"
import "sync"
import "time"
import "github.com/josegonzalez/metricsd/collectors"
import "github.com/josegonzalez/metricsd/config"
import "github.com/josegonzalez/metricsd/shippers"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

var conf ini.File

func main() {
	conf = config.Setup()
	initializeLogging()
	shippers := getShippers()
	collectorList := getCollectors()
	interval := getInterval()
	loop, _ := conf.Get("metricsd", "loop")

	runCollect(shippers, collectorList)
	if loop == "true" {
		for _ = range time.Tick(interval) {
			runCollect(shippers, collectorList)
		}
	}
}

func getInterval() time.Duration {
	defaultInterval := 30
	interval, ok := conf.Get("metricsd", "interval")

	if ok {
		interval, err := strconv.Atoi(interval)
		if err == nil {
			return time.Duration(interval) * time.Second
		}
	}
	return time.Duration(defaultInterval) * time.Second
}

func runCollect(shippers []shippers.ShipperInterface, collectorList []collectors.CollectorInterface) {
	var c chan *structs.Metric = make(chan *structs.Metric)
	var collector_wg sync.WaitGroup
	var reporter_wg sync.WaitGroup
	var active = 0

	for _, collector := range collectorList {
		if collector.Enabled() {
			active += 1
		}
	}

	collector_wg.Add(active)
	reporter_wg.Add(1)

	for _, collector := range collectorList {
		if !collector.Enabled() {
			continue
		}

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

func initializeLogging() {
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
		item.Process(conf)
		list = append(list, item)

		if len(list) == 10 {
			logrus.Debug(fmt.Sprintf("shipping %d messages", len(list)))
			for _, shipper := range shippers {
				if shipper.Enabled() {
					shipper.Ship(list)
				}
			}
			list = nil
		}
	}

	if len(list) > 0 {
		logrus.Debug(fmt.Sprintf("shipping %d messages", len(list)))
		for _, shipper := range shippers {
			if shipper.Enabled() {
				shipper.Ship(list)
			}
		}
		list = nil
	}
}

func getShippers() []shippers.ShipperInterface {
	var shipperList []shippers.ShipperInterface
	var enabled string

	shipperList = append(shipperList, &shippers.GraphiteShipper{})
	shipperList = append(shipperList, &shippers.LogstashElasticsearchShipper{})
	shipperList = append(shipperList, &shippers.StdoutShipper{})
	shipperList = append(shipperList, &shippers.LogstashRedisShipper{})

	for _, shipper := range shipperList {
		collectorName := strings.Split(reflect.TypeOf(shipper).String(), ".")[1]
		enabled, _ = conf.Get(collectorName, "enabled")
		if enabled == "true" {
			logrus.Debug(fmt.Sprintf("enabling %s", collectorName))
			shipper.Setup(conf)
			shipper.State(true)
		} else {
			shipper.State(false)
		}
	}

	return shipperList
}

func getCollectors() []collectors.CollectorInterface {
	var collectorList []collectors.CollectorInterface
	var enabled string

	// iostat: (diskstat.go + mangling) /proc/diskstats
	collectorList = append(collectorList, &collectors.CpuCollector{})
	collectorList = append(collectorList, &collectors.DiskspaceCollector{})
	collectorList = append(collectorList, &collectors.IostatCollector{})
	collectorList = append(collectorList, &collectors.LoadAvgCollector{})
	collectorList = append(collectorList, &collectors.MemoryCollector{})
	collectorList = append(collectorList, &collectors.RedisCollector{})
	collectorList = append(collectorList, &collectors.SocketsCollector{})
	collectorList = append(collectorList, &collectors.VmstatCollector{})

	for _, collector := range collectorList {
		collectorName := strings.Split(reflect.TypeOf(collector).String(), ".")[1]
		enabled, _ = conf.Get(collectorName, "enabled")
		if enabled == "true" {
			logrus.Debug(fmt.Sprintf("enabling %s", collectorName))
			collector.Setup(conf)
			collector.State(true)
		} else {
			collector.State(false)
		}
	}

	return collectorList
}
