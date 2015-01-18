package main

import (
	"log"

	"github.com/dancannon/go-micro/server"
	"github.com/dancannon/go-micro/template/handler"
)

func main() {
	server.Name = "go.micro.service.template"

	// Initialise Server
	server.Init()

	// Register Handlers
	server.Register(
		server.NewReceiver(
			new(handler.Example),
		),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
