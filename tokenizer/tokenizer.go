package tokenizer

import (
	"bytes"
	"unicode"
)

type Tokenizer struct {
	input        string
	peekPosition int
	position     int
	current      byte
}

func New(input string) *Tokenizer {
	return &Tokenizer{
		input: input,
	}
}

func (t *Tokenizer) Next() Token {
	ch := t.read()

	switch ch {
	case '{':
		if t.expect('{') {
			return newToken(WILDCARD_OPEN, "{{")
		} else {
			return newToken(LBRACE, string(ch))
		}
	case '}':
		if t.expect('}') {
			return newToken(WILDCARD_CLOSE, "}}")
		} else {
			return newToken(RBRACE, string(ch))
		}
	case '\'':
		if c := t.read(); c == 0 {
			return newToken(EOF, "")
		}
		return newToken(TEXT, t.readWord('\''))
	case '"':
		if c := t.read(); c == 0 {
			return newToken(EOF, "")
		}
		return newToken(TEXT, t.readWord('"'))
	case '?':
		if t.expect('?') {
			return newToken(NULL_COALESCE, "??")
		} else {
			return newToken(QUESTION_MARK, "?")
		}
	case '|':
		return newToken(PIPE, string(ch))
	case '[':
		return newToken(LBRACKET, string(ch))
	case ']':
		return newToken(RBRACKET, string(ch))
	case '(':
		return newToken(LPAREN, string(ch))
	case ')':
		return newToken(RPAREN, string(ch))
	case '.':
		return newToken(DOT, string(ch))
	case ' ':
		return t.Next()
	case 0:
		return newToken(EOF, "")
	default:
		if t.isLetter(ch) || t.isNumber(ch) {
			word := t.readWord(0)
			return newToken(TEXT, word)
		}
		return newToken(ILLEGAL, string(ch))
	}
}

func (t *Tokenizer) peek() byte {
	if t.peekPosition < len(t.input) {
		return t.input[t.peekPosition]
	}
	return 0
}

func (t *Tokenizer) expect(b byte) bool {
	if ch := t.peek(); ch == b {
		t.position = t.peekPosition
		t.peekPosition += 1
		return true
	}
	return false
}

func (t *Tokenizer) isLetter(b byte) bool {
	return unicode.IsLetter(rune(b))
}

func (t *Tokenizer) isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}

func (t *Tokenizer) readWord(terminator byte) string {
	var out bytes.Buffer
	out.WriteByte(t.current)

	ch := t.peek()
	if terminator != 0 {
		for ch != 0 && ch != terminator {
			out.WriteByte(t.read())
			ch = t.peek()
		}
		if ch != 0 {
			t.read()
		}
	} else {
		for t.isLetter(ch) || t.isNumber(ch) {
			out.WriteByte(t.read())
			ch = t.peek()
		}
	}

	return out.String()
}

func (t *Tokenizer) read() byte {
	t.position = t.peekPosition
	t.peekPosition += 1

	if t.position < len(t.input) {
		t.current = t.input[t.position]
		return t.current
	} else {
		return 0
	}
}
