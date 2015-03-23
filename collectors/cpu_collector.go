package collectors

import "fmt"
import "github.com/c9s/goprocinfo/linux"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"

type CpuCollector struct{}

func (c *CpuCollector) Collect() (map[string]mappings.MetricMap, error) {
	stat, err := linux.ReadStat("/proc/stat")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	cpuMapping := map[string]mappings.MetricMap{}

	for _, s := range stat.CPUStats {
		cpuMapping[s.Id] = mappings.MetricMap{
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

func (c *CpuCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	data, _ := c.Collect()

	if data != nil {
		for cpu, values := range data {
			for k, v := range values {
				metric := structs.BuildMetric("cpu", "gauge_pct", k, v, structs.FieldsMap{
					"core": cpu,
					"unit": "Jiff",
				})
				metric.Path = fmt.Sprintf("cpu.%s", cpu)
				report = append(report, metric)
			}
		}
	}

	return report, nil
}
