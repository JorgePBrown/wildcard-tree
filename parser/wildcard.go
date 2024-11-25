package parser

import (
	"fmt"

	"github.com/jorgepbrown/wildcard-tree/tokenizer"
)

const WILDCARD ExpressionType = "WILDCARD"

type Wildcard struct {
	value Expression
}

func (w *Wildcard) Value() string {
	return w.value.Value()
}
func (w *Wildcard) Type() ExpressionType {
	return WILDCARD
}

func (w *Wildcard) Literal() string {
	return fmt.Sprintf("{{%s}}", w.value.Literal())
}

func (p *Parser) parseWildcard() (*Wildcard, error) {
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if !p.expectCurrent(tokenizer.WILDCARD_CLOSE) {
		return nil, newSyntaxError("}}", p.currentToken.Literal)
	}
	return &Wildcard{
		value: expr,
	}, nil
}
