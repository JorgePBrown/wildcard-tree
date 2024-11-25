package tokenizer

type Token struct {
	T       TokenType
	Literal string
}

type TokenType string

const (
	EOF            TokenType = "EOF"
	LBRACE         TokenType = "LBRACE"         // {
	RBRACE         TokenType = "RBRACE"         // }
	LBRACKET       TokenType = "LBRACKET"       // [
	RBRACKET       TokenType = "RBRACKET"       // ]
	WILDCARD_OPEN  TokenType = "WILDCARD_OPEN"  // {{
	WILDCARD_CLOSE TokenType = "WILDCARD_CLOSE" // }}
	LPAREN         TokenType = "LPAREN"         // (
	RPAREN         TokenType = "RPAREN"         // )
	QUESTION_MARK  TokenType = "QUESTION_MARK"  // ?
	NULL_COALESCE  TokenType = "NULL_COALESCE"  // ??
	DOT            TokenType = "DOT"            // .
	PIPE           TokenType = "PIPE"           // |
	TEXT           TokenType = "TEXT"
	ILLEGAL        TokenType = "ILLEGAL"
)

func newToken(t TokenType, literal string) Token {
	return Token{
		T:       t,
		Literal: literal,
	}
}
