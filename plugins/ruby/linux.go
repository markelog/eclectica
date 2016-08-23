package ruby

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

var (
	dependencies = []string{
		"autoconf", "bison", "build-essential",
		"libssl-dev", "libyaml-dev", "libreadline6-dev",
		"zlib1g-dev", "libncurses5-dev", "libffi-dev",
		"libgdbm3", "libgdbm-dev",
	}
)

func checkDependencies() (has bool, deps []string) {
	if runtime.GOOS != "linux" {
		return
	}

	out, _ := exec.Command("dpkg", "-l").Output()
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

func printMissingDependencies() {
	has, deps := checkDependencies()

	if has == false {
		return
	}

	messageStart := `Ruby has been installed, but it requires global dependencies which weren't found on your system,
  please execute following command to complete installation (you would need to do it only`
	messageMiddle := " once"
	messageEnd := "):"

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
	fmt.Println("                          sudo apt-get install", strings.Join(deps, " "))
	fmt.Println()
}
