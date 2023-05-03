package xfl

type Token int

const (
	//Special Token
	ILLEGAL Token = iota
	EOF

	//Literal
	HTML
	CONTENT
	STRING

	//Characters
	COLON
	COMMA
	SLASH
	GREATER_THEN
	EQUAL

	//Keyword
	START
	END
)

type Mode int

const (
	SEARCH_START Mode = iota
	SEARCH_HTML_END
	NORMAL
)
