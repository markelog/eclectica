package rvm

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
