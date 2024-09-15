package main

import (
	_ "github.com/lib/pq"
	"log"
	"runners/config"
	"runners/server"
)

func main() {
	log.Println("Starting Runners App")
	log.Println("Initializing configuration")
	config, err := config.Init("runners")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initializing databese")
	dbHandler, err := server.InitDatabase(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initializing HTTP server")
	httpServer := server.InitHttp(config, dbHandler)
	httpServer.Start()
}
