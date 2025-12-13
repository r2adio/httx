package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.ReadFile("message.txt")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("read: %s\n", string(f))
}
