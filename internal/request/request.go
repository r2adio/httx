package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

// request-line  = method SP request-target SP HTTP-version
func parseRequestLine(line string) (*RequestLine, string, error) {
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

// parses the request line from the reader
func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to read request line"), err)
	}

	str := string(data)
	rL, _, err := parseRequestLine(str)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to parse request line"), err)
	}

	return &Request{*rL}, err
}
