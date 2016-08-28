package ruby

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/variables"
)

var (
	dependencies = []string{
		"autoconf", "bison", "build-essential",
		"libssl-dev", "libyaml-dev", "libreadline6-dev",
		"zlib1g-dev", "libncurses5-dev", "libffi-dev",
		"libgdbm3", "libgdbm-dev",
	}
)

func checkDependencies() (has bool, deps []string, err error) {
	if runtime.GOOS != "linux" {
		return
	}

	out, err := exec.Command("dpkg", "-l").Output()
	if err != nil {
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

func dealWithShell() (bool, error) {
	has, deps, err := checkDependencies()

	if err != nil {
		return false, err
	}

	if has == false {
		return true, nil
	}

	start := `Ruby has been installed, but it requires global dependencies which weren't found on your system,
  please execute following command to complete installation (you would need to do it only`
	middle := " once"
	end := "):"
	command := "sudo apt-get install " + strings.Join(deps, " ")

	if variables.NeedToRestartShell("ruby") {
		command += " && " + variables.GetShellName()
	}

	print.PostInstall(start, middle, end, command)

	return false, nil
}
