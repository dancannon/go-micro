package transport

type Request interface {
	Service() string
	Method() string
	Payload() interface{}

	Headers() Headers

	Send() (interface{}, error)
}

type Headers interface {
	Add(string, string)
	Del(string)
	Get(string) string
	Set(string, string)
}