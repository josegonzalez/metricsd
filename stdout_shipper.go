package main

import "fmt"
import "github.com/vaughan0/go-ini"

type StdoutShipper struct{}

func (shipper *StdoutShipper) Setup(_ ini.File) {
}

func (shipper *StdoutShipper) Ship(logs MetricMapSlice) error {
	for _, item := range logs {
		serialized := MarshalData(item)
		fmt.Printf("%s\n", string(serialized))
	}

	return nil
}
