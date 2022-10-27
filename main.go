package main

import (
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Panic(err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	DockerEvents()
}
