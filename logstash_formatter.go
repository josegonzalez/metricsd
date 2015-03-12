package main

import "encoding/json"
import "fmt"
import log "github.com/Sirupsen/logrus"

type LogstashFormatter struct{}

func (f *LogstashFormatter) Format(entry *log.Entry) ([]byte, error) {
	data := MarshalData(entry)

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
