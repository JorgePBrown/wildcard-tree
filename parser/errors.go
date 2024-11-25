package parser

import (
	"fmt"

	"github.com/jorgepbrown/wildcard-tree/tokenizer"
)

const (
	INVALID_SYNTAX    = "invalid wildcard syntax: expected '%s' found '%s'"
	UNKNOWN_EXPR_TYPE = "parser error unknown epression type %s"
	MALFORMED_EXPR    = "parser error malformed expression '%s' followed by '%s'"
)

func newSyntaxError(expected, found string) error {
	return fmt.Errorf(INVALID_SYNTAX, expected, found)
}

func newParserUnkownExprTypeError(t tokenizer.TokenType) error {
	return fmt.Errorf(UNKNOWN_EXPR_TYPE, t)
}

func newParserMalformedExprError(l1, l2 string) error {
	return fmt.Errorf(MALFORMED_EXPR, l1, l2)
}
