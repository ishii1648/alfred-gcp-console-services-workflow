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
		count++
		switch count {
		case 1:
			query.ServiceId = p.scanner.Text()
		case 2:
			query.SubServiceId = p.scanner.Text()
		case 3:
			query.Filter = p.scanner.Text()
		default:
			break
		}
	}

	return query
}

type Query struct {
	ServiceId    string
	SubServiceId string
	Filter       string
	// HasSearchFunc bool
}

func (q *Query) IsEmpty() bool {
	return q.ServiceId == "" && q.SubServiceId == ""
}
