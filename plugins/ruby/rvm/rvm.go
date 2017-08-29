package rvm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/markelog/release"
)

// RemoveArtefacts removes RVM artefacts (ignore errors)
func RemoveArtefacts(base string) (err error) {
	gems := filepath.Join(base, "lib/ruby/gems")

	// Remove `cache` folder since it supposed to work with RVM cache
	folders, err := ioutil.ReadDir(gems)
	if err != nil {
		return
	}
	for _, folder := range folders {
		err = os.RemoveAll(filepath.Join(gems, folder.Name(), "cache"))
		if err != nil {
			return
		}
	}

	return nil
}

func GetUrl(versionLink string) string {
	typa, _, version := release.All()
	arch := "x86_64"

	versions := strings.Split(version, ".")
	version = versions[0] + "." + versions[1]

	return fmt.Sprintf("%s/%s/%s/%s", versionLink, typa, version, arch)
}
