package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"

	log "github.com/cihub/seelog"

	"github.com/asim/go-micro/errors"
)

var (
	HealthPath = "/_status/health"
	RPCPath    = "/_rpc"
)

func executeRequestSafely(c *serverContext, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			log.Criticalf("Panicked on request: %v", r)
			log.Criticalf("%v: %v", x, string(debug.Stack()))
			err := errors.InternalServerError("go.micro.server", "Unexpected error")
			c.WriteHeader(500)
			c.Write([]byte(err.Error()))
		}
	}()

	http.DefaultServeMux.ServeHTTP(c, r)
}

func (t *Transport) handler(w http.ResponseWriter, r *http.Request) {
	c := &serverContext{
		req:       r,
		outHeader: w.Header(),
	}

	ctxs.Lock()
	ctxs.m[r] = c
	ctxs.Unlock()
	defer func() {
		ctxs.Lock()
		delete(ctxs.m, r)
		ctxs.Unlock()
	}()

	// Patch up RemoteAddr so it looks reasonable.
	if addr := r.Header.Get("X-Forwarded-For"); len(addr) > 0 {
		r.RemoteAddr = addr
	} else {
		// Should not normally reach here, but pick a sensible default anyway.
		r.RemoteAddr = "127.0.0.1"
	}
	// The address in the headers will most likely be of these forms:
	//	123.123.123.123
	//	2001:db8::1
	// net/http.Request.RemoteAddr is specified to be in "IP:port" form.
	if _, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		// Assume the remote address is only a host; add a default port.
		r.RemoteAddr = net.JoinHostPort(r.RemoteAddr, "80")
	}

	executeRequestSafely(c, r)
	c.outHeader = nil // make sure header changes aren't respected any more

	// Avoid nil Write call if c.Write is never called.
	if c.outCode != 0 {
		w.WriteHeader(c.outCode)
	}
	if c.outBody != nil {
		w.Write(c.outBody)
	}
}

func (t *Transport) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	serveCtx := getServerContext(req)

	// TODO: get user scope from context
	// check access

	if req.Method != "POST" {
		err := errors.BadRequest("go.micro.server", "Method not allowed")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	defer req.Body.Close()

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errr := errors.InternalServerError("go.micro.server", fmt.Sprintf("Error reading request body: %v", err))
		w.WriteHeader(500)
		w.Write([]byte(errr.Error()))
		log.Errorf("Erroring reading request body: %v", err)
		return
	}

	rbq := bytes.NewBuffer(b)
	rsp := bytes.NewBuffer(nil)
	defer rsp.Reset()
	defer rbq.Reset()

	buf := &buffer{
		rbq,
		rsp,
	}

	cc := t.marshaler.NewServerCodec(buf)
	ctx := newContext(&ctx{}, serveCtx)
	err = t.rpc.ServeRequestWithContext(ctx, cc)
	if err != nil {
		// This should not be possible.
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Errorf("Erroring serving request: %v", err)
		return
	}

	w.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", strconv.Itoa(rsp.Len()))
	w.Write(rsp.Bytes())
}

func (t *Transport) Register(handler interface{}) {
	t.rpc.Register(handler)
}

func (t *Transport) RegisterNamed(name string, handler interface{}) {
	t.rpc.RegisterName(name, handler)
}

func (t *Transport) StartListening(addr string) error {
	registerHealthChecker(http.DefaultServeMux)
	http.Handle(RPCPath, t)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Debugf("Listening on %s", l.Addr().String())

	t.mtx.Lock()
	t.mtx.Unlock()

	go http.Serve(l, http.HandlerFunc(t.handler))

	go func() {
		ch := <-t.exit
		ch <- l.Close()
	}()

	return nil
}

func (t *Transport) StopListening() error {
	ch := make(chan error)
	t.exit <- ch
	return <-ch
}
