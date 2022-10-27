package main

import (
	"flag"
	"os"
)

var config Config

func init() {
	configFilePath := flag.String("config", "", "config file path")
	flag.Parse()
	if *configFilePath == "" {
		flag.Usage()
		os.Exit(1)
	}

	configFile, err := os.Open(*configFilePath)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	ConfigFileParser(configFile.Name())
}
