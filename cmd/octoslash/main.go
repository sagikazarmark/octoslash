package main

import (
	"fmt"

	"github.com/sagikazarmark/octoslash/builtin"
	"github.com/sagikazarmark/octoslash/cli"
)

func main() {
	opts := cli.DefaultOptions()

	app := cli.Application{
		Provider: builtin.Provider{},
	}

	err := app.Main(opts)
	if err != nil {
		fmt.Println(err)
	}
}
