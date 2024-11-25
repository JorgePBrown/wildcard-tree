package parser

type Literal struct {
	value string
}

const LITERAL ExpressionType = "LITERAL"

func (l *Literal) Value() string {
	return l.value
}
func (l *Literal) Type() ExpressionType {
	return LITERAL
}
func (l *Literal) Literal() string {
	return l.value
}
