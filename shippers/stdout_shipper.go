package shippers

import "fmt"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/vaughan0/go-ini"

type StdoutShipper struct{
	enabled bool
}

func (shipper *StdoutShipper) Enabled() (bool) {
	return shipper.enabled
}

func (shipper *StdoutShipper) State(state bool) {
	shipper.enabled = state
}

func (shipper *StdoutShipper) Setup(_ ini.File) {
	shipper.State(true)
}

func (shipper *StdoutShipper) Ship(logs structs.MetricSlice) error {
	for _, item := range logs {
		serialized := item.ToJson()
		fmt.Printf("%s\n", string(serialized))
	}

	return nil
}
