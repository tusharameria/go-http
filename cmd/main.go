package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// Open the file for reading
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	// Always close the file when we're done using defer
	defer file.Close()

	// Read the file till it reaches the end
	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		if err != nil {
			// Always check for End of File while reading a file to properly close it
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Printf("read : %s\n", string(data[:n]))
	}
}
