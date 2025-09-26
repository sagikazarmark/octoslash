package app

import (
	"errors"
	"log/slog"

	"github.com/google/go-github/v74/github"
	"github.com/sagikazarmark/octoslash"
	"github.com/sagikazarmark/octoslash/command"
)

func NewCommandDispatcher(provider Provider, def LazyResult[octoslash.CommandDispatcher]) (octoslash.CommandDispatcher, error) {
	switch p := provider.(type) {
	case interface {
		NewCommandDispatcher() octoslash.CommandDispatcher
	}:
		return p.NewCommandDispatcher(), nil

	default:
		return def.Resolve()
	}
}

func DefaultCommandDispatcher(authorizer LazyResult[command.Authorizer], commandProvider LazyResult[command.CommandProvider]) LazyResult[octoslash.CommandDispatcher] {
	return func() (octoslash.CommandDispatcher, error) {
		authorizer, err := authorizer.Resolve()
		if err != nil {
			return nil, err
		}

		commandProvider, err := commandProvider.Resolve()
		if err != nil {
			return nil, err
		}

		return command.CobraDispatcher{
			Authorizer:      authorizer,
			CommandProvider: commandProvider,
		}, nil
	}
}

func NewCommandProvider(provider Provider, client *github.Client, logger *slog.Logger) LazyResult[command.CommandProvider] {
	return func() (command.CommandProvider, error) {
		switch p := provider.(type) {
		case interface {
			NewCommandProvider() command.CommandProvider
		}:
			return p.NewCommandProvider(), nil

		case interface {
			NewCommandProvider(client *github.Client) command.CommandProvider
		}:
			return p.NewCommandProvider(client), nil

		case interface {
			NewCommandProvider(client *github.Client, logger *slog.Logger) command.CommandProvider
		}:
			return p.NewCommandProvider(client, logger), nil

		default:
			return nil, errors.New("no command provider")
		}
	}
}
