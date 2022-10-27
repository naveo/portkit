package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/docker/docker/client"
)

// A DockerCli for initializing new Docker CLI client instant.
type DockerCli struct {
	cli *client.Client
}

// A Config represents the portkit config.json file structure.
type Config struct {
	NaveoSocket string `json:"naveoSocket"`
	NaveoHost   string `json:"naveoHost"`
	LogFilePath string `json:"logFilePath"`
	TargetHost  string `json:"targetHost"`
}

// ConfigFileParser parses the JSON-encoded data and stored
// in the portkit config.json file which hold operational data
// and file path for the naveo host, log file path, etc.
func ConfigFileParser(configFile string) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Printf("error parsing config file: %v", err)
	}
}
