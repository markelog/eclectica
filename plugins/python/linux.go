package python

import (
	"os"
	"os/exec"
	"strings"

	"github.com/go-errors/errors"
	"github.com/mgutz/ansi"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/variables"
)

var (

	// LinuxDependencies is a list of all linux system dependencies
	LinuxDependencies = []string{
		"make", "build-essential", "libssl-dev",
		"zlib1g-dev", "libbz2-dev", "libreadline-dev",
		"libsqlite3-dev", "llvm", "libncurses5-dev",
		"xz-utils",
	}
)

func checkLinuxDependencies() (has bool, deps []string, err error) {
	out, err := exec.Command("dpkg", "-l").Output()
	if err != nil {
		err = errors.New(err)
		return
	}

	output := string(out)

	for _, dep := range LinuxDependencies {
		if strings.Contains(output, dep) == false {
			deps = append(deps, dep)
		}
	}

	if len(deps) > 0 {
		has = true
	}

	return
}

func dealWithLinuxShell() error {
	has, deps, err := checkLinuxDependencies()

	if err != nil {
		return errors.New(err)
	}

	if has == false {
		return nil
	}

	message := `Python cannot be installed without external Linux dependencies,
  please execute following command before trying it again (you need to do it only ` + ansi.Color("once", "red") + "):"
	command := "sudo apt-get update && sudo apt-get install -y " + strings.Join(deps, " ")

	print.Warning(message, command)
	print.LastPrint()
	os.Exit(1)

	return nil
}

func getLinuxLineArguments(version string) []string {
	var (
		path      = variables.Path("python", version)
		prefix    = "--prefix=" + path
		ensurepip = "--with-ensurepip=upgrade"
	)

	result := []string{
		prefix, ensurepip,
	}

	return result
}
