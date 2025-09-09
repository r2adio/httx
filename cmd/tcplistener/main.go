package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8)
			x, err := f.Read(data)
			if err != nil {
				break
			}
			data = data[:x] // take first 8 bytes present in data

			// IndexByte: returns the index of first instance of '\n' in data
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				out <- str // replacing the fmt.Printf
				str = ""
			}
			str += string(data)
		}

		// also reads the empty lines at the last of file
		if len(str) != 0 {
			out <- str // replacing the fmt.Printf
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		// waits for and return the next connection, to listener
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		// replace the the file a tcp connection as an argument to func
		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}
	}
}
