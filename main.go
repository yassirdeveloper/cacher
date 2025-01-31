package main

import (
	"log"
	_ "net/http/pprof"
	"os"
	"strconv"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("CACHER_PORT"))
	if err != nil {
		log.Fatal("Error during reading CACHER_PORT variable from env: ", err)
	}

	nbrWorkers, err := strconv.Atoi(os.Getenv("CACHER_NBR_WORKERS"))
	if err != nil {
		log.Fatal("Error during reading CACHER_NBR_WORKERS variable from env: ", err)
	}

	logFilePath := "server.log"
	logPrefix := "- "
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error during creating the log file: ", err)
	}
	logger := NewLogger(logFile, logPrefix, log.Ldate|log.Ltime)

	server, err := InitServer(port, nbrWorkers, logger)
	if err != nil {
		log.Fatal("Error during init server: ", err)
	}
	server.Start()
}
