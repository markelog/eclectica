package modules

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

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

func (modules Modules) List() (packages []string, err error) {
	path := filepath.Dir(modules.Path(modules.previous))
	prefix := "--prefix=" + path

	bytes, _ := exec.Command("npm", "list", prefix, "--depth=0").Output()
	output := string(bytes)

	rPackages, _ := regexp.Compile(`\s(\w+)@`)
	tmp := rPackages.FindAllStringSubmatch(output, -1)

	for _, pack := range tmp {
		name := pack[1]

		if name == "npm" {
			continue
		}

		if name == "yarn" {
			continue
		}

		packages = append(packages, name)
	}

	return
}

func (modules Modules) Install() (err error) {
	if modules.SameMajors() {
		return modules.Copy()
	}

	packages, err := modules.List()
	if err != nil {
		return errors.New(err)
	}

	for _, name := range packages {
		err = modules.install(name)
		if err != nil {
			return errors.New(err)
		}
	}

	return
}

func (modules Modules) install(name string) (err error) {
	// Doesn't work with yarn 0.22.x and 0.23.x
	// if node.isYarnPossible() {
	// bin := variables.Path("node", modules.version)
	// _, err = exec.Command("yarn", "global", "add", name, "--prefix", bin).CombinedOutput()
	//
	// return
	// }

	_, err = exec.Command("npm", "install", "--global", name).CombinedOutput()

	return
}

func (modules Modules) Copy() (err error) {
	packages := modules.Path(modules.previous)
	dest := filepath.Dir(modules.Path(modules.current))

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
