package python

import (
	"os"
	"os/exec"
	"strings"

	"github.com/markelog/eclectica/cmd/print"
)

var (
	OSXDependencies = []string{
		"openssl", "readline",
	}
	XCodeDependencies = []string{
		"xcrun", "make", "gcc",
	}
)

func checkXCodeDependencies() bool {
	for _, dep := range XCodeDependencies {
		_, err := exec.Command(dep, "--version").CombinedOutput()

		if err != nil {
			return false
		}
	}

	return true
}

func checkOSXDependencies() (has bool, deps []string, err error) {
	out, err := exec.Command("brew", "list").Output()
	if err != nil {
		return
	}

	output := string(out)

	for _, dep := range OSXDependencies {
		if strings.Contains(output, dep) == false {
			deps = append(deps, dep)
		}
	}

	if len(deps) > 0 {
		has = true
	}

	return
}

func printErrForOSXDependencies(deps []string) {
	start := `Python cannot be installed without external dependencies,
  please execute following command before trying it again (you need to do it only`
	middle := " once"
	end := "):"
	command := "brew update && brew install " + strings.Join(deps, " ")

	print.Install(start, middle, end, command)
	os.Exit(1)
}

func printErrForXCodeDependencies() {
	start := `Python cannot be installed without external dependencies,
  like Xcode, please download it from https://developer.apple.com/download/
  before trying it again (you need to do it only`
	middle := " once"
	end := ")"
	command := ""

	print.Install(start, middle, end, command)
	os.Exit(1)
}

func dealWithOSXShell() error {
	has := checkXCodeDependencies()

	if has == false {
		printErrForXCodeDependencies()
		return nil
	}

	has, deps, err := checkOSXDependencies()

	if err != nil {
		return err
	}

	if has == false {
		return nil
	}

	printErrForOSXDependencies(deps)

	return nil
}
