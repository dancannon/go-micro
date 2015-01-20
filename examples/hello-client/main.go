package main

import (
	"fmt"
	"log"

	"code.google.com/p/goprotobuf/proto"

	"github.com/asim/go-micro"

	greeterpb "github.com/asim/go-micro/examples/hello-service/proto/greeter"
)

func main() {
	config := micro.DefaultConsulConfig()
	config.RegistryAddr = "boot2docker:8500"
	config.StoreAddr = "boot2docker:8500"
	service, err := micro.New("client.hello", config)
	if err != nil {
		log.Fatal(err)
	}

	client, err := service.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	rsp := &greeterpb.Response{}
	req, err := client.NewRequest("service.hello", "Greeter.Hello", &greeterpb.Request{
		Name: proto.String("John"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// req.Headers().Set("X-User-Id", "john")
	// req.Headers().Set("X-From-Id", "script")

	// Call service
	if err := req.Execute(rsp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rsp.GetMsg())
}
