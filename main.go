package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

// getLinesReader returns a channel that will receive the lines of the file
// io.ReadCloser: interface for reading files and closing them
// os.File implements io.ReadCloser, to pass it directly to getLinesReader
func getLinesReader(f io.ReadCloser) <-chan string {
	ch := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(ch)

		str := ""
		// for file read operations, default buffer size is 4096 bytes
		buf := make([]byte, 8)
		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("read error: %v", err)
			}

			data := buf[:n]
			if i := bytes.IndexByte(data, '\n'); i >= 0 {
				str += string(data[:i])
				ch <- str
				str = ""
				data = data[i+1:] // skips the newline

			}
			str += string(data)
		}

		// if there is any data left in the buffer, prints it
		if len(str) > 0 {
			ch <- str
		}
	}()

	return ch
}

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
		for line := range getLinesReader(conn) {
			fmt.Println(line)
		}
	}
}
