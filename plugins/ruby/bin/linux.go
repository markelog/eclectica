// Package bin provides ruby related methods for the "bin" type installation
package bin

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-errors/errors"
	"github.com/mgutz/ansi"

	"github.com/markelog/eclectica/cmd/print"
)

var (
	dependencies = []string{
		"autoconf", "bison", "build-essential",
		"libssl-dev", "libyaml-dev", "libreadline6-dev",
		"zlib1g-dev", "libncurses5-dev", "libffi-dev",
		"libgdbm3", "libgdbm-dev", "libpq-dev",
	}
)

func checkDependencies() (has bool, deps []string, err error) {
	if runtime.GOOS != "linux" {
		return
	}

	out, err := exec.Command("dpkg", "-l").Output()
	if err != nil {
		err = errors.New(err)
		return
	}

	output := string(out)

	for _, dep := range dependencies {
		if strings.Contains(output, dep) == false {
			deps = append(deps, dep)
		}
	}

	if len(deps) > 0 {
		has = true
	}

	return
}

func dealWithShell() error {
	has, deps, err := checkDependencies()

	if err != nil {
		return errors.New(err)
	}

	if has == false {
		return nil
	}

	message := `Ruby has been installed, but it requires global dependencies which weren't found on your system,
  please execute following command to complete installation (you need to do it only ` + ansi.Color("once", "red") + "):"
	command := "sudo apt-get update && sudo apt-get install -y " + strings.Join(deps, " ")

	print.Warning(message, command)
	print.LastPrint()

	return nil
}
