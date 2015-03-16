package shippers

import "github.com/josegonzalez/metricsd/mappings"
import "github.com/vaughan0/go-ini"

type ShipperInterface interface {
	Ship(mappings.MetricMapSlice) error
	Setup(ini.File)
}
