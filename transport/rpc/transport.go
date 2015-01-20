package rpc

import (
	"io"
	"sync"

	"github.com/youtube/vitess/go/rpcplus"

	"github.com/asim/go-micro/registry"
	"github.com/asim/go-micro/transport"
)

type marshaler interface {
	ContentType() string
	NewClientCodec(data io.ReadWriteCloser) rpcplus.ClientCodec
	NewServerCodec(data io.ReadWriteCloser) rpcplus.ServerCodec
}

type Transport struct {
	marshaler marshaler
	registry  registry.Registry

	// server related fields
	mtx  sync.RWMutex
	rpc  *rpcplus.Server
	exit chan chan error
}

func init() {
	transport.Register("pbrpc", &Transport{marshaler: pbMarshaler{}})
	transport.Register("jsonrpc", &Transport{marshaler: jsonMarshaler{}})
}

func (t *Transport) Init(reg registry.Registry) error {
	t.registry = reg
	t.rpc = rpcplus.NewServer()
	t.exit = make(chan chan error)

	return nil
}
