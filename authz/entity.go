package authz

import (
	"fmt"

	"github.com/cedar-policy/cedar-go"
	"github.com/google/go-github/v74/github"
)

const (
	User cedar.EntityType = "User"

	Repository cedar.EntityType = "Repository"
	Owner      cedar.EntityType = "Owner"

	Issue       cedar.EntityType = "Issue"
	PullRequest cedar.EntityType = "PullRequest"

	Role cedar.EntityType = "Role"

	Action cedar.EntityType = "Action"
)

func NewUserID(user *github.User) cedar.EntityUID {
	return NewEntityUID(User, user)
}

func NewIssueOrPullRequestID(issue *github.Issue) cedar.EntityUID {
	t := Issue

	if issue.IsPullRequest() {
		t = PullRequest
	}

	return NewEntityUID(t, issue)
}

func NewEntityUID(entityType cedar.EntityType, entity interface{ GetID() int64 }) cedar.EntityUID {
	return cedar.NewEntityUID(entityType, cedar.String(fmt.Sprintf("%d", entity.GetID())))
}

func NewOwner(owner *github.User) cedar.Entity {
	uid := NewEntityUID(Owner, owner)

	return cedar.Entity{
		UID: uid,
		Attributes: cedar.NewRecord(cedar.RecordMap{
			cedar.String("login"): cedar.String(owner.GetLogin()),
		}),
	}
}

func NewRepository(repo *github.Repository) cedar.Entity {
	uid := NewEntityUID(Repository, repo)

	return cedar.Entity{
		UID:     uid,
		Parents: cedar.NewEntityUIDSet(NewEntityUID(Owner, repo.GetOwner())),
		Attributes: cedar.NewRecord(cedar.RecordMap{
			cedar.String("name"): cedar.String(repo.GetName()),
		}),
	}
}

func NewIssueOrPullRequest(issue *github.Issue, repo *github.Repository) cedar.Entity {
	uid := NewIssueOrPullRequestID(issue)

	attributes := cedar.RecordMap{
		cedar.String("number"): cedar.Long(issue.GetNumber()),
	}

	var labels []cedar.Value

	for _, label := range issue.Labels {
		labels = append(labels, cedar.String(label.GetName()))
	}

	if len(labels) > 0 {
		attributes[cedar.String("labels")] = cedar.NewSet(labels...)
	}

	entity := cedar.Entity{
		UID:        uid,
		Parents:    cedar.NewEntityUIDSet(NewEntityUID(Repository, issue.GetRepository())),
		Attributes: cedar.NewRecord(attributes),
	}

	return entity
}
