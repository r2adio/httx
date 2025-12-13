package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer f.Close()

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
			fmt.Printf("read: %s\n", str)
			str = ""
			data = data[i+1:] // skips the newline

		}
		str += string(data)
	}

	// if there is any data left in the buffer, prints it
	if len(str) > 0 {
		fmt.Printf("read: %s\n", str)
	}
}
