package xfl

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var eof = rune(0)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

type Scanner struct {
	r       *bufio.Reader
	pos     int
	str     []rune
	current rune
	next    rune
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	scan := Scanner{r: bufio.NewReader(r)}
	scan.current = scan.readBuffer()
	scan.next = scan.readBuffer()
	return &scan
}

// readBuffer reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) readBuffer() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) readNextRune() rune {
	s.pos += 1
	if s.pos >= len(s.str) {
		return eof
	} else {
		return s.str[s.pos]
	}
}

func (s *Scanner) consume() rune {
	res := s.current
	s.current = s.next
	if s.r == nil {
		s.next = s.readNextRune()
	} else {
		s.next = s.readBuffer()
	}
	return res
}

func (s *Scanner) fillBuffer() {
	s.consume()
	s.consume()
}

func (s *Scanner) getCurrent() rune {
	return s.current
}

func (s *Scanner) peak() rune {
	return s.next
}

// unread places the previously readBuffer rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) Scan(m Mode) (tok Token, lit string) {
	// Read the next rune.
	ch := s.getCurrent()

	if m == SEARCH_START {
		for !(ch == '<' && s.peak() == '!') {
			ch = s.consume()
			if ch == eof {
				return EOF, ""
			}
		}
		s.fillBuffer()
		if s.getCurrent() == '-' && s.peak() == '-' {
			s.fillBuffer()
			return START, "<!--"
		} else {
			return ILLEGAL, "<!" + string(s.getCurrent())
		}
	}

	if m == SEARCH_HTML_END {
		var buf bytes.Buffer
		for !(ch == '<' && s.getCurrent() == '!') {
			if ch == eof {
				break
			}
			buf.WriteRune(ch)
			ch = s.consume()
		}
		if ch == '<' && s.getCurrent() == '!' {
			s.unread()
			s.current = '<'
			s.next = '!'
		}
		return HTML, strings.TrimSpace(buf.String())
	}

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.scanWhitespace()
		ch = s.getCurrent()
	}
	if isLetter(ch) {
		return s.scanIdent()
	}

	switch true {
	case ch == eof:
		return EOF, ""
	case ch == ':':
		s.consume()
		return COLON, ":"
	case ch == ',':
		s.consume()
		return COMMA, ","
	case ch == '-' && s.peak() == '-':
		s.fillBuffer()
		if s.getCurrent() == '>' {
			s.consume()
			return END, "-->"
		} else {
			return ILLEGAL, "--" + string(ch)
		}
	}

	// Otherwise readBuffer the individual character.

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanWhitespace() {
	// Non-whitespace characters and EOF will cause the loop to exit.
	s.consume()
	for {
		if ch := s.getCurrent(); !isWhitespace(ch) {
			break
		}
		s.consume()
	}
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and readBuffer the current character into it.
	var buf bytes.Buffer

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	var ch rune
	for {
		buf.WriteRune(s.consume())
		if ch = s.getCurrent(); !isLetter(ch) && !isDigit(ch) && !isAllowedInName(ch) {
			break
		}
	}

	// Otherwise return as a regular identifier.
	return CONTENT, buf.String()
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isAllowedInName(ch rune) bool {
	return ch == '_' || ch == '-'
}

func NewScannerFromString(r string) *Scanner {
	scan := Scanner{str: []rune(r), pos: 1, r: nil}
	scan.current = scan.str[0]
	scan.next = scan.str[1]
	return &scan
}

func (s *Scanner) ScanHTMLPart() (tok Token, lit string) {
	// Read the next rune.
	ch := s.getCurrent()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.scanWhitespace()
		ch = s.getCurrent()
	}
	if isLetter(ch) {
		return s.scanIdent()
	}
	if ch == '"' {
		return s.scanString()
	}

	switch true {
	case ch == eof:
		return EOF, ""
	case ch == '/':
		s.consume()
		return SLASH, "/"
	case ch == '>':
		s.consume()
		return GREATER_THEN, ">"
	case ch == '=':
		s.consume()
		return EQUAL, "="
	}

	// Otherwise readBuffer the individual character.

	return ILLEGAL, string(ch)
}

func (s *Scanner) scanString() (tok Token, lit string) {
	// Create a buffer and readBuffer the current character into it.
	var buf bytes.Buffer

	// Read every subsequent character into the buffer as long as they are not a ".
	s.consume()
	for {
		if s.getCurrent() == eof {
			return ILLEGAL, buf.String()
		} else if s.getCurrent() != '"' {
			buf.WriteRune(s.consume())
		} else {
			break
		}
	}
	s.consume()

	// Otherwise return as a regular identifier.
	return STRING, buf.String()
}

func (s *Scanner) skipToFirstWhiteSpace() {
	for !isWhitespace(s.getCurrent()) && s.getCurrent() != '/' && s.getCurrent() != '>' {
		s.consume()
	}
}

func (s *Scanner) scanRest() string {
	var buf bytes.Buffer
	for {
		if s.getCurrent() == eof {
			break
		}
		buf.WriteRune(s.consume())
	}
	if buf.Len() == 0 {
		return ""
	}
	str := s.removeLastPart(buf.String())
	return str
}

func (s *Scanner) removeLastPart(str string) string {
	runes := []rune(str)
	last := len(runes) - 1
	for runes[last] != '<' {
		last--
		runes = runes[:last+1]
	}
	runes = runes[:last]
	return string(runes)
}
