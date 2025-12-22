package request

import (
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

	startLine := before
	restOfMsg := after

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, restOfMsg, fmt.Errorf("invalid number of parts in request line: %s", startLine)
	}

	rL := &RequestLine{Method: parts[0], RequestTarget: parts[1], HttpVersion: parts[2]}

	return rL, restOfMsg, nil
}

// parses the request line from the reader
func RequestFromReader(reader io.Reader) (*Request, error)
