package parser

import "github.com/jorgepbrown/wildcard-tree/tokenizer"

func (p *Parser) parseParenExpression() (Expression, error) {
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if !p.expectCurrent(tokenizer.RPAREN) {
		return nil, newSyntaxError(")", p.currentToken.Literal)
	}
	return expr, nil
}
