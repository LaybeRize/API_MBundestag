package xfl

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type HTMLItem struct {
	Attributes []string
	Content    string
}

type parseStruct struct {
	path         string
	latestString string
	t            string
	a            Token
	loop         bool
	s            *Scanner
	name         string
}

func ParseFile(filepath string, m *map[string]HTMLItem) {
	fi, err := os.Open(filepath)
	if err != nil {
		err = fi.Close()
		log.Fatal("Error while opening file ", err)
	}
	r := bufio.NewReader(fi)
	p := parseStruct{
		path:         filepath,
		latestString: "",
		t:            "",
		a:            0,
		loop:         true,
		s:            NewScanner(r),
	}
	for p.loop {
		p.startParsing(m)
	}
}

func (p *parseStruct) startParsing(m *map[string]HTMLItem) {
	p.searchStart()
	if !p.loop {
		return
	}
	p.getName()
	htmlItem := HTMLItem{Attributes: []string{}}
	p.getAttributes(&htmlItem)
	p.getHTML(&htmlItem)
	(*m)[p.name] = htmlItem
}

func (p *parseStruct) searchStart() {
	token, literal := p.s.Scan(SEARCH_START)
	switch token {
	case ILLEGAL:
		log.Fatal("Error while trying to parse new HTML Element malformed token '", literal, "'")
	case EOF:
		p.loop = false
		return
	case START:
		p.latestString = literal
	default:
		log.Fatal("Error: This statment should never be reached")
	}
}

func (p *parseStruct) getName() {
	token, literal := p.s.Scan(NORMAL)
	if token != CONTENT {
		log.Fatal("Error while trying to parse new HTML Element malformed token '", literal, "' after '", p.latestString, "'")
	}
	token, _ = p.s.Scan(NORMAL)
	if token != COLON {
		log.Fatal("Expected ':' after name in decleration of '", literal, "'")
	}
	p.name = literal
}

func (p *parseStruct) getAttributes(htmlItem *HTMLItem) {
	token := CONTENT
	var literal string
	for token != END {
		token, literal = p.s.Scan(NORMAL)
		if token == END {
			break
		}
		if token != CONTENT {
			log.Fatal("Expected variable name, got '", literal, "' in decleration of '", p.name, "' instead")
		}
		htmlItem.Attributes = append(htmlItem.Attributes, literal)
		token, literal = p.s.Scan(NORMAL)
		if token != COMMA && token != END {
			log.Fatal("Expected ',' or '-->' got '", literal, "' in decleration of '", p.name, "' instead")
		}
	}
}

func (p *parseStruct) getHTML(htmlItem *HTMLItem) {
	token, literal := p.s.Scan(SEARCH_HTML_END)
	if token != HTML {
		log.Fatal("Expected HTML part got '", literal, "' in decleration of '", p.name, "' instead")
	}
	htmlItem.Content = literal
}

func ParseHTML(errorIn string, html string) map[string]string {
	s := NewScannerFromString(html)
	s.skipToFirstWhiteSpace()
	token, literal := s.ScanHTMLPart()
	m := map[string]string{}
	for token == CONTENT {
		name := literal
		token, literal = s.ScanHTMLPart()
		if token != EQUAL {
			log.Fatal("Error in '" + errorIn + "' could not find '=' after attribute '" + name + "'")
		}
		token, literal = s.ScanHTMLPart()
		if token != STRING {
			log.Fatal("Error in '" + errorIn + "' could not find string assigned attribute '" + name + "'")
		}
		m[name] = literal
		token, literal = s.ScanHTMLPart()
	}
	if token == SLASH {
		return m
	}
	if token == GREATER_THEN {
		m["content"] = strings.TrimSpace(s.scanRest())
		return m
	}
	log.Fatal("Error in '" + errorIn + "' could not find correct closing got '" + literal + "' instead")
	return m
}
