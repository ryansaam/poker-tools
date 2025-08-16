package cpokersjs2

import (
	"bufio"
	"strings"
)

type LineReader struct {
	scanner    *bufio.Scanner
	pushedBack *string
}

func NewLineReader(scanner *bufio.Scanner) *LineReader {
	// scanner.Buffer(make([]byte, 0, 64<<10), 1<<20) // enable if needed
	return &LineReader{scanner: scanner}
}

func (lineReader *LineReader) Next() (string, bool) {
	if lineReader.pushedBack != nil {
		str := *lineReader.pushedBack
		lineReader.pushedBack = nil
		return str, true
	}
	if !lineReader.scanner.Scan() {
		return "", false
	}
	return strings.TrimSpace(lineReader.scanner.Text()), true
}

func (lineReader *LineReader) Unread(line string) { str := line; lineReader.pushedBack = &str }
