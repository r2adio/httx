package main

import (
	"fmt"
	"log"
	"net"

	"github.com/r2adio/httx/internal/request"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer l.Close()

	// lines := getLinesReader(f)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}
}
