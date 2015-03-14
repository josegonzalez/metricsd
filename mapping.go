package main

type FloatMetricMap map[string]float64

type IntMetricMap map[string]uint64

type IndexMap map[string]string

type ActionMap map[string]IndexMap

type MetricMap map[string]interface{}

type MetricMapSlice []MetricMap
