package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/jorgepbrown/wildcard-tree/repl"
)

var step = flag.String("step", "parse", "--step=<step> to specify the step")

func init() {
	flag.Parse()
}

func main() {
	r := repl.New()
	var s repl.ParseStep
	switch strings.ToLower(*step) {
	case "tokenize":
		s = repl.TOKENIZE
	case "parse":
		s = repl.PARSE
	default:
		fmt.Printf("unknown step %s", *step)
		return
	}
	r.Start(s)
}
