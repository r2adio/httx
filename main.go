package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error opening file:", err)
	}

	defer f.Close()

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
			fmt.Printf("read: %s\n", str)
			str = ""
		}
		str += string(data)
	}

	// also reads the empty lines at the last of file
	if len(str) != 0 {
		fmt.Printf("read: %s\n", str)
	}
}
