package parser

import (
	"bytes"
)

type Parser struct {
	aggregator []byte
}

func (p *Parser) Feed(data []byte) [][]byte {
	res := [][]byte{}
	p.aggregator = append(p.aggregator, data...)
	for {
		firstIdx := bytes.IndexByte(p.aggregator, '\n')
		if firstIdx >= 0 {
			msg := p.aggregator[:firstIdx]
			res = append(res, msg)
		} else {
			break
		}
		p.aggregator = p.aggregator[firstIdx+1:]
	}
	return res
}
