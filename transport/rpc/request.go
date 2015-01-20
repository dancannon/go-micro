package rpc

import (
	"net/http"

	"github.com/asim/go-micro/transport"
)

type Request struct {
	transport       *Transport
	service, method string
	payload         interface{}
	headers         http.Header
}

func (r *Request) Headers() transport.Headers {
	return r.headers
}

func (r *Request) Service() string {
	return r.service
}

func (r *Request) Method() string {
	return r.method
}

func (r *Request) Payload() interface{} {
	return r.payload
}

func (r *Request) Execute(response interface{}) error {
	return r.transport.Execute(r, response)
}
