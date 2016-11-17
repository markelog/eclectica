package python

import (
	"os"
	"os/exec"
	"strings"

	"github.com/markelog/eclectica/cmd/print"
)

var (
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
		return err
	}

	if has == false {
		return nil
	}

	start := `Python cannot be installed without external LinuxDependencies,
  please execute following command before trying it again (you need to do it only`
	middle := " once"
	end := "):"
	command := "sudo apt-get update && sudo apt-get install -y " + strings.Join(deps, " ")

	print.Install(start, middle, end, command)
	os.Exit(1)

	return nil
}
