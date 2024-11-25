package parser

import "bytes"

const DOT_EXPR ExpressionType = "DOT"

type DotExpression struct {
	Target Expression
	Key    Expression
}

func (p *Parser) parseDotExpression(target Expression) (*DotExpression, error) {
	key, err := p.parseExpression(INDEX)
	if err != nil {
		return nil, err
	}
	return &DotExpression{
		Target: target,
		Key:    key,
	}, nil
}

func (e *DotExpression) Type() ExpressionType {
	return DOT_EXPR
}
func (e *DotExpression) Value() string {
	return e.Literal()
}
func (e *DotExpression) Literal() string {
	var out bytes.Buffer
	out.WriteString(e.Target.Literal())
	out.WriteByte('.')
	out.WriteString(e.Key.Literal())
	return out.String()
}
