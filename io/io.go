package io

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Walker func(path string) bool

func walkUp(path string, fn Walker) {
	current := path

	stop := fn(current)
	if stop == true {
		return
	}

	for {
		if current == "" || current == "/" {
			return
		}

		current = filepath.Dir(current)

		stop = fn(current)
		if stop == true {
			return
		}
	}

	return
}

func GetVersion(args ...interface{}) (version string, err error) {
	var path string

	if len(args) > 1 {
		path, err = FindDotFile(args[0], args[1])
	} else {
		path, err = FindDotFile(args[0])
	}

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

func FindDotFile(args ...interface{}) (versionPath string, err error) {
	var path string
	dots := args[0].([]string)

	if len(args) > 1 {
		path = args[1].(string)
	} else {

		path, err = os.Getwd()
		if err != nil {
			return
		}
	}

	walkUp(path, func(path string) bool {
		for _, file := range dots {
			p := filepath.Join(path, file)

			if _, err := os.Stat(p); err == nil {
				versionPath = p
				return true
			}
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

func Read(path string) string {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return ""
	}

	return string(bytes)
}
