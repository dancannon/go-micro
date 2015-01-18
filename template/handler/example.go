package handler

import (
	"code.google.com/p/go.net/context"
	"code.google.com/p/goprotobuf/proto"

	log "github.com/cihub/seelog"
	"github.com/dancannon/go-micro/server"
	example "github.com/dancannon/go-micro/template/proto/example"
)

type Example struct{}

func (e *Example) Call(ctx context.Context, req *example.Request, rsp *example.Response) error {
	log.Debug("Received Example.Call request")

	rsp.Msg = proto.String(server.Id + ": Hello " + req.GetName())

	return nil
}
