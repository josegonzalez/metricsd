package shippers

import "fmt"
import "github.com/josegonzalez/metricsd/structs"
import "github.com/vaughan0/go-ini"

// StdoutShipper is an exported type that
// allows shipping metrics to stdout
type StdoutShipper struct {
	enabled bool
}

// Enabled allows checking whether the shipper is enabled or not
func (s *StdoutShipper) Enabled() bool {
	return s.enabled
}

// State allows setting the enabled state of the shipper
func (s *StdoutShipper) State(state bool) {
	s.enabled = state
}

// Setup configures the shipper
func (s *StdoutShipper) Setup(_ ini.File) {
	s.State(true)
}

// Ship sends a list of MetricSlices to stdout
func (s *StdoutShipper) Ship(logs structs.MetricSlice) error {
	for _, item := range logs {
		serialized := item.ToJSON()
		fmt.Printf("%s\n", string(serialized))
	}

	return nil
}
