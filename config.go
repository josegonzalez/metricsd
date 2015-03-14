package main

import "flag"
import "os"
import "github.com/vaughan0/go-ini"
import "github.com/Sirupsen/logrus"

var LogLevel string

func Setup() ini.File {
	configFile := flag.String("config", "/etc/metricsd/metricsd.ini", "full path to config file.")
    loglevel := flag.String("loglevel", "warning", "one of the following loglevels: [debug, info, warning, error, fatal, panic]")
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

	LogLevel = *loglevel

	return file
}
