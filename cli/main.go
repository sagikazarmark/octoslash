package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sagikazarmark/octoslash/internal/app"
)

func init() {
	cobra.EnableTraverseRunHooks = true
}

type Options struct {
	Args     []string
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
	Getenv   func(key string) string
	Open     func(name string) (*os.File, error)
	OpenRoot func(name string) (*os.Root, error)
}

func DefaultOptions() Options {
	return Options{
		Args:     os.Args,
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Getenv:   os.Getenv,
		Open:     os.Open,
		OpenRoot: os.OpenRoot,
	}
}

type Application struct {
	Provider Provider
}

type Provider = app.Provider

func (a Application) Main(os Options) error {
	if a.Provider == nil {
		return errors.New("octoslash: no provider is set")
	}

	flags := pflag.NewFlagSet("octoslash", pflag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	var eventName string
	flags.StringVar(&eventName, "event-name", os.Getenv("GITHUB_EVENT_NAME"), "")

	var eventPath string
	flags.StringVar(&eventPath, "event-path", os.Getenv("GITHUB_EVENT_PATH"), "")

	defaultConfigPath := filepath.Join(".github", "octoslash")

	var configPath string
	flags.StringVar(&configPath, "config-path", defaultConfigPath, "")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	if eventName != "issue_comment" {
		return fmt.Errorf("unsupported event: %s", eventName)
	}

	event, err := LoadEventFromFile(os, eventName, eventPath)
	if err != nil {
		return fmt.Errorf("loading event: %w", err)
	}

	var localFS fs.FS

	handler, err := app.InitializeEventHandler(
		a.Provider,
		app.Token(os.Getenv("GITHUB_TOKEN")),
		event.GetRepo(),
		localFS,
	)
	if err != nil {
		return fmt.Errorf("initializing event handler: %w", err)
	}

	return handler.Handle(context.Background(), event)
}
