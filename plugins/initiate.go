package plugins

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/markelog/eclectica/variables"
)

var (
	command = `
# Eclectic stuff
export PATH=` + variables.DefaultInstall + `/bin:$PATH
export ECLECTICA="true"
`
	rcs = []string{
		".bash_profile", ".bashrc", ".profile",
		".zshrc", ".fishrc",
	}
)

func Initiate(name string) (err error) {
	if name == "rust" {
		name = "rustc"
	}

	if variables.ShouldBeLocalBin(name) {
		return
	}

	rc := findrc()
	if inShell(rc) {
		return
	}

	return add(rc)
}

func inShell(path string) bool {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	str := string(contents)

	return strings.Contains(str, `ECLECTICA="true"`)
}

func add(path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.WriteString(command)
	return err
}

func findrc() string {
	home := os.Getenv("HOME")

	files, _ := ioutil.ReadDir(home)

	for _, rc := range rcs {
		for _, file := range files {
			if file.Name() == rc {
				return filepath.Join(home, rc)
			}
		}
	}

	return ""
}
