package rpc

type Request struct {
	transport         *RPCTransport
	service, endpoint string
	payload           interface{}
	headers           http.Header
}

func (r *Request) Headers() Headers {
	return r.headers
}

func (r *Request) Service() string {
	return r.service
}

func (r *Request) Endpoint() string {
	return r.method
}

func (r *Request) Payload() interface{} {
	return r.payload
}

func (r *Request) Send(response interface{}) error {
	return r.transport.Send(r, response)
}
