package rvm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/markelog/release"
)

var (

	// Right now lowest possible version on rvm is for "10.12"
	min = float64(10.12)
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
	floatVersion, _ := strconv.ParseFloat(version, 64)
	arch := "x86_64"

	if min < floatVersion {
		version = fmt.Sprintf("%v", min)
	}

	versions := strings.Split(version, ".")
	version = versions[0] + "." + versions[1]

	return fmt.Sprintf("%s/%s/%s/%s", versionLink, typa, version, arch)
}
