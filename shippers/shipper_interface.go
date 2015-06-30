package shippers

import "github.com/josegonzalez/metricsd/structs"
import "github.com/vaughan0/go-ini"

type ShipperInterface interface {
	Setup(ini.File)
	Ship(structs.MetricSlice) error
}
