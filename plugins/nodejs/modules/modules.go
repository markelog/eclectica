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

type Modules struct {
	current  string
	previous string
}

func New(previous, current string) *Modules {
	return &Modules{
		current:  current,
		previous: previous,
	}
}

func (modules Modules) Install() (err error) {
	err = modules.copy()
	if err != nil {
		return
	}

	if modules.SameMajors() == false {
		return
	}

	err = modules.rebuild()
	if err != nil {
		return errors.New(err)
	}

	return
}

func (modules Modules) rebuild() (err error) {
	dest := filepath.Join(modules.getDest(), "!(npm|yarn)")
	output, err := exec.Command("npm", "rebuild", "--global", "--verbose", dest).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}

	return
}

func (modules Modules) getDest() string {
	return filepath.Dir(modules.Path(modules.current))
}

func (modules Modules) copy() (err error) {
	packages := modules.Path(modules.previous)
	dest := modules.getDest()

	err = cprf.Copy(packages, dest)
	if err != nil {
		return errors.New(err)
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

func (modules Modules) Path(version string) string {
	path := variables.Path("node", version)

	return filepath.Join(path, "lib/node_modules")
}

func (modules Modules) SameMajors() bool {
	previous, _ := semver.Make(modules.previous)
	current, _ := semver.Make(modules.current)

	return previous.Major == current.Major
}
