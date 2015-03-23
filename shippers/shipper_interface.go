package shippers

import "github.com/josegonzalez/metricsd/structs"
import "github.com/vaughan0/go-ini"

type ShipperInterface interface {
	Ship(structs.MetricSlice) error
	Setup(ini.File)
}
