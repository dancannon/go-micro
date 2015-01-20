package rpc

import (
	"github.com/youtube/vitess/go/rpcplus"

	"github.com/asim/go-micro/registry"
)

type marshaler interface {
	NewClientCodec(data interface{}) rpcplus.ClientCodec
	NewServerCodec(data interface{}) rpcplus.ServerCodec
}

type Transport struct {
	marshaler marshaler
	registry  *registry.Registry

	// server related fields
	mtx sync.RWMutex
	rpc *rpc.Server
}

type Listener struct {
	exit chan chan error
}

func (l *Listener) Close() error {
	ch := make(chan error)
	l.exit <- ch
	return <-ch
}

func (t *RPCTransport) Init(reg *registry.Registry) error {
	t.registry = reg
}
