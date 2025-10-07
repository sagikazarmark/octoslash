package parser

import (
	"errors"
	"io"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/syntax"
)

type Parser struct {
	parser *syntax.Parser
}

func NewParser() *Parser {
	return &Parser{
		parser: syntax.NewParser(),
	}
}

func (p *Parser) Parse(r io.Reader) ([]string, error) {
	file, err := p.parser.Parse(r, "")
	if err != nil {
		return nil, err
	}

	if len(file.Stmts) == 0 {
		return nil, errors.New("empty input")
	}

	cmd, ok := file.Stmts[0].Cmd.(*syntax.CallExpr)
	if !ok {
		return nil, errors.New("not a call expression")
	}

	printer := syntax.NewPrinter()

	config := &expand.Config{
		Env: expand.FuncEnviron(func(variable string) string { return "$" + variable }),
		CmdSubst: func(w io.Writer, stmt *syntax.CmdSubst) error {
			return printer.Print(w, stmt)
		},
	}

	var args []string
	for _, word := range cmd.Args {
		expanded, err := expand.Literal(config, word)
		if err != nil {
			return nil, err
		}

		args = append(args, expanded)
	}

	return args, nil
}
