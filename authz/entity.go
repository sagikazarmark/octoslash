package authz

import (
	"encoding/json"
	"errors"
	"io/fs"

	"github.com/cedar-policy/cedar-go"
	"github.com/google/go-github/v74/github"
)

var _ cedar.EntityGetter = (EntityGetters)(nil)

// EntityGetters is a list of [cedar.EntityGetter] implementations.
type EntityGetters []cedar.EntityGetter

func (g EntityGetters) Get(uid cedar.EntityUID) (cedar.Entity, bool) {
	for _, getter := range g {
		entity, ok := getter.Get(uid)
		if ok {
			return entity, true
		}
	}

	return cedar.Entity{}, false
}

type EntityLoader interface {
	LoadEntities() (cedar.EntityGetter, error)
}

type EntityLoaders []EntityLoader

func (l EntityLoaders) LoadEntities() (cedar.EntityGetter, error) {
	var entities EntityGetters

	for _, loader := range l {
		loadedEntities, err := loader.LoadEntities()
		if err != nil {
			return nil, err
		}

		entities = append(entities, loadedEntities)
	}

	return entities, nil
}

type EventEntityLoader struct {
	Event github.IssueCommentEvent
}

func (l EventEntityLoader) LoadEntities() (cedar.EntityGetter, error) {
	entities := cedar.EntityMap{}

	owner := NewOwner(l.Event.GetRepo().GetOwner())
	repo := NewRepository(l.Event.GetRepo())
	issue := NewIssueOrPullRequest(l.Event.GetIssue(), l.Event.GetRepo())

	entities[owner.UID] = owner
	entities[repo.UID] = repo
	entities[issue.UID] = issue

	return entities, nil
}

type FileEntityLoader struct {
	Fsys fs.FS
}

func (l FileEntityLoader) LoadEntities() (cedar.EntityGetter, error) {
	var entities cedar.EntityMap

	file, err := l.Fsys.Open("principals.json")
	if errors.Is(err, fs.ErrNotExist) {
		return entities, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&entities); err != nil {
		return nil, err
	}

	return entities, nil
}
