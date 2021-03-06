// Package modules is a way for installation node modules for node.js plugin
package modules

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/go-errors/errors"
	"github.com/markelog/cprf"
	"github.com/markelog/eclectica/variables"
)

// Modules essential struct
type Modules struct {
	current  string
	previous string
}

// New returns modules struct
func New(previous, current string) *Modules {
	return &Modules{
		current:  current,
		previous: previous,
	}
}

// Install modules
func (modules Modules) Install() (err error) {
	if modules.SameMajors() {
		err = modules.copy()
		return
	}

	err = modules.reinstall()
	if err != nil {
		return errors.New(err)
	}

	return
}

func (modules Modules) read() (result []string, err error) {
	dest := modules.getPrevious()

	files, err := ioutil.ReadDir(dest)

	if err != nil {
		err = errors.New(err)
		return
	}

	for _, file := range files {
		name := file.Name()

		if name == "npm" {
			continue
		}

		if name == "yarn" {
			continue
		}

		result = append(result, filepath.Join(dest, name))
	}

	return
}

func (modules Modules) reinstall() (err error) {
	install, err := modules.read()
	if err != nil {
		return
	}

	if len(install) == 0 {
		return
	}

	install = append(
		[]string{
			"install",
			"--offline",
			"--global",
			"--verbose",
		},
		install...,
	)
	output, err := exec.Command("npm", install...).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return
}

func (modules Modules) getDest() string {
	return modules.Path(modules.current)
}

func (modules Modules) getPrevious() string {
	return modules.Path(modules.previous)
}

func (modules Modules) copy() (err error) {
	dest := modules.getDest()

	install, err := modules.read()
	if err != nil {
		return
	}

	if len(install) == 0 {
		return
	}

	for _, pack := range install {
		err = cprf.Copy(pack, dest)
		if err != nil {
			return errors.New(err)
		}
	}

	previousBin := filepath.Join(variables.Path("node", modules.previous), "bin")
	currentBin := filepath.Join(variables.Path("node", modules.current), "bin")

	bins, _ := ioutil.ReadDir(previousBin)
	for _, bin := range bins {
		name := bin.Name()

		if name == "node" {
			continue
		}

		if name == "npm" {
			continue
		}

		if name == "yarn" {
			continue
		}

		linkPath := filepath.Join(previousBin, name)
		newLink := filepath.Join(currentBin, name)

		_, statErr := os.Stat(newLink)
		if statErr == nil {
			continue
		}

		link, errLink := os.Readlink(linkPath)
		if errLink != nil {
			return errLink
		}

		newBin := filepath.Join(currentBin, name)

		removeErr := os.RemoveAll(newBin)
		if removeErr != nil {
			return removeErr
		}

		symErr := os.Symlink(link, newBin)
		if symErr != nil {
			return
		}
	}

	return
}

// Path returns path to the node_modules folder
func (modules Modules) Path(version string) string {
	path := variables.Path("node", version)

	return filepath.Join(path, "lib/node_modules")
}

// SameMajors compares previous and current node version
// returning true if they have same major versions
func (modules Modules) SameMajors() bool {
	previous, _ := semver.Make(modules.previous)
	current, _ := semver.Make(modules.current)

	return previous.Major == current.Major
}
