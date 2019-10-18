package main

import (
	"com/merkinsio/oasis-api/config"
	"com/merkinsio/oasis-api/logs"
	"com/merkinsio/oasis-api/server"

	"fmt"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&logs.LogFormatter{})
}

func main() {
	app := App{}

	fmt.Println("Starting API")
	app.Initialize()
	defer app.ShutDown()

	go server.HandleSocketMessages()

	fmt.Println("Server running at localhost:", config.Config.GetString("server.port"))
	app.Run(":" + config.Config.GetString("server.port"))
}
