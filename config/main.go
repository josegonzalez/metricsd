package config

import "os"
import "github.com/ogier/pflag"
import "github.com/Sirupsen/logrus"
import "github.com/vaughan0/go-ini"

var LogLevel string

func Setup() ini.File {
	configFile := pflag.String("config", "/etc/metricsd/metricsd.ini", "full path to config file.")
	loglevel := pflag.String("loglevel", "warning", "one of the following loglevels: [debug, info, warning, error, fatal, panic]")
	pflag.Parse()

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
