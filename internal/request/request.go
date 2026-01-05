package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{state: StateInit}
}

// request-line = method SP request-target SP HTTP-version
func parseRequestLine(buf []byte) (*RequestLine, int, error) {
	idx := bytes.Index(buf, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, fmt.Errorf("invalid request line: %s", buf)
	}
	reqLine := buf[:idx]
	read := idx + 2 // +2 for \r\n

	httparts := bytes.Split(reqLine, []byte(" ")) // method, request-target, HTTP-version
	if len(httparts) != 3 {
		return nil, 0, fmt.Errorf("invalid number of parts in request line: %s", reqLine)
	}

	if string(httparts[0]) != "GET" && string(httparts[0]) != "POST" && string(httparts[0]) != "PUT" && string(httparts[0]) != "DELETE" {
		return nil, 0, fmt.Errorf("invalid method: %s", httparts[0])
	}

	if httparts[1] == nil {
		return nil, 0, fmt.Errorf("empty request-target")
	}

	httpartsV := bytes.Split(httparts[2], []byte("/")) // protocol, its version
	if len(httpartsV) != 2 || string(httpartsV[0]) != "HTTP" || string(httpartsV[1]) != "1.1" {
		return nil, 0, fmt.Errorf("invalid HTTP version: %s", httparts[2])
	}

	rL := &RequestLine{Method: string(httparts[0]), RequestTarget: string(httparts[1]), HttpVersion: string(httpartsV[1])}

	return rL, read, nil
}

func (r *Request) parse(buf []byte) (int, error) {
	return 0, nil
}

func (r *Request) done() bool { return r.state == StateDone }

// parses the request line from the reader
func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// NOTE: header could overflow the buffer(1024 bytes)
	// header or body could be more than 1k
	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err // TODO: handle error
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen]) // shifts the remaining data to the beginning of the buffer
		bufLen -= readN              // updates the buffer length
	}

	return request, nil
}
