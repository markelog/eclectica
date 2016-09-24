package io

import (
	"io/ioutil"
	"os"
)

func CreateDir(path string) (string, error) {
	err := os.MkdirAll(path, 0700)

	if err != nil {
		return "", err
	}

	return path, nil
}

func WriteFile(path, content string) error {
	data := []byte(content)

	err := ioutil.WriteFile(path, data, 0700)
	if err != nil {
		return err
	}

	return nil
}
