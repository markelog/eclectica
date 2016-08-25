package ruby

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"

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

func printMissingDependencies() error {
	has, deps, err := checkDependencies()

	if err != nil {
		return err
	}

	if has == false {
		return nil
	}

	messageStart := `Ruby has been installed, but it requires global dependencies which weren't found on your system,
  please execute following command to complete installation (you would need to do it only`
	messageMiddle := " once"
	messageEnd := "):"
	fullMessage := "                          sudo apt-get install " + strings.Join(deps, " ")

	if variables.HasLocalBin() == false {
		fullMessage += " && " + variables.GetShellName()
	}

	fmt.Println()

	color.Set(color.FgRed)
	fmt.Print("> ")
	color.Unset()

	color.Set(color.Bold)
	fmt.Print(messageStart)
	color.Set(color.FgRed)
	fmt.Print(messageMiddle)
	color.Unset()

	color.Set(color.Bold)
	fmt.Print(messageEnd)
	color.Unset()

	fmt.Println()
	fmt.Println()
	fmt.Println(fullMessage)
	fmt.Println()

	return nil
}
