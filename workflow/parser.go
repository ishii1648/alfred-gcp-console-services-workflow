package workflow

import (
	"bufio"
	"io"
)

type Parser struct {
	scanner *bufio.Scanner
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{scanner: bufio.NewScanner(reader)}
}

func (p *Parser) Parse() *Query {
	query := &Query{}
	count := 0

	p.scanner.Split(bufio.ScanWords)
	for p.scanner.Scan() {
		switch count {
		case 0:
			query.ServiceId = p.scanner.Text()
			count++
		case 1:
			query.SubServiceId = p.scanner.Text()
			count++
		default:
			break
		}
	}

	return query
}

type Query struct {
	ServiceId    string
	SubServiceId string
	// HasSearchFunc bool
}

func (q *Query) IsEmpty() bool {
	return q.ServiceId == "" && q.SubServiceId == ""
}
