package app

import (
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/cedar-policy/cedar-go"

	"github.com/sagikazarmark/octoslash/authz"
	"github.com/sagikazarmark/octoslash/command"
)

func DefaultAuthorizer(
	provider Provider,
	defaultPolicyLoader LazyResult[authz.PolicyLoader],
	defaultEntityLoader LazyResult[authz.EntityLoader],
	logger *slog.Logger,
) LazyResult[command.Authorizer] {
	return func() (command.Authorizer, error) {
		switch p := provider.(type) {
		case interface {
			NewAuthorizer() command.Authorizer
		}:
			return p.NewAuthorizer(), nil

		case interface {
			NewAuthorizer() (command.Authorizer, error)
		}:
			return p.NewAuthorizer()

		default:
			policies, err := newPolicyIterator(provider, defaultPolicyLoader, logger)
			if err != nil {
				return nil, err
			}

			entities, err := newEntityGetter(provider, defaultEntityLoader, logger)
			if err != nil {
				return nil, err
			}

			return authz.NewAuthorizer(policies, entities, logger), nil
		}
	}
}

func newPolicyIterator(
	provider Provider,
	def LazyResult[authz.PolicyLoader],
	logger *slog.Logger,
) (cedar.PolicyIterator, error) {
	switch p := provider.(type) {
	case interface {
		PolicyIterator() cedar.PolicyIterator
	}:
		return p.PolicyIterator(), nil

	case interface {
		PolicyIterator() (cedar.PolicyIterator, error)
	}:
		return p.PolicyIterator()

	default:
		policyLoader, err := newPolicyLoader(provider, def, logger)
		if err != nil {
			return nil, err
		}

		policies, err := policyLoader.LoadPolicies()
		if err != nil {
			return nil, fmt.Errorf("loading policies: %w", err)
		}

		return policies, nil
	}
}

func newPolicyLoader(
	provider Provider,
	def LazyResult[authz.PolicyLoader],
	logger *slog.Logger,
) (authz.PolicyLoader, error) {
	switch p := provider.(type) {
	case interface {
		PolicyLoader() authz.PolicyLoader
	}:
		return p.PolicyLoader(), nil

	case interface {
		PolicyLoader() (authz.PolicyLoader, error)
	}:
		return p.PolicyLoader()

	default:
		var loader authz.PolicyLoaders

		l, err := def.Resolve()
		if err != nil {
			return nil, err
		}

		if l != nil {
			loader = append(loader, l)
		}

		if p, ok := provider.(interface{ AdditionalPolicyLoaders() authz.PolicyLoaders }); ok {
			loader = append(loader, p.AdditionalPolicyLoaders()...)
		}

		if len(loader) == 0 {
			logger.Warn("no policy loader available, all authorization requests will be denied")
		}

		return loader, nil
	}
}

func DefaultPolicyLoader(fsys LazyResult[fs.FS]) LazyResult[authz.PolicyLoader] {
	return func() (authz.PolicyLoader, error) {
		fsys, err := fsys.Resolve()
		if err != nil {
			return nil, err
		}

		if fsys == nil {
			return nil, nil
		}

		return authz.FilePolicyLoader{Fsys: fsys}, nil
	}
}

func newEntityGetter(
	provider Provider,
	def LazyResult[authz.EntityLoader],
	logger *slog.Logger,
) (cedar.EntityGetter, error) {
	switch p := provider.(type) {
	case interface {
		EntityGetter() cedar.EntityGetter
	}:
		return p.EntityGetter(), nil

	case interface {
		EntityGetter() (cedar.EntityGetter, error)
	}:
		return p.EntityGetter()

	default:
		entityLoader, err := newEntityLoader(provider, def, logger)
		if err != nil {
			return nil, err
		}

		entities, err := entityLoader.LoadEntities()
		if err != nil {
			return nil, fmt.Errorf("loading entities: %w", err)
		}

		return entities, nil
	}
}

func newEntityLoader(
	provider Provider,
	def LazyResult[authz.EntityLoader],
	logger *slog.Logger,
) (authz.EntityLoader, error) {
	switch p := provider.(type) {
	case interface {
		EntityLoader() authz.EntityLoader
	}:
		return p.EntityLoader(), nil

	case interface {
		EntityLoader() (authz.EntityLoader, error)
	}:
		return p.EntityLoader()

	default:
		var loader authz.EntityLoaders

		l, err := def.Resolve()
		if err != nil {
			return nil, err
		}

		if l != nil {
			loader = append(loader, l)
		}

		if p, ok := provider.(interface{ AdditionalEntityLoaders() authz.EntityLoaders }); ok {
			loader = append(loader, p.AdditionalEntityLoaders()...)
		}

		if len(loader) == 0 {
			logger.Warn("no entity loader available, all authorization requests will be denied")
		}

		return loader, nil
	}
}

func DefaultEntityLoader(fsys LazyResult[fs.FS]) LazyResult[authz.EntityLoader] {
	return func() (authz.EntityLoader, error) {
		fsys, err := fsys.Resolve()
		if err != nil {
			return nil, err
		}

		if fsys == nil {
			return nil, nil
		}

		return authz.FileEntityLoader{Fsys: fsys}, nil
	}
}
