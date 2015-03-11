package main

import log "github.com/Sirupsen/logrus"

type CollectorInterface interface {
	Report() ([]log.Fields, error)
}
