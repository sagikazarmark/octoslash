package cli

import (
	"errors"
	"fmt"
	"io/fs"
)

const defaultConfigPath = ".github/octoslash"

func DefaultFilesystem(os Options) (fs.FS, error) {
	fsys, err := openLocal(os, defaultConfigPath)
	if err != nil {
		return nil, fmt.Errorf("opening local config: %w", err)
	}

	return fsys, nil
}

func openLocal(os Options, configPath string) (fs.FS, error) {
	root, err := os.OpenRoot(configPath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return root.FS(), nil
}
