package main

import log "github.com/Sirupsen/logrus"
import linuxproc "github.com/c9s/goprocinfo/linux"

type MemoryCollector struct{}

func (c *MemoryCollector) Collect() (IntMetricMapping, error) {
	stat, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatal("stat read fail")
		return nil, err
	}

	return IntMetricMapping{
		"memory_total":  stat.MemTotal,
		"memory_free":   stat.MemFree,
		"buffers":       stat.Buffers,
		"cached":        stat.Cached,
		"active":        stat.Active,
		"dirty":         stat.Dirty,
		"inactive":      stat.Inactive,
		"shmem":         stat.Shmem,
		"swap_total":    stat.SwapTotal,
		"swap_free":     stat.SwapFree,
		"swap_cached":   stat.SwapCached,
		"vmalloc_total": stat.VmallocTotal,
		"vmalloc_used":  stat.VmallocUsed,
		"vmalloc_chunk": stat.VmallocChunk,
		"committed_as":  stat.Committed_AS,
	}, nil
}

func (c *MemoryCollector) Report() ([]log.Fields, error) {
	var report []log.Fields
	values, _ := c.Collect()

	if values != nil {
		for k, v := range values {
			report = append(report, log.Fields{
				"target_type": "gauge",
				"type":        k,
				"unit":        "B",
				"where":       "system_memory",
				"result":      v,
			})
		}
	}

	return report, nil
}
