package http

type Request struct {
	requestLine RequestLine
	headers     *Headers
	body        []byte
}
