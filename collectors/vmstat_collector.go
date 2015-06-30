package collectors

import "github.com/c9s/goprocinfo/linux"
import "github.com/josegonzalez/metricsd/mappings"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

type VmstatCollector struct{}

func (c *VmstatCollector) Setup(conf ini.File) {
}

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

func (c *VmstatCollector) Report() (structs.MetricSlice, error) {
	var report structs.MetricSlice
	values, _ := c.Collect()

	if values != nil {
		for k, v := range values {
			metric := structs.BuildMetric("vmstat", "rate", k, v, structs.FieldsMap{
				"unit":      "Page",
				"raw_key":   k,
				"raw_value": v,
			})
			report = append(report, metric)
		}
	}

	return report, nil
}
