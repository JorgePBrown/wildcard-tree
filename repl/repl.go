package repl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jorgepbrown/wildcard-tree/parser"
	"github.com/jorgepbrown/wildcard-tree/tokenizer"
)

const PREFIX = "> "

type Repl struct{}

func New() *Repl {
	return &Repl{}
}

type ParseStep int

const (
	TOKENIZE ParseStep = iota
	PARSE
)

func (r *Repl) Start(step ParseStep) {
	reader := bufio.NewReader(os.Stdin)
	for {
		_, err := os.Stdout.WriteString(PREFIX)
		if err != nil {
			panic(err)
		}
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		var out bytes.Buffer
		out.WriteString(string(line))
		for isPrefix {
			line, isPrefix, err = reader.ReadLine()
			if err != nil {
				panic(err)
			}
			out.WriteString(string(line))
		}

		t := tokenizer.New(out.String())

		if step == TOKENIZE {
			for {
				token := t.Next()
				if token.T == tokenizer.EOF {
					break
				}

				_, err := os.Stdout.WriteString(fmt.Sprintf("%s %s\n", token.T, token.Literal))
				if err != nil {
					panic(err)
				}
			}
		} else if step == PARSE {
			p := parser.New(t)

			ast, err := p.Parse()
			if err == nil {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("  ", "  ")
				err := enc.Encode(ast)
				if err != nil {
					panic(err)
				}
			} else {
				_, er := os.Stderr.WriteString(err.Error())
				if er != nil {
					panic(er)
				}
			}
		}
	}
}
