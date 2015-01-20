package rpc

import (
	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcplus/jsonrpc"
)

type jsonMarshaler struct{}

func (e jsonMarshaler) ClientCodec(data interface{}) *rpcplus.ClientCodec {
	return jsonrpc.NewClientCodec(data)
}

func (e jsonMarshaler) ServerCodec(data interface{}) *rpcplus.ServerCodec {
	return jsonrpc.NewServerCodec(data)
}
