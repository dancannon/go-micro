package rpc

import (
	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcplus/pbrpc"
)

type pbMarshaler struct{}

func (e jsonEncoder) ClientCodec(data interface{}) *rpcplus.ClientCodec {
	return pbrpc.NewClientCodec(data)
}

func (e jsonEncoder) ServerCodec(data interface{}) *rpcplus.ServerCodec {
	return pbrpc.NewServerCodec(data)
}
