package parser

import (
	"bytes"

	"github.com/jorgepbrown/wildcard-tree/tokenizer"
)

const INDEX_EXPR ExpressionType = "INDEX"

type IndexExpression struct {
	Target Expression
	Key    Expression
}

func (p *Parser) parseIndexExpression(target Expression) (*IndexExpression, error) {
	key, err := p.parseExpression(INDEX)
	if err != nil {
		return nil, err
	}
	if !p.expectCurrent(tokenizer.RBRACKET) {
		return nil, newSyntaxError("]", p.currentToken.Literal)
	}
	return &IndexExpression{
		Target: target,
		Key:    key,
	}, nil
}

func (e *IndexExpression) Type() ExpressionType {
	return INDEX_EXPR
}
func (e *IndexExpression) Value() string {
	return e.Literal()
}
func (e *IndexExpression) Literal() string {
	var out bytes.Buffer
	out.WriteString(e.Target.Literal())
	out.WriteByte('[')
	out.WriteString(e.Key.Literal())
	out.WriteByte(']')
	return out.String()
}
