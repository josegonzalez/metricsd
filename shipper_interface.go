package main

type ShipperInterface interface {
	Ship(MetricMapSlice) error
}
