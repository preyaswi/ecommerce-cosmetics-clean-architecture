package main

import (
	"clean/pkg/config"
	"clean/pkg/di"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading the env file")
	}
	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config:", configErr)
	}
	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("can not start server:", diErr)
	} else {
		server.Start()
	}
}
