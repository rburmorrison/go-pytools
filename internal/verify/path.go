package verify

import (
	"errors"
	"os"
	"path/filepath"
)

// Error definitions
var (
	ErrIsNotAbs     = errors.New("verify: path is not absolute")
	ErrIsNotFile    = errors.New("verify: path does not lead to a file")
	ErrPathNotExist = errors.New("verify: path does not exist")
)

// FilePath accepts a path to a file and returns an
// error if the path does not exist, is not absolute,
// or if it does not lead to a file.
func FilePath(path string) error {
	if !filepath.IsAbs(path) {
		return ErrIsNotAbs
	}

	fi, err := os.Stat(path)
	if err != nil {
		return ErrPathNotExist
	}

	if !fi.Mode().IsRegular() {
		return ErrIsNotFile
	}

	return nil
}
