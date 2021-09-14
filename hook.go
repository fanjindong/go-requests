package requests

type Hook interface {
	// BeforeProcess Before the HTTP request is executed
	BeforeProcess(req *Request)
	// AfterProcess After the HTTP request is executed
	AfterProcess(req *Request, resp *Response, err error)
}
