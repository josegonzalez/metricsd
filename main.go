package main

import "os"
import "sync"
import log "github.com/Sirupsen/logrus"

func main() {
	initializeLogging()
	SetupTemplate()

	// iostat: (diskstat.go + mangling) /proc/diskstats
	// sockets: (sockstat.go in a pr) /proc/net/sockstat

  collectors := []CollectorInterface{
		&CpuCollector{},
		&DiskspaceCollector{},
		&LoadAvgCollector{},
		&MemoryCollector{},
		&VmstatCollector{},
	}

	var c chan log.Fields = make(chan log.Fields)
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
		report(c)
	}()

	collector_wg.Wait()
	close(c)
	reporter_wg.Wait()
}

func initializeLogging() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&LogstashFormatter{})
	log.AddHook(new(ElasticsearchHook))
}

func collect(c chan log.Fields, collector CollectorInterface) {
	data, err := collector.Report()
	if err != nil {
		close(c)
		return
	}

	for _, element := range data {
		c <- element
	}
}

func report(c chan log.Fields) {
	var list []log.Fields

	for item := range c {
		log.WithFields(item).Info()
		list = append(list, item)

		if len(list) == 10 {
			// TODO: Ship list
			list = nil
		}
	}

  if len(list) > 0 {
  	// TODO: Ship list
		list = nil
	}
}

func errHndlr(err error) {
	if err != nil {
		log.Error("error:", err)
		os.Exit(1)
	}
}
