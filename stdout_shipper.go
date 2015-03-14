package main

import "fmt"

type StdoutShipper struct{}

func (hook *StdoutShipper) Ship(logs MetricMapSlice) error {
	for _, item := range logs {
		serialized := MarshalData(item)
		fmt.Printf("%s\n", string(serialized))
	}

	return nil
}
