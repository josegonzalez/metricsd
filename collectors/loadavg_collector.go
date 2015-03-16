package collectors

import "github.com/c9s/goprocinfo/linux"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/Sirupsen/logrus"

type LoadAvgCollector struct{}

func (c *LoadAvgCollector) Collect() (mappings.MetricMap, error) {
	stat, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	// TODO: Add processes_running and processes_total,
	// unit:processes, type:(running|total)
	return mappings.MetricMap{
		"01": stat.Last1Min,
		"05": stat.Last5Min,
		"15": stat.Last15Min,
	}, nil
}

func (c *LoadAvgCollector) Report() (mappings.MetricMapSlice, error) {
	var report mappings.MetricMapSlice
	values, _ := c.Collect()

	if values != nil {
		for k, v := range values {
			report = append(report, mappings.MetricMap{
				"target_type": "gauge",
				"type":        k,
				"unit":        "Load",
				"result":      v,
			})
		}
	}

	return report, nil
}
