package parser

import (
	"fmt"
)

const NULL_COALESCE ExpressionType = "NULL_COALESCE"

type NullCoalesceExpression struct {
	Primary  Expression
	Fallback Expression
}

func (w *NullCoalesceExpression) Value() string {
	return w.Literal()
}
func (w *NullCoalesceExpression) Type() ExpressionType {
	return NULL_COALESCE
}

func (w *NullCoalesceExpression) Literal() string {
	return fmt.Sprintf("(%s ?? %s)", w.Primary.Literal(), w.Fallback.Literal())
}

func (p *Parser) parseNullCoalesceExpression(primary Expression) (*NullCoalesceExpression, error) {
	expr, err := p.parseExpression(NULL)
	if err != nil {
		return nil, err
	}
	return &NullCoalesceExpression{
		Primary:  primary,
		Fallback: expr,
	}, nil
}
