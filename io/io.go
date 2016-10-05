package io

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Walker func(path string) bool

func walkUp(path string, fn Walker) {
	current := path

	for {
		if current == "" || current == "/" {
			return
		}

		current = filepath.Dir(current)

		stop := fn(current)
		if stop == true {
			return
		}
	}

	return
}

func GetVersion(language string) (version string, err error) {
	path, err := FindDotFile(language)
	if err != nil {
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "current", nil
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		return scanner.Text(), nil
	}

	if scannerErr := scanner.Err(); scannerErr != nil {
		return "", scannerErr
	}

	return "current", nil
}

func FindDotFile(language string) (versionPath string, err error) {
	pwd, err := os.Getwd()
	file := fmt.Sprintf(".%s-version", language)
	if err != nil {
		return
	}

	walkUp(pwd, func(path string) bool {
		p := filepath.Join(path, file)

		if _, err := os.Stat(p); err == nil {
			versionPath = p
			return true
		}

		return false
	})

	return
}

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
