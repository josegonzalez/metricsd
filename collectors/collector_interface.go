package collectors

import "github.com/josegonzalez/metricsd/mappings"

type CollectorInterface interface {
	Report() (mappings.MetricMapSlice, error)
}
