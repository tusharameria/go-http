package main

import (
	"bytes"
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
	s := ""
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
		data = data[:n]
		// Find index of a byte in a byte array
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			s += string(data[:i])
			fmt.Printf("read in : %s\n", s)
			data = data[i+1:]
			s = ""
		}
		s += string(data)
	}
	if len(s) != 0 {
		fmt.Printf("read : %s\n", s)
	}
}
