package main

import log "github.com/Sirupsen/logrus"
import linuxproc "github.com/c9s/goprocinfo/linux"

type LoadAvgCollector struct { }

func (c LoadAvgCollector) Collect() (map[string]float64, error) {
    stat, err := linuxproc.ReadLoadAvg("/proc/loadavg")
    if err != nil {
        log.Fatal("stat read fail")
        return nil, err
    }

    // TODO: Add processes_running and processes_total,
    // unit:processes, type:(running|total)
    return map[string]float64{
        "01": stat.Last1Min,
        "05": stat.Last5Min,
        "15": stat.Last15Min,
    }, nil
}

func (c LoadAvgCollector) Report() {
    values, _ := c.Collect()

    if values != nil {
        for k, v := range values {
            log.WithFields(log.Fields{
                "target_type": "gauge",
                "type": k,
                "unit": "Load",
                "result": v,
            }).Info()
        }
    }
}
