package headers

import (
	"bytes"
	"fmt"
	"strings"
)

func isToken(str []byte) bool {
	for _, ch := range str {
		if ch >= 'A' && ch <= 'Z' ||
			ch >= 'a' && ch <= 'z' ||
			ch >= '0' && ch <= '9' {
			continue
		}
		switch ch {
		case '!', '#', '$', '%', '&', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		}
		return false
	}
	return true
}

type Headers struct {
	header map[string]string
}

var rn = []byte("\r\n")

func NewHeaders() *Headers {
	return &Headers{
		header: make(map[string]string),
	}
}

func (h *Headers) Get(name string) string {
	return h.header[strings.ToLower(name)]
}

func (h *Headers) Set(name, value string) {
	h.header[strings.ToLower(name)] = value
}

func parseHeader(fieldLine []byte) (string, string, error) {
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

func (h *Headers) Parse(data []byte) (int, bool, error) {
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
			return 0, false, err
		}

		if !isToken([]byte(fieldName)) {
			return 0, false, fmt.Errorf("malformed header name")
		}

		h.Set(fieldName, fieldValue)
		read += idx + len(rn)
	}
	return read, done, nil
}
