package collectors

import "github.com/josegonzalez/metricsd/structs"

type CollectorInterface interface {
	Report() (structs.MetricSlice, error)
}
