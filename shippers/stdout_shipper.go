package shippers

import "fmt"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/vaughan0/go-ini"

type StdoutShipper struct{
	enabled bool
}

func (this *StdoutShipper) Enabled() (bool) {
	return this.enabled
}

func (this *StdoutShipper) State(state bool) {
	this.enabled = state
}

func (this *StdoutShipper) Setup(_ ini.File) {
	this.State(true)
}

func (this *StdoutShipper) Ship(logs structs.MetricSlice) error {
	for _, item := range logs {
		serialized := item.ToJson()
		fmt.Printf("%s\n", string(serialized))
	}

	return nil
}
