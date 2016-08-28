package directory

import (
	"os"
)

func Create(path string) (string, error) {
	err := os.MkdirAll(path, 0700)

	if err != nil {
		return "", err
	}

	return path, nil
}
