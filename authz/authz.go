package authz

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cedar-policy/cedar-go"
	"github.com/google/go-github/v74/github"
)

type Authorizer struct {
	Policies cedar.PolicyIterator
	Entities cedar.EntityGetter

	Logger *slog.Logger
}

func NewAuthorizer(
	policies cedar.PolicyIterator,
	entities cedar.EntityGetter,
	logger *slog.Logger,
) Authorizer {
	return Authorizer{
		Policies: policies,
		Entities: entities,
		Logger:   logger,
	}
}

func (a Authorizer) Authorize(
	ctx context.Context,
	event github.IssueCommentEvent,
	action string,
) error {
	request := newRequest(event, action)

	a.Logger.Debug(
		"authorizing request",
		slog.String("principal", request.Principal.String()),
		slog.String("resource", request.Resource.String()),
		slog.String("action", request.Action.String()),
	)

	// TODO: clean this up
	eventEntities, _ := (EventEntityLoader{Event: event}).LoadEntities()

	entities := EntityGetters{
		eventEntities,
		a.Entities,
	}

	ok, _ := cedar.Authorize(a.Policies, entities, request)
	if !ok {
		return fmt.Errorf(
			"principal %s is not authorized to perform %s on %s",
			request.Principal.String(),
			request.Action.String(),
			request.Resource.String(),
		)
	}

	return nil
}

func newRequest(event github.IssueCommentEvent, action string) cedar.Request {
	return cedar.Request{
		Principal: NewUserID(event.GetComment().GetUser()),
		Action:    cedar.NewEntityUID(Action, cedar.String(action)),
		Resource:  NewIssueOrPullRequestID(event.GetIssue()),
	}
}
