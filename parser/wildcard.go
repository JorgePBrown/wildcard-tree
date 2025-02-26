package parser

import (
	"fmt"

	"github.com/jorgepbrown/wildcard-tree/tokenizer"
)

const WILDCARD ExpressionType = "WILDCARD"

type Wildcard struct {
	Expression Expression
}

func (w *Wildcard) Value() string {
	return w.Expression.Value()
}
func (w *Wildcard) Type() ExpressionType {
	return WILDCARD
}

func (w *Wildcard) Literal() string {
	return fmt.Sprintf("{{%s}}", w.Expression.Literal())
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
		Expression: expr,
	}, nil
}
