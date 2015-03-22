package collectors

import "github.com/c9s/goprocinfo/linux"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/Sirupsen/logrus"

type MemoryCollector struct{}

func (c *MemoryCollector) Collect() (mappings.MetricMap, error) {
	stat, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	return mappings.MetricMap{
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

func (c *MemoryCollector) Report() (mappings.MetricMapSlice, error) {
	var report mappings.MetricMapSlice
	values, _ := c.Collect()

	if values != nil {
		for k, v := range values {
			report = append(report, mappings.MetricMap{
				"_from":       "memory",
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
