package main

import (
	routes "issue-service/routes"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	log.Printf("Server started")

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return
	}

	var config routes.EnvConfiguration

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Error reading env configuration")
		return
	}

	router := routes.NewRouter(config)

	log.Fatal(http.ListenAndServe(":8080", router))
}
