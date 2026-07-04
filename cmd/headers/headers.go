package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")

func NewHeaders() Headers {
	return make(Headers)
}

func parseHeader(fieldLine []byte) (string, string, error) {
	fmt.Println("==== parseHeader ====")
	fmt.Println(string(fieldLine))
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field-line")
	}

	fieldName, fieldValue := parts[0], parts[1]
	if fieldName[0] == ' ' || fieldName[len(fieldName)-1] == ' ' {
		return "", "", fmt.Errorf("malformed field-name")
	}
	fieldValue = bytes.TrimSpace(fieldValue)
	return string(fieldName), string(fieldValue), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false

	for {
		idx := bytes.Index(data[read:], rn)
		// No "\r\n" found
		if idx == -1 {
			break
		}
		// "\r\n" found on start of stream
		if idx == 0 {
			done = true
			read += len(rn)
			break
		}

		fieldName, fieldValue, err := parseHeader(data[read : read+idx])
		if err != nil {
			return read, done, err
		}
		h[fieldName] = fieldValue
		read += idx + len(rn)
	}
	return read, done, nil
}
