package nodejs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blang/semver"
	"github.com/markelog/archive"
	"github.com/markelog/cprf"
	"gopkg.in/cavaliercoder/grab.v1"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"
)

var (
	minimalForYarn, _ = semver.Make("6.0.0")
	version           = "1.13.0"
	yarnURL           = fmt.Sprintf(
		"https://yarnpkg.com/downloads/%s/yarn-v%s.tar.gz",
		version,
		version,
	)
)

func (node Node) download(path string) (err error) {
	_, err = grab.Get(path, yarnURL)
	if err != nil {
		return
	}

	return
}

func (node Node) isYarnPossible() bool {
	version, _ := semver.Make(node.Version)

	return version.GTE(minimalForYarn)
}

// Yarn does everything that needs to be done to install yarn
func (node Node) Yarn() (ok bool, err error) {
	var (
		path       = variables.Path("node", node.Version)
		modules    = filepath.Join(path, "lib/node_modules")
		archived   = filepath.Join(variables.TempDir(), "yarn-archived")
		unarchived = filepath.Join(variables.TempDir(), "yarn-unarchived")
		from       = filepath.Join(unarchived, fmt.Sprintf("yarn-v%s/", version))
		dest       = filepath.Join(modules, "yarn")
	)

	if node.isYarnPossible() == false {
		return true, errors.New("\"" + node.Version + "\" version is not supported by yarn")
	}

	if _, statErr := os.Stat(unarchived); statErr != nil {
		err = node.download(archived)
		if err != nil {
			return
		}

		err = archive.Extract(archived, unarchived)
		if err != nil {
			return
		}
	}

	err = cprf.Copy(from+"/", dest)
	if err != nil {
		return
	}

	os.RemoveAll(archived)

	current := filepath.Join(modules, "yarn/bin/yarn.js")
	base := filepath.Join(path, "bin/yarn")

	err = io.Symlink(base, current)
	if err != nil {
		return
	}

	return true, nil
}
