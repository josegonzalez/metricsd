package main

import log "github.com/Sirupsen/logrus"
import linuxproc "github.com/c9s/goprocinfo/linux"

type CpuCollector struct{}

func (c *CpuCollector) Collect() (map[string]IntMetricMap, error) {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
		return nil, err
	}

	cpuMapping := map[string]IntMetricMap{}

	for _, s := range stat.CPUStats {
		cpuMapping[s.Id] = IntMetricMap{
			"user":       s.User,
			"nice":       s.Nice,
			"system":     s.System,
			"idle":       s.Idle,
			"iowait":     s.IOWait,
			"irq":        s.IRQ,
			"softirq":    s.SoftIRQ,
			"steal":      s.Steal,
			"guest":      s.Guest,
			"guest_nice": s.GuestNice,
		}
	}

	return cpuMapping, nil
}

func (c *CpuCollector) Report() (MetricMapSlice, error) {
	var report MetricMapSlice
	data, _ := c.Collect()

	if data != nil {
		for cpu, values := range data {
			for k, v := range values {
				report = append(report, MetricMap{
					"target_type": "gauge_pct",
					"core":        cpu,
					"type":        k,
					"unit":        "Jiff",
					"result":      v,
				})
			}
		}
	}

	return report, nil
}
