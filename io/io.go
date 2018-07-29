// Package io provides some helpful IO functions with simplified
// signatures in the context of eclectica package
package io

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-errors/errors"
)

const (
	perm = 0700
)

var (
	versionPattern = `(\d+(\.\d+)?(\.\d+)?)|(latest)`
	rVersion       = regexp.MustCompile(versionPattern)
)

// Walker signature function
type Walker func(path string) bool

// walkUp walker up to filesystem tree
func walkUp(path string, fn Walker) {
	current := path

	stop := fn(current)
	if stop == true {
		return
	}

	for {
		if current == "" || current == "/" {
			break
		}

		current = filepath.Dir(current)

		stop = fn(current)
		if stop == true {
			break
		}
	}

	return
}

// ExtractVersion from the string
func ExtractVersion(file string) (string, error) {
	match := rVersion.FindAllStringSubmatch(file, 1)
	if len(match) == 0 {
		return "", errors.New("There is no version here")
	}

	version := match[0][0]

	return version, nil
}

// GetVersion finds a file by provided argument and extracts
// the version defined in it
func GetVersion(args ...interface{}) (version, path string, err error) {
	current := "current"

	if len(args) > 1 {
		path, err = FindDotFile(args[0], args[1])
	} else {
		path, err = FindDotFile(args[0])
	}

	if err != nil {
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return current, "", nil
	}

	file, err := os.Open(path)
	if err != nil {
		err = errors.New(err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		version, extractErr := ExtractVersion(scanner.Text())

		if extractErr != nil {
			return "", "", extractErr
		}

		return version, path, extractErr
	}

	if scannerErr := scanner.Err(); scannerErr != nil {
		return "", "", scannerErr
	}

	return current, "", nil
}

// FindDotFile finds file up in the filesystem tree
// by provided list of possible files
func FindDotFile(args ...interface{}) (versionPath string, err error) {
	var path string
	dots := args[0].([]string)

	if len(args) > 1 {
		path = args[1].(string)
	} else {

		path, err = os.Getwd()
		if err != nil {
			err = errors.New(err)
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

// CreateDir creates dir with predefined perms
func CreateDir(path string) (string, error) {
	err := os.MkdirAll(path, perm)

	if err != nil {
		return "", errors.New(err)
	}

	return path, nil
}

// WriteFile writes file with default perms and accepts a string as data
func WriteFile(path, content string) error {
	data := []byte(content)

	err := ioutil.WriteFile(path, data, perm)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Read file and return a content string
func Read(path string) string {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return ""
	}

	return string(bytes)
}

// ListVersions lists installed versions for the given path
func ListVersions(path string) (vers []string) {
	vers = []string{}

	folders, _ := ioutil.ReadDir(path)
	length := len(folders)

	for i := length - 1; i > -1; i-- {
		name := folders[i].Name()

		if name == "current" {
			continue
		}

		vers = append(vers, name)
	}

	return
}

// Symlink removes the link if file already
// present and the creates another symlink
func Symlink(current, base string) (err error) {
	// Remove symlink just in case it's already present
	err = os.RemoveAll(current)
	if err != nil {
		return errors.New(err)
	}

	// Set up new symlink
	err = os.Symlink(base, current)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
