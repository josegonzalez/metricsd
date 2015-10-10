package shippers

import "github.com/josegonzalez/metricsd/structs"
import "github.com/vaughan0/go-ini"

// ShipperInterface is an exported type that
// defines the a generic interface for a shipper
type ShipperInterface interface {
	Enabled() bool
	Setup(ini.File)
	Ship(structs.MetricSlice) error
	State(bool)
}
