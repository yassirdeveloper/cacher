package main

import (
	"log"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"
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

	useSyncCache, err := strconv.ParseBool(os.Getenv("CACHER_USE_SYNC_CACHE"))
	if err != nil {
		log.Fatal("Error during reading CACHER_USE_SYNC_CACHE variable from env: ", err)
	}

	logFilePath := "server.log"
	logPrefix := "- "
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error during creating the log file: ", err)
	}
	logger := NewLogger(logFile, logPrefix, log.Ldate|log.Ltime)

	commandManager := NewCommandManager()
	cacheManager := NewCacheManager[string, string](logger)

	err = cacheManager.SetupMainCache(time.Minute)
	if err != nil {
		log.Fatal("Error setting up main cache: ", err)
	}

	err = cacheManager.SetupMainCacheJanitor(time.Minute * 5)
	if err != nil {
		log.Fatal("Error setting up main cache janitor: ", err)
	}

	if useSyncCache {
		err = cacheManager.SetupSyncCache(time.Minute * 5)
		if err != nil {
			log.Fatal("Error setting up sync cache: ", err)
		}

		err = cacheManager.SetupSyncCacheJanitor(time.Minute * 25)
		if err != nil {
			log.Fatal("Error setting up sync cache janitor: ", err)
		}
	}

	server, err := NewServer(port, nbrWorkers, logger, commandManager, cacheManager)
	if err != nil {
		log.Fatal("Error during init server: ", err)
		os.Exit(1)
	}
	server.Start(5 * time.Second)
}
