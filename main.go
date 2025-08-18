package main

import (
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

	data := make([]byte, 8)
	for {
		n, err := f.Read(data)
		if err != nil {
			break
		}

		fmt.Printf("read: %s\n", string(data[:n]))
	}
}
