package collectors

import "github.com/josegonzalez/metricsd/structs"
import "github.com/vaughan0/go-ini"

type CollectorInterface interface {
	Report() (structs.MetricSlice, error)
	Setup(ini.File)
}
