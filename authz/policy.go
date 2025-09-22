package authz

import (
	"errors"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/cedar-policy/cedar-go"
)

type FilePolicyLoader struct {
	Fsys fs.FS
}

func NewFilePolicyLoader(fsys fs.FS) FilePolicyLoader {
	return FilePolicyLoader{
		Fsys: fsys,
	}
}

func (l FilePolicyLoader) LoadPolicies() (cedar.PolicyIterator, error) {
	if l.Fsys == nil {
		return nil, errors.New("filesystem is not configured")
	}

	policies := cedar.NewPolicySet()

	_, err := fs.Stat(l.Fsys, "policy.cedar")
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	if err == nil {
		return nil, err
	}

	file, err := l.Fsys.Open("policy.cedar")
	if err == nil {
		defer file.Close()

		policy, err := l.unmarshalPolicy(file)
		if err != nil {
			return nil, err
		}

		policies.Add("policy", policy)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	fsys, err := fs.Sub(l.Fsys, "policies")
	if err != nil {
		return nil, err
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			file, err := fsys.Open(entry.Name())
			if err != nil {
				return nil, err
			}
			defer file.Close()

			policy, err := l.unmarshalPolicy(file)
			if err != nil {
				return nil, err
			}

			base := filepath.Base(entry.Name())
			ext := filepath.Ext(base)

			policies.Add(cedar.PolicyID(base[:len(base)-len(ext)]), policy)
		}

	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	return policies, nil
}

func (l FilePolicyLoader) unmarshalPolicy(r io.Reader) (*cedar.Policy, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var policy cedar.Policy
	if err := policy.UnmarshalCedar(b); err != nil {
		return nil, err
	}

	return &policy, nil
}
