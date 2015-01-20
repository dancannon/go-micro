package main

import (
	"github.com/Sirupsen/logrus"
	"log"

	"github.com/asim/go-micro"

	"github.com/asim/go-micro/examples/hello-service/handler"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	config := micro.DefaultConsulConfig()
	config.RegistryAddr = "boot2docker:8500"
	config.StoreAddr = "boot2docker:8500"
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
	if err := server.ListenAndServe("127.0.0.1:8080"); err != nil {
		log.Fatal(err)
	}
}
