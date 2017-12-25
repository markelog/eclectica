package compile

import (
	"os"
	"os/exec"
	"strings"

	"github.com/go-errors/errors"
	"github.com/mgutz/ansi"

	"github.com/markelog/eclectica/cmd/print"
)

var (

	// OSXDependencies is a list of OSX system dependencies
	OSXDependencies = []string{
		"openssl", "libyaml", "automake",
	}

	// XCodeDependencies is a list of XCode dependencies
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
		err = errors.New(err)
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
	message := `Ruby cannot be installed without external dependencies,
  please execute following command before trying it again (you need to do it only ` + ansi.Color("once", "red") + "):"
	command := "brew update && brew install " + strings.Join(deps, " ")

	print.Warning(message, command)
	os.Exit(1)
}

func printErrForXCodeDependencies() {
	message := `Ruby cannot be installed without Xcode,
	please download it from https://developer.apple.com/download/
  before trying it again (you need to do it only ` + ansi.Color("once", "red") + "):"

	print.Warning(message, "")
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
