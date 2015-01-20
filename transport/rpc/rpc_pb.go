package rpc

import (
	"io"

	"github.com/youtube/vitess/go/rpcplus"
	"github.com/youtube/vitess/go/rpcplus/pbrpc"
)

type pbMarshaler struct{}

func (m pbMarshaler) ContentType() string {
	return "application/octet-stream"
}

func (m pbMarshaler) NewClientCodec(data io.ReadWriteCloser) rpcplus.ClientCodec {
	return pbrpc.NewClientCodec(data)
}

func (m pbMarshaler) NewServerCodec(data io.ReadWriteCloser) rpcplus.ServerCodec {
	return pbrpc.NewServerCodec(data)
}
