package main

import "flag"
import "os"
import "github.com/vaughan0/go-ini"
import log "github.com/Sirupsen/logrus"

func Config() ini.File {
	configFile := flag.String("config", "", "full path to config file.")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("config file not specified")
	}

	if _, err := os.Stat(*configFile); err != nil {
		log.Fatal("config file does not exist")
	}

	file, err := ini.LoadFile(*configFile)
	if err != nil {
		log.Fatal("config file read failure")
	}

	return file
}
