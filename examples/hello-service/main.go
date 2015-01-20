package main

import (
	"log"

	"github.com/asim/go-micro"

	"github.com/asim/go-micro/examples/hello-service/handler"
)

func main() {
	config := micro.DefaultConsulConfig()
	service, err := micro.New("service.hello", config)
	if err != nil {
		log.Fatal(err)
	}

	server, err := service.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	server.Register(new(handler.Greeter))

	// Run server
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
