package main

import "os"
import log "github.com/Sirupsen/logrus"

func prefixFieldClashes(data log.Fields) {
	_, ok := data["message"]
	if ok {
		data["fields.message"] = data["message"]
	}

	_, ok = data["level"]
	if ok {
		data["fields.level"] = data["level"]
	}
}

func MarshalData(entry *log.Entry) log.Fields {
	data := make(log.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		// Otherwise errors are ignored by `encoding/json`
		// https://github.com/Sirupsen/logrus/issues/137
		if err, ok := v.(error); ok {
			data[k] = err.Error()
		} else {
			data[k] = v
		}
	}
	prefixFieldClashes(data)
	data["@version"] = "1"
	data["@timestamp"] = entry.Time.Format("2006-01-02T15:04:05.000Z")
	if entry.Message != "" {
		data["message"] = entry.Message
		data["level"] = entry.Level.String()
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}

	data["host"] = hostname
	return data
}
