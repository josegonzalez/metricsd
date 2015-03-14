package main

type CollectorInterface interface {
	Report() (MetricMapSlice, error)
}
