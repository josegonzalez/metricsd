package main

import log "github.com/Sirupsen/logrus"
import linuxproc "github.com/c9s/goprocinfo/linux"

type VmstatCollector struct{}

func (c *VmstatCollector) Collect() (IntMetricMapping, error) {
	stat, err := linuxproc.ReadVMStat("/proc/vmstat")
	if err != nil {
		log.Fatal("stat read fail")
		return nil, err
	}

	return IntMetricMapping{
		"paging_in": stat.PagePagein,
		"pagingout": stat.PagePageout,
		"swap_in":   stat.PageSwapin,
		"swap_out":  stat.PageSwapout,
	}, nil
}

func (c *VmstatCollector) Report() ([]log.Fields, error) {
	var report []log.Fields
	values, _ := c.Collect()

	if values != nil {
		for k, v := range values {
			report = append(report, log.Fields{
				"target_type": "rate",
				"type":        k,
				"unit":        "Page",
				"result":      v,
			})
		}
	}

	return report, nil
}
