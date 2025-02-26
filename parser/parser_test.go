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
				&Wildcard{Expression: &Literal{V: "a"}},
			},
		},
		{
			"{{a.a}}",
			AST{
				&Wildcard{Expression: &DotExpression{
					Target: &Literal{V: "a"},
					Key:    &Literal{V: "a"},
				}},
			},
		},
		{
			"{{a.{{a}}}}",
			AST{
				&Wildcard{Expression: &DotExpression{
					Target: &Literal{V: "a"},
					Key:    &Wildcard{Expression: &Literal{V: "a"}},
				}},
			},
		},
		{
			"{{a.{{a.{{b}}}}}}",
			AST{
				&Wildcard{Expression: &DotExpression{
					Target: &Literal{V: "a"},
					Key: &Wildcard{Expression: &DotExpression{
						Target: &Literal{V: "a"},
						Key:    &Wildcard{Expression: &Literal{V: "b"}},
					}},
				}},
			},
		},
		{
			`{{a[""a"]}}`,
			AST{
				&Wildcard{Expression: &IndexExpression{
					Target: &Literal{V: "a"},
					Key:    &Literal{V: "\"a"},
				}},
			},
		},
		{
			"{{a[{{a}}]}}",
			AST{
				&Wildcard{Expression: &IndexExpression{
					Target: &Literal{V: "a"},
					Key:    &Wildcard{Expression: &Literal{V: "a"}},
				}},
			},
		},
		{
			`{{"{{a}}"}}`,
			AST{
				&Wildcard{Expression: &Literal{V: "{{a}}"}},
			},
		},
		{
			`{{ "{{a}}" }}`,
			AST{
				&Wildcard{Expression: &Literal{V: "{{a}}"}},
			},
		},
		{
			`{{ "{{a}}" ?? a }}`,
			AST{
				&Wildcard{Expression: &NullCoalesceExpression{
					Primary:  &Literal{V: "{{a}}"},
					Fallback: &Literal{V: "a"},
				}},
			},
		},
		{
			`{{ "{{a}}" ?? a | toUpper }}`,
			AST{
				&Wildcard{Expression: &FunctionExpression{Argument: &NullCoalesceExpression{
					Primary:  &Literal{V: "{{a}}"},
					Fallback: &Literal{V: "a"},
				}, Name: &Literal{V: "toUpper"}}},
			},
		},
		{
			`{{ "{{a}}" ?? (a | toUpper) }}`,
			AST{
				&Wildcard{Expression: &NullCoalesceExpression{
					Primary: &Literal{V: "{{a}}"},
					Fallback: &FunctionExpression{
						Argument: &Literal{V: "a"},
						Name:     &Literal{V: "toUpper"},
					},
				}},
			},
		},
		{
			`{{ a | toUpper ?? "{{a}}" }}`,
			AST{
				&Wildcard{Expression: &FunctionExpression{
					Argument: &Literal{V: "a"},
					Name: &NullCoalesceExpression{
						Primary:  &Literal{"toUpper"},
						Fallback: &Literal{V: "{{a}}"},
					},
				}},
			},
		},
		{
			`{{ a | toUpper ?? {{a}} }}`,
			AST{
				&Wildcard{Expression: &FunctionExpression{
					Argument: &Literal{V: "a"},
					Name: &NullCoalesceExpression{
						Primary:  &Literal{"toUpper"},
						Fallback: &Wildcard{Expression: &Literal{V: "a"}},
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
	if expected.Root == nil && actual.Root != nil {
		t.Errorf("expected nil ast wildcard")
		return
	} else if expected.Root != nil && actual.Root == nil {
		t.Error("expected non-nil ast wildcard")
		return
	} else if expected.Root == nil && actual.Root == nil {
		return
	}
	testExpr(expected.Root, actual.Root, t)
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
		testExpr(v.Expression, actual.(*Wildcard).Expression, t)
	}
}
