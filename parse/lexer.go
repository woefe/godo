package parse

import (
	"bufio"
	"bytes"
	"io"
)

//Token identifies the type of data that was read
type Token int

const (
	//Special
	ILLEGAL = iota //ILLEGAL represents anything that cannot be identified by any other token
	EOF            //EOF represents the end of file token
	WS             //WS identifies a whitespace

	//Literals
	IDENT //todos and dates

	//Key Symbols
	STATUS_OPEN  //[
	STATUS_CLOSE //]

	//Misc
	SLASH         // /
	SEMICOLON     // ;
	COLON         // :
	ASTERISK      // *
	COMMA         // ,
	DOT           // .
	HASHTAG       // #
	BRACKET       // ( )
	CURRENCY_SIGN // $ €
	PARAGRAPH     // §
	AMPERSAND     // &
	EQUALS        // =
	TILDE         // ~
	AT            // @
	PERCENT       // %
	DASH          // -
	UNDERSCORE    // _
)

var eof = rune(0)

//Scanner represents a lexical scanner
type scanner struct {
	*bufio.Reader
}

//NewScanner returns a new instance of Scanner
func NewScanner(r io.Reader) *scanner {
	return &scanner{bufio.NewReader(r)}
}

//read reads the next rune from the buffered reader.
//Returns the rune(0) if an error occurs(or io.EOF is returned).
func (s *scanner) read() rune {
	r, _, err := s.ReadRune()
	if err != nil {
		return eof
	}

	return r
}

//Scan returns the next token and its value
func (s *scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.UnreadRune()
		return s.scanWhitespace()
	} else if isLetter(ch) || isDigit(ch) {
		s.UnreadRune()
		return s.scanIdent()
	}

	//Otherwise read individual character
	switch ch {
	case '#':
		return HASHTAG, "#"
	case '[':
		return STATUS_OPEN, "["
	case ']':
		return STATUS_CLOSE, "]"
	case ',':
		return COMMA, ","
	case '.':
		return DOT, "."
	case ':':
		return COLON, ":"
	case ';':
		return SEMICOLON, ";"
	case '/':
		return SLASH, "/"
	case '*':
		return ASTERISK, "*"
	case '(':
		fallthrough
	case ')':
		return BRACKET, string(ch)
	case '~':
		return TILDE, "~"
	case '€':
		fallthrough
	case '$':
		fallthrough
	case '£':
		fallthrough
	case '¥':
		return CURRENCY_SIGN, string(ch)
	case '§':
		return PARAGRAPH, "§"
	case '&':
		return AMPERSAND, "&"
	case '=':
		return EQUALS, "="
	case '@':
		return AT, "@"
	case '%':
		return PERCENT, "%"
	case '-':
		return DASH, "-"
	case '_':
		return UNDERSCORE, "_"
	case eof:
		return EOF, string(ch)
	}

	return ILLEGAL, string(ch)
}

func (s *scanner) scanWhitespace() (tok Token, lit string) {
	//Create buffer and read the current character into it
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	//Read every subsequent whitespace into the buffer.
	//Non Whitspace Characters and EOF will cause the loop to exit
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.UnreadRune()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *scanner) scanIdent() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	//Read every subsequent ident character into the buffer.
	//Non ident Characters and EOF will cause the loop to exit

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.UnreadRune()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return IDENT, buf.String()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == 'ä' || ch == 'Ö' ||
		ch == 'ö' || ch == 'Ä' ||
		ch == 'ü' || ch == 'Ü' ||
		ch == 'ß'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
