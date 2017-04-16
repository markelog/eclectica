package nodejs

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/variables"
)

func (node Node) installModules(modules string) (err error) {
	folders, _ := ioutil.ReadDir(modules)

	for _, folder := range folders {
		name := folder.Name()

		if name == "npm" {
			continue
		}

		if name == "yarn" {
			continue
		}

		err = node.installModule(name)
		if err != nil {
			return
		}
	}

	return
}

func (node Node) installModule(name string) (err error) {
	// Doesn't work with yarn 0.22.x and 0.23.x
	// if node.isYarnPossible() {
	// bin := variables.Path("node", node.Version)
	// _, err = exec.Command("yarn", "global", "add", name, "--prefix", bin).CombinedOutput()
	//
	// return
	// }

	_, err = exec.Command("npm", "install", "--global", name).CombinedOutput()

	return
}

func (node Node) copyModules(modules string) (err error) {
	dest := node.modulesPath(node.Version)

	err = cprf.Copy(modules+"/", dest)
	if err != nil {
		return
	}

	previousBin := filepath.Join(variables.Path("node", node.previous), "bin")
	currentBin := filepath.Join(variables.Path("node", node.Version), "bin")

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

		if _, err := os.Stat(newLink); err == nil {
			continue
		}

		link, errLink := os.Readlink(linkPath)
		if errLink != nil {
			return errLink
		}

		err = os.Symlink(link, currentBin+"/"+name)
		if err != nil {
			return
		}
	}

	return
}

func (node Node) modulesPath(version string) string {
	path := variables.Path("node", version)

	return filepath.Join(path, "lib/node_modules")
}

func (node Node) sameMajors() bool {
	previous, _ := semver.Make(node.previous)
	current, _ := semver.Make(node.Version)

	return previous.Major == current.Major
}
