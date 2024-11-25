package parser

import (
	"testing"

	"github.com/jorgepbrown/wildcard-tree/tokenizer"
)

func TestParser(t *testing.T) {
	tt := []struct {
		input    string
		expected AST
	}{
		{
			"{{a}}",
			AST{
				&Wildcard{value: &Literal{value: "a"}},
			},
		},
		{
			"{{a.a}}",
			AST{
				&Wildcard{value: &DotExpression{
					Target: &Literal{value: "a"},
					Key:    &Literal{value: "a"},
				}},
			},
		},
		{
			"{{a.{{a}}}}",
			AST{
				&Wildcard{value: &DotExpression{
					Target: &Literal{value: "a"},
					Key:    &Wildcard{value: &Literal{value: "a"}},
				}},
			},
		},
		{
			"{{a.{{a.{{b}}}}}}",
			AST{
				&Wildcard{value: &DotExpression{
					Target: &Literal{value: "a"},
					Key: &Wildcard{value: &DotExpression{
						Target: &Literal{value: "a"},
						Key:    &Wildcard{value: &Literal{value: "b"}},
					}},
				}},
			},
		},
		{
			"{{a[a]}}",
			AST{
				&Wildcard{value: &IndexExpression{
					Target: &Literal{value: "a"},
					Key:    &Literal{value: "a"},
				}},
			},
		},
		{
			"{{a[{{a}}]}}",
			AST{
				&Wildcard{value: &IndexExpression{
					Target: &Literal{value: "a"},
					Key:    &Wildcard{value: &Literal{value: "a"}},
				}},
			},
		},
		{
			`{{"{{a}}"}}`,
			AST{
				&Wildcard{value: &Literal{value: "{{a}}"}},
			},
		},
		{
			`{{ "{{a}}" }}`,
			AST{
				&Wildcard{value: &Literal{value: "{{a}}"}},
			},
		},
		{
			`{{ "{{a}}" ?? a }}`,
			AST{
				&Wildcard{value: &NullCoalesceExpression{
					Primary:  &Literal{value: "{{a}}"},
					Fallback: &Literal{value: "a"},
				}},
			},
		},
		{
			`{{ "{{a}}" ?? a | toUpper }}`,
			AST{
				&Wildcard{value: &FunctionExpression{Argument: &NullCoalesceExpression{
					Primary:  &Literal{value: "{{a}}"},
					Fallback: &Literal{value: "a"},
				}, Name: &Literal{value: "toUpper"}}},
			},
		},
		{
			`{{ "{{a}}" ?? (a | toUpper) }}`,
			AST{
				&Wildcard{value: &NullCoalesceExpression{
					Primary: &Literal{value: "{{a}}"},
					Fallback: &FunctionExpression{
						Argument: &Literal{value: "a"},
						Name:     &Literal{value: "toUpper"},
					},
				}},
			},
		},
		{
			`{{ a | toUpper ?? "{{a}}" }}`,
			AST{
				&Wildcard{value: &FunctionExpression{
					Argument: &Literal{value: "a"},
					Name: &NullCoalesceExpression{
						Primary:  &Literal{"toUpper"},
						Fallback: &Literal{value: "{{a}}"},
					},
				}},
			},
		},
		{
			`{{ a | toUpper ?? {{a}} }}`,
			AST{
				&Wildcard{value: &FunctionExpression{
					Argument: &Literal{value: "a"},
					Name: &NullCoalesceExpression{
						Primary:  &Literal{"toUpper"},
						Fallback: &Wildcard{value: &Literal{value: "a"}},
					},
				}},
			},
		},
	}

	for i, test := range tt {
		t.Logf("parser-%d %s", i, test.input)
		p := New(tokenizer.New(test.input))
		ast, err := p.Parse()
		if err != nil {
			t.Fatal(err)
		}
		testAST(&test.expected, &ast, t)
	}
}

func testAST(expected, actual *AST, t *testing.T) {
	t.Helper()
	if expected.root == nil && actual.root != nil {
		t.Errorf("expected nil ast wildcard")
		return
	} else if expected.root != nil && actual.root == nil {
		t.Error("expected non-nil ast wildcard")
		return
	} else if expected.root == nil && actual.root == nil {
		return
	}
	testExpr(expected.root, actual.root, t)
}

func testExpr(expected, actual Expression, t *testing.T) {
	t.Helper()
	if expected == nil && actual != nil {
		t.Errorf("expected nil expression, got=%s", actual.Type())
		return
	} else if expected != nil && actual == nil {
		t.Errorf("expected non-nil expression, expected=%s", expected.Type())
		return
	} else if expected == nil && actual == nil {
		return
	}

	if expected.Type() != actual.Type() {
		t.Errorf("wrong expression type, expected=%s got=%s", expected.Type(), actual.Type())
	}
	if expected.Literal() != actual.Literal() {
		t.Errorf("wrong expression literal, expected=%s got=%s", expected.Literal(), actual.Literal())
	}
	if expected.Value() != actual.Value() {
		t.Errorf("wrong expression value, expected=%s got=%s", expected.Value(), actual.Value())
	}
	if t.Failed() {
		t.FailNow()
	}

	switch v := expected.(type) {
	case *Wildcard:
		testExpr(v.value, actual.(*Wildcard).value, t)
	}
}
