package tokenizer

import "testing"

func TestTokenizer(t *testing.T) {
	tt := []struct {
		input    string
		expected []Token
	}{
		{
			"{}", []Token{
				{T: LBRACE, Literal: "{"},
				{T: RBRACE, Literal: "}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			"{{", []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: EOF, Literal: ""},
			},
		},
		{
			"}}", []Token{
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			"{{}}", []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			"{{text}}", []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			"{{ text}}", []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{" text"}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: " text"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			"*", []Token{
				{T: ILLEGAL, Literal: "*"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{"text"}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{"text"?}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text"},
				{T: QUESTION_MARK, Literal: `?`},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{"text"??}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text"},
				{T: NULL_COALESCE, Literal: `??`},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{"text"|}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text"},
				{T: PIPE, Literal: `|`},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{"text|"|}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text|"},
				{T: PIPE, Literal: `|`},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{"text'a'"}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text'a'"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{'text"a"'}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: `text"a"`},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{'text'|}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text"},
				{T: PIPE, Literal: `|`},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{'text|'|}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "text|"},
				{T: PIPE, Literal: `|`},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{node.output[1]}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "node"},
				{T: DOT, Literal: "."},
				{T: TEXT, Literal: "output"},
				{T: LBRACKET, Literal: "["},
				{T: TEXT, Literal: "1"},
				{T: RBRACKET, Literal: "]"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{ node.output[1]}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: TEXT, Literal: "node"},
				{T: DOT, Literal: "."},
				{T: TEXT, Literal: "output"},
				{T: LBRACKET, Literal: "["},
				{T: TEXT, Literal: "1"},
				{T: RBRACKET, Literal: "]"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			`{{ ( node.output[1] )}}`, []Token{
				{T: WILDCARD_OPEN, Literal: "{{"},
				{T: LPAREN, Literal: "("},
				{T: TEXT, Literal: "node"},
				{T: DOT, Literal: "."},
				{T: TEXT, Literal: "output"},
				{T: LBRACKET, Literal: "["},
				{T: TEXT, Literal: "1"},
				{T: RBRACKET, Literal: "]"},
				{T: RPAREN, Literal: ")"},
				{T: WILDCARD_CLOSE, Literal: "}}"},
				{T: EOF, Literal: ""},
			},
		},
		{
			"", []Token{
				{T: EOF, Literal: ""},
			},
		},
	}

	for _, test := range tt {
		t.Log(test.input)
		tokenizer := New(test.input)
		for _, tok := range test.expected {
			actual := tokenizer.Next()
			if actual.T != tok.T {
				t.Errorf("wrong token type, expected=%s got=%s", tok.T, actual.T)
			}
			if actual.Literal != tok.Literal {
				t.Errorf("wrong token literal, expected=%s got=%s", tok.Literal, actual.Literal)
			}
			if t.Failed() {
				t.FailNow()
			}
		}
	}
}
