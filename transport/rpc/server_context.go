package rpc

import (
	"net/http"
	"sync"

	log "github.com/cihub/seelog"

	"github.com/asim/go-micro/transport"
)

var ctxs = struct {
	sync.Mutex
	m map[*http.Request]*serverContext
}{
	m: make(map[*http.Request]*serverContext),
}

// A server context interface
type Context interface {
	Request() *http.Request     // the request made to the server
	Headers() transport.Headers // the response headers
}

// context represents the context of an in-flight HTTP request.
// It implements the appengine.Context and http.ResponseWriter interfaces.
type serverContext struct {
	req       *http.Request
	outCode   int
	outHeader http.Header
	outBody   []byte
}

// Copied from $GOROOT/src/pkg/net/http/transfer.go. Some response status
// codes do not permit a response body (nor response entity headers such as
// Content-Length, Content-Type, etc).
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}

func getServerContext(req *http.Request) *serverContext {
	ctxs.Lock()
	c := ctxs.m[req]
	ctxs.Unlock()

	if c == nil {
		// Someone passed in an http.Request that is not in-flight.
		panic("NewContext passed an unknown http.Request")
	}
	return c
}

// The response headers
func (c *serverContext) Headers() transport.Headers {
	return c.outHeader
}

// The response headers
func (c *serverContext) Header() http.Header {
	return c.outHeader
}

// The request made to the server
func (c *serverContext) Request() *http.Request {
	return c.req
}

func (c *serverContext) Write(b []byte) (int, error) {
	if c.outCode == 0 {
		c.WriteHeader(http.StatusOK)
	}
	if len(b) > 0 && !bodyAllowedForStatus(c.outCode) {
		return 0, http.ErrBodyNotAllowed
	}
	c.outBody = append(c.outBody, b...)
	return len(b), nil
}

func (c *serverContext) WriteHeader(code int) {
	if c.outCode != 0 {
		log.Errorf("WriteHeader called multiple times on request.")
		return
	}
	c.outCode = code
}

func GetContext(r *http.Request) *serverContext {
	return getServerContext(r)
}
