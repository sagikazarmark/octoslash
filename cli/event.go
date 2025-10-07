package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/google/go-github/v74/github"
)

func LoadEventFromFile(
	os Options,
	eventName string,
	eventPath string,
) (github.IssueCommentEvent, error) {
	file, err := os.Open(eventPath)
	if err != nil {
		return github.IssueCommentEvent{}, fmt.Errorf("loading event: %w", err)
	}
	defer file.Close()

	return LoadEvent(file)
}

func LoadEvent(r io.Reader) (github.IssueCommentEvent, error) {
	var event github.IssueCommentEvent

	err := json.NewDecoder(r).Decode(&event)
	if err != nil {
		return event, fmt.Errorf("decoding event: %w", err)
	}

	return event, nil
}
