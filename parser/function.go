package parser

import "fmt"

const FUNCTION ExpressionType = "FUNCTION"

type FunctionExpression struct {
	Argument Expression
	Name     Expression
}

func (w *FunctionExpression) Value() string {
	return w.Literal()
}
func (w *FunctionExpression) Type() ExpressionType {
	return NULL_COALESCE
}

func (w *FunctionExpression) Literal() string {
	return fmt.Sprintf("(%s | %s)", w.Argument.Literal(), w.Name.Literal())
}

func (p *Parser) parseFunctionExpression(primary Expression) (*FunctionExpression, error) {
	expr, err := p.parseExpression(PIPE)
	if err != nil {
		return nil, err
	}
	return &FunctionExpression{
		Argument: primary,
		Name:     expr,
	}, nil
}
