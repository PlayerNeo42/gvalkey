package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Parser is a RESP parser
type Parser struct {
	reader *bufio.Reader
}

// NewParser creates a new parser instance
func NewParser(rd io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(rd),
	}
}

// Parse is the entry point of the parser
// it returns the parsed value (any) and a potential error
func (p *Parser) Parse() (any, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	// RESP2
	case '*': // array
		return p.parseArray(line)
	case '$': // bulk string
		return p.parseBulkString(line)
	case '+': // simple string
		return p.parseSimpleString(line)
	case ':': // integer
		return p.parseInteger(line)
	// RESP3
	case '_', ',', '#', '!', '=', '(', '%', '~', '|', '>': // nil
		return nil, fmt.Errorf("RESP3 type not supported yet: %q", line)
	default:
		return nil, fmt.Errorf("unsupported RESP type: %q", line)
	}
}

// readLine reads a line (terminated by \r\n)
func (p *Parser) readLine() ([]byte, error) {
	line, err := p.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	// remove the trailing \r\n
	return line[:len(line)-2], nil
}

// parseArray parses an array
func (p *Parser) parseArray(line []byte) (Array, error) {
	// line example: *3
	count, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, fmt.Errorf("parse array length failed: %v", err)
	}

	// Redis's empty array or null array
	if count <= 0 {
		return Array{}, nil
	}

	result := make(Array, 0, count)
	for range count {
		// recursively call Parse to parse each element in the array
		// here we simplify the processing, assuming that the array elements are all Bulk String
		val, err := p.Parse()
		if err != nil {
			return nil, err
		}

		result = append(result, val)
	}

	return result, nil
}

// parseBulkString parses a bulk string
func (p *Parser) parseBulkString(line []byte) (BulkString, error) {
	// line example: $5
	length, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return "", fmt.Errorf("parse bulk string length failed: %v", err)
	}

	// Redis's null bulk string
	if length == -1 {
		return "", nil // return an empty string to represent nil
	}

	// read the string itself
	data := make([]byte, length+2) // +2 to read the trailing \r\n
	_, err = io.ReadFull(p.reader, data)
	if err != nil {
		return "", err
	}

	return BulkString(data[:length]), nil
}

func (p *Parser) parseSimpleString(line []byte) (SimpleString, error) {
	// line example: +OK
	return SimpleString(line[1:]), nil
}

func (p *Parser) parseInteger(line []byte) (Integer, error) {
	// line example: :100
	num, err := strconv.ParseInt(string(line[1:]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse integer failed: %v", err)
	}
	return Integer(num), nil
}
