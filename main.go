package main

import "os"
import log "github.com/Sirupsen/logrus"

func main() {
	log.SetFormatter(&LogstashFormatter{})
	log.SetLevel(log.InfoLevel)
	log.Info("metrics_d")
	log.AddHook(new(ElasticsearchHook))

	SetupTemplate()

	cpu := &CpuCollector{}
	cpu.Report()

	// diskspace: (mounts.go + statvfs) /proc/mounts
	diskspace := &DiskspaceCollector{}
	diskspace.Report()

	// iostat: (diskstat.go + mangling) /proc/diskstats

	loadavg := &LoadAvgCollector{}
	loadavg.Report()

	memory := &MemoryCollector{}
	memory.Report()

	// sockets: (sockstat.go in a pr) /proc/net/sockstat

	vmstat := &VmstatCollector{}
	vmstat.Report()
}

func errHndlr(err error) {
	if err != nil {
		log.Error("error:", err)
		os.Exit(1)
	}
}
