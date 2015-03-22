package collectors

import "github.com/c9s/goprocinfo/linux"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/Sirupsen/logrus"

type VmstatCollector struct{}

func (c *VmstatCollector) Collect() (mappings.MetricMap, error) {
	stat, err := linux.ReadVMStat("/proc/vmstat")
	if err != nil {
		logrus.Fatal("stat read fail")
		return nil, err
	}

	return mappings.MetricMap{
		"paging_in": stat.PagePagein,
		"pagingout": stat.PagePageout,
		"swap_in":   stat.PageSwapin,
		"swap_out":  stat.PageSwapout,
	}, nil
}

func (c *VmstatCollector) Report() (mappings.MetricMapSlice, error) {
	var report mappings.MetricMapSlice
	values, _ := c.Collect()

	if values != nil {
		for k, v := range values {
			report = append(report, mappings.MetricMap{
				"_from":       "vmstat",
				"target_type": "rate",
				"type":        k,
				"unit":        "Page",
				"result":      v,
			})
		}
	}

	return report, nil
}
