package rpc

import (
	"io"

	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcplus/jsonrpc"
)

type jsonMarshaler struct{}

func (m jsonMarshaler) ContentType() string {
	return "application/json"
}

func (m jsonMarshaler) NewClientCodec(data io.ReadWriteCloser) rpcplus.ClientCodec {
	return jsonrpc.NewClientCodec(data)
}

func (m jsonMarshaler) NewServerCodec(data io.ReadWriteCloser) rpcplus.ServerCodec {
	return jsonrpc.NewServerCodec(data)
}
