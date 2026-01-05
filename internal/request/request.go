package request

import (
	"fmt"
	"io"
	"strings"
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
func parseRequestLine(line string) (*RequestLine, int, error) {
	before, after, ok := strings.Cut(line, "\r\n")
	if !ok {
		return nil, line, fmt.Errorf("invalid request line: %s", line)
	}

	reqLine := before
	restOfMsg := after

	httparts := strings.Split(reqLine, " ") // method, request-target, HTTP-version
	if len(httparts) != 3 {
		return nil, restOfMsg, fmt.Errorf("invalid number of parts in request line: %s", reqLine)
	}

	if httparts[0] != "GET" && httparts[0] != "POST" && httparts[0] != "PUT" && httparts[0] != "DELETE" {
		return nil, restOfMsg, fmt.Errorf("invalid method: %s", httparts[0])
	}

	if httparts[1] == "" {
		return nil, restOfMsg, fmt.Errorf("empty request-target")
	}

	httpartsV := strings.Split(httparts[2], "/") // protocol, its version
	if len(httpartsV) != 2 || httpartsV[0] != "HTTP" || httpartsV[1] != "1.1" {
		return nil, restOfMsg, fmt.Errorf("invalid HTTP version: %s", httparts[2])
	}

	rL := &RequestLine{Method: httparts[0], RequestTarget: httparts[1], HttpVersion: httpartsV[1]}

	return rL, restOfMsg, nil
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
