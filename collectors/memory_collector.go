package collectors

import "github.com/c9s/goprocinfo/linux"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type MemoryCollector struct {
	enabled bool
}

func (this *MemoryCollector) Enabled() bool {
	return this.enabled
}

func (this *MemoryCollector) State(state bool) {
	this.enabled = state
}

func (this *MemoryCollector) Setup(conf ini.File) {
	this.State(true)
}

func (this *MemoryCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	values, _ := this.collect()

	if values != nil {
		for k, v := range values {
			metric := structs.BuildMetric("MemoryCollector", "memory", "gauge", k, v, structs.FieldsMap{
				"unit":      "B",
				"where":     "system_memory",
				"raw_key":   k,
				"raw_value": v,
			})
			report = append(report, metric)
		}
	}

	return report, nil
}

func (this *MemoryCollector) collect() (mappings.MetricMap, error) {
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
