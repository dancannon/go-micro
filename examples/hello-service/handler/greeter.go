package handler

import (
	"code.google.com/p/go.net/context"
	"code.google.com/p/goprotobuf/proto"
	log "github.com/cihub/seelog"

	greeterpb "github.com/asim/go-micro/examples/hello-service/proto/greeter"
)

type Greeter struct{}

func (e *Greeter) Hello(ctx context.Context, req *greeterpb.Request, rsp *greeterpb.Response) error {
	log.Debug("Received Example.Call request")

	rsp.Msg = proto.String("Hello " + req.GetName())

	return nil
}
