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
func parseRequestLine(line string) (*RequestLine, error) {
	indx := strings.Index(line, "\r\n")
	if indx == -1 {
		return nil, nil
	}

	startLine := line[:indx]

	parts := strings.Split(startLine, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid number of parts in request line: %s", startLine)
	}

	return &RequestLine{
		Method: parts[0], RequestTarget: parts[1], HttpVersion: parts[2],
	}, nil
}

// parses the request line from the reader
func RequestFromReader(reader io.Reader) (*Request, error)
