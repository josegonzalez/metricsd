package main

import "flag"
import "os"
import "github.com/vaughan0/go-ini"
import "github.com/Sirupsen/logrus"

func Config() ini.File {
	configFile := flag.String("config", "", "full path to config file.")
	flag.Parse()

	if *configFile == "" {
		logrus.Fatal("config file not specified")
	}

	if _, err := os.Stat(*configFile); err != nil {
		logrus.Fatal("config file does not exist")
	}

	file, err := ini.LoadFile(*configFile)
	if err != nil {
		logrus.Fatal("config file read failure")
	}

	return file
}
