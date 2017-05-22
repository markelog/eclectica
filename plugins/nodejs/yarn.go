package nodejs

import (
	"errors"
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
	minimalForYarn, _ = semver.Make("4.0.0")
	yarnUrl           = "https://yarnpkg.com/latest.tar.gz"
)

func (node Node) download(path string) (err error) {
	_, err = grab.Get(path, yarnUrl)
	if err != nil {
		return
	}

	return
}

func (node Node) isYarnPossible() bool {
	version, _ := semver.Make(node.Version)

	return version.GTE(minimalForYarn)
}

func (node Node) Yarn() (ok bool, err error) {
	path := variables.Path("node", node.Version)
	modules := filepath.Join(path, "lib/node_modules")
	dist := filepath.Join(modules, "yarn-archive")
	temp := filepath.Join(modules, "yarn-temp")
	from := filepath.Join(temp, "dist/")
	dest := filepath.Join(modules, "yarn")

	if node.isYarnPossible() == false {
		return true, errors.New("\"" + node.Version + "\" version is not supported by yarn")
	}

	err = node.download(dist)
	if err != nil {
		return
	}

	err = archive.Extract(dist, temp)
	if err != nil {
		return
	}

	err = cprf.Copy(from+"/", dest)
	if err != nil {
		return
	}

	os.RemoveAll(dist)
	os.RemoveAll(temp)

	current := filepath.Join(modules, "yarn/bin/yarn.js")
	base := filepath.Join(path, "bin/yarn")

	err = io.Symlink(base, current)
	if err != nil {
		return
	}

	return true, nil
}
