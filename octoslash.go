package octoslash

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/go-github/v74/github"
	"github.com/sagikazarmark/octoslash/slash"
)

type EventHandler struct {
	Dispatcher CommandDispatcher
	// ErrorHandler ErrorHandler
}

type CommandDispatcher interface {
	Dispatch(ctx context.Context, event github.IssueCommentEvent, args []string) error
}

// Error handling behavior: ignore (debug log), warnButIgnore (error log), return (fails the command)
// TODO: wait for other commands?
type ErrorHandler interface {
	Handle(ctx context.Context, err error)
}

func (h EventHandler) Handle(ctx context.Context, event github.IssueCommentEvent) error {
	logger := slog.Default()

	rawCommands := slash.ScanString(event.GetComment().GetBody())

	if len(rawCommands) == 0 {
		logger.Info("no commands to run")

		return nil
	}

	parser := slash.NewParser()

	for _, rawCommand := range rawCommands {
		args, err := parser.Parse(strings.NewReader(rawCommand))
		if err != nil {
			logger.Error(fmt.Sprintf("parsing command: %s", err.Error()), slog.String("command", rawCommand))

			// TODO: make this behavior configurable
			continue
		}

		logger.Debug("running command", slog.String("command", rawCommand))

		err = h.Dispatcher.Dispatch(ctx, event, args)
		if err != nil {
			return err
		}
	}

	return nil
}
