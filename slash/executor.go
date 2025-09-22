package slash

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

type Executor struct {
	NewCommand func() *cobra.Command
}

func (e Executor) Execute(ctx context.Context, text string) error {
	cmdLines := ScanString(text)

	parser := NewParser()

	for _, line := range cmdLines {
		args, err := parser.Parse(strings.NewReader(line))
		if err != nil {
			return err
		}

		cmd := e.NewCommand()

		cmd.SetArgs(args)
		cmd.SetContext(ctx)

		if err := cmd.Execute(); err != nil {
			return err
		}
	}

	return nil
}
