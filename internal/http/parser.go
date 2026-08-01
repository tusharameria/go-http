package http

import (
	"bytes"
	"fmt"
	"strings"
)

func Parse(data []byte) (*Request, error) {
	requesLine, restOfData, found := bytes.Cut(data, []byte("\r\n"))
	if !found {
		return nil, fmt.Errorf("request line not ended")
	}

	rlParts := bytes.Split(requesLine, []byte(" "))
	if len(rlParts) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	headers := make(map[string]string)
	if len(restOfData) != 0 {
		headerParts := bytes.Split(restOfData, []byte("\r\n"))
		lenHeadParts := len(headerParts)
		if len(headerParts[lenHeadParts-1]) != 0 || len(headerParts[lenHeadParts-2]) != 0 {
			return nil, fmt.Errorf("headers not ended properly")
		}

		for i := 0; i < lenHeadParts-2; i++ {
			headerLine := headerParts[i]
			headerLineParts := bytes.SplitN(headerLine, []byte(":"), 2)
			if len(headerLineParts) != 2 {
				return nil, fmt.Errorf("invalid header line")
			}
			headers[strings.TrimSpace(string(headerLineParts[0]))] = strings.TrimSpace(string(headerLineParts[1]))
		}
	}

	return &Request{
		Method:  string(rlParts[0]),
		Path:    string(rlParts[1]),
		Version: string(rlParts[2]),
		Headers: headers,
	}, nil
}
