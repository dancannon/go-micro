package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/youtube/vitess/go/rpcplus"

	"github.com/asim/go-micro/errors"
	"github.com/asim/go-micro/transport"
)

type headerRoundTripper struct {
	r http.RoundTripper
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (t *headerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("X-Client-Version", "1.0")
	return t.r.RoundTrip(r)
}

func (t *Transport) call(address, path string, request *Request, response interface{}) error {
	pReq := &rpcplus.Request{
		ServiceMethod: request.Method(),
	}

	reqB := bytes.NewBuffer(nil)
	defer reqB.Reset()
	buf := &buffer{reqB, reqB}

	cc := t.marshaler.NewClientCodec(buf)
	err := cc.WriteRequest(pReq, request.Payload())
	if err != nil {
		return errors.InternalServerError("go.micro.client", fmt.Sprintf("Error writing request: %v", err))
	}

	client := &http.Client{}
	client.Transport = &headerRoundTripper{http.DefaultTransport}

	request.Headers().Set("Content-Type", t.marshaler.ContentType())

	hreq := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   address,
			Path:   path,
		},
		Header:        request.Headers().(http.Header),
		Body:          buf,
		ContentLength: int64(reqB.Len()),
		Host:          address,
	}

	rsp, err := client.Do(hreq)
	if err != nil {
		return errors.InternalServerError("go.micro.client", fmt.Sprintf("Error sending request: %v", err))
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return errors.InternalServerError("go.micro.client", fmt.Sprintf("Error reading response: %v", err))
	}

	rspB := bytes.NewBuffer(b)
	defer rspB.Reset()
	rBuf := &buffer{rspB, rspB}

	pRsp := &rpcplus.Response{}
	cc = t.marshaler.NewClientCodec(rBuf)
	err = cc.ReadResponseHeader(pRsp)
	if err != nil {
		return errors.InternalServerError("go.micro.client", fmt.Sprintf("Error reading response headers: %v", err))
	}

	if len(pRsp.Error) > 0 {
		return errors.Parse(pRsp.Error)
	}

	err = cc.ReadResponseBody(response)
	if err != nil {
		return errors.InternalServerError("go.micro.client", fmt.Sprintf("Error reading response body: %v", err))
	}

	return nil
}

func (t *Transport) Execute(request *Request, response interface{}) error {
	service, err := t.registry.GetService(request.Service())
	if err != nil {
		return errors.InternalServerError("go.micro.client", err.Error())
	}

	if len(service.Nodes) == 0 {
		return errors.NotFound("go.micro.client", "Service not found")
	}

	n := rand.Int() % len(service.Nodes)
	node := service.Nodes[n]
	address := fmt.Sprintf("%s:%d", node.Address, node.Port)

	return t.call(address, "/_rpc", request, response)
}

func (t *Transport) NewRequest(service, method string, payload interface{}) (transport.Request, error) {
	return &Request{
		service:   service,
		method:    method,
		payload:   payload,
		headers:   make(http.Header),
		transport: t,
	}, nil
}
