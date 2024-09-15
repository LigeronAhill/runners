package main

import (
	"log"
	"runners/config"
)

func main() {
	config, err := config.Init("runners")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config)
}
