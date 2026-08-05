package http

type Request struct {
	RequestLine RequestLine
	Headers     *Headers
}
