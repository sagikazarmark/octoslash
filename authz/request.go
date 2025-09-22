package authz

import (
	"github.com/cedar-policy/cedar-go"
	"github.com/google/go-github/v74/github"
	"github.com/spf13/cobra"
)

func NewRequest(cmd *cobra.Command, event github.IssueCommentEvent) cedar.Request {
	return cedar.Request{
		Principal: NewUserID(event.GetComment().GetUser()),
		Action:    cedar.NewEntityUID(Action, cedar.String(cmd.Name())),
		Resource:  NewIssueOrPullRequestID(event.GetIssue()),
	}
}
