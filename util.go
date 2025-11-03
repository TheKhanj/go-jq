package jq

import (
	"errors"
	"os"
)

func assertDir(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return errors.New("path is not a directory: " + path)
		}

		return nil
	}

	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return errors.New("failed to create directory: " + path)
		}
		return nil
	}

	return err
}
