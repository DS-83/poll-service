package main

import (
	"log"
	"poll-service/config"
	"poll-service/server"

	"github.com/spf13/viper"
)

const appPort = "app_port"

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app := server.NewApp()

	if err := app.Run(viper.GetString(appPort)); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
