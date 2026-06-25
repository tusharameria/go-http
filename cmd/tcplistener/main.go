package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadWriteCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)
		// Read the file till it reaches the end
		s := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
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
				out <- s
				data = data[i+1:]
				s = ""
			}
			s += string(data)
		}
		if len(s) != 0 {
			out <- s
		}
	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "err", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error in Accept", "err", err)
		}
		for line := range getLinesChannel(conn) {
			fmt.Printf("read : %s\n", line)
		}
	}
}
