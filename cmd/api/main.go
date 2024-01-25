package main

import (
	"clean/pkg/config"
	"clean/pkg/di"
	"log"

	"github.com/joho/godotenv"
	"github.com/swaggo/swag/example/basic/docs"
)

func main() {

	// // swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "ZEE-Cosmetics-E-commerce"
	docs.SwaggerInfo.Description = "ZEE-Cosmetics-E-commerce"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:3000"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http"}

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
