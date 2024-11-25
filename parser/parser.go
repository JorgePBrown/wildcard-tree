package parser

import (
	"github.com/jorgepbrown/wildcard-tree/tokenizer"
)

type Parser struct {
	t            *tokenizer.Tokenizer
	peekToken    tokenizer.Token
	currentToken tokenizer.Token
}

func New(t *tokenizer.Tokenizer) *Parser {
	p := &Parser{
		t: t,
	}
	p.peekToken = p.t.Next()
	return p
}

type AST struct {
	root *Wildcard
}

type ExpressionType string
type Expression interface {
	Value() string
	Literal() string // debugging
	Type() ExpressionType
}

type OperatorPriority uint

const (
	LOWEST OperatorPriority = iota
	PIPE
	NULL
	INDEX
	PAREN
)

var opMap = map[tokenizer.TokenType]OperatorPriority{
	tokenizer.PIPE:          PIPE,
	tokenizer.NULL_COALESCE: NULL,
	tokenizer.DOT:           INDEX,
	tokenizer.LBRACKET:      INDEX,
	tokenizer.LPAREN:        PAREN,
}

func (a AST) String() string {
	return a.root.Literal()
}

func (p *Parser) Parse() (AST, error) {
	ast := AST{}

	if p.expect(tokenizer.WILDCARD_OPEN) {
		// c wcopen p text
		p.read()
		// c text p wcclose
		wc, err := p.parseWildcard()
		if err != nil {
			return ast, err
		}
		ast.root = wc
		return ast, nil
	}

	return ast, newSyntaxError("{{", p.currentToken.Literal)
}

func (p *Parser) parseExpression(prio OperatorPriority) (Expression, error) {
	var leftExpr Expression
	switch p.currentToken.T {
	case tokenizer.TEXT:
		leftExpr = &Literal{value: p.currentToken.Literal}
		p.read()
	case tokenizer.WILDCARD_OPEN:
		p.read()
		wc, err := p.parseWildcard()
		if err != nil {
			return nil, err
		}
		leftExpr = wc
	case tokenizer.LPAREN:
		p.read()
		e, err := p.parseParenExpression()
		if err != nil {
			return nil, err
		}
		leftExpr = e
	default:
		return nil, newParserUnkownExprTypeError(p.currentToken.T)
	}

	var err error
	for {
		nextPrio, ok := opMap[p.currentToken.T]
		if !ok {
			return leftExpr, nil
		}
		if nextPrio > prio {
			switch p.currentToken.T {
			case tokenizer.DOT:
				p.read()
				leftExpr, err = p.parseDotExpression(leftExpr)
				if err != nil {
					return nil, err
				}
			case tokenizer.LBRACKET:
				p.read()
				leftExpr, err = p.parseIndexExpression(leftExpr)
				if err != nil {
					return nil, err
				}
			case tokenizer.NULL_COALESCE:
				p.read()
				leftExpr, err = p.parseNullCoalesceExpression(leftExpr)
				if err != nil {
					return nil, err
				}
			case tokenizer.PIPE:
				p.read()
				leftExpr, err = p.parseFunctionExpression(leftExpr)
				if err != nil {
					return nil, err
				}
			default:
				return leftExpr, nil
			}
		} else {
			return leftExpr, nil
		}
	}
}

func (p *Parser) expect(t tokenizer.TokenType) bool {
	if p.peekToken.T == t {
		p.read()
		return true
	}
	return false
}
func (p *Parser) expectCurrent(t tokenizer.TokenType) bool {
	if p.currentToken.T == t {
		p.read()
		return true
	}
	return false
}

func (p *Parser) read() bool {
	p.currentToken = p.peekToken
	if p.currentToken.T == tokenizer.EOF {
		return false
	}
	p.peekToken = p.t.Next()
	return true
}
