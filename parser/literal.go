package parser

type Literal struct {
	V string
}

const LITERAL ExpressionType = "LITERAL"

func (l *Literal) Value() string {
	return l.V
}
func (l *Literal) Type() ExpressionType {
	return LITERAL
}
func (l *Literal) Literal() string {
	return l.V
}
