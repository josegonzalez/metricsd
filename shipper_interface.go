package main

import "github.com/vaughan0/go-ini"

type ShipperInterface interface {
	Ship(MetricMapSlice) error
	Setup(ini.File)
}
