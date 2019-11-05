package python

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/blang/semver"
	"github.com/go-errors/errors"
	"github.com/mgutz/ansi"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/release"
)

var (

	// OSXDependencies is a list of OSX system dependencies
	OSXDependencies = []string{
		"openssl", "readline",
	}

	// XCodeDependencies is a list of XCode dependencies
	XCodeDependencies = []string{
		"xcrun", "make", "gcc",
	}

	// Minimal version for --with-openssl flag,
	// without this flag python will no longer build with SSL :((
	minimalForWithOpenSSLFlag, _ = semver.Make("3.7.0")
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
	message := `Python cannot be installed without external dependencies,
  please execute following command before trying it again (you need to do it only ` + ansi.Color("once", "red") + "):"
	command := "brew update && brew install " + strings.Join(deps, " ")

	print.Warning(message, command)
	os.Exit(1)
}

func printErrForXCodeDependencies() {
	message := `Python cannot be installed without Xcode,
	please download it from https://developer.apple.com/download/
  before trying it again (you need to do it only ` + ansi.Color("once", "red") + "):"

	print.Warning(message, "")
	os.Exit(1)
}

func isWithOpenSSLFlag(version string) bool {
	madeVersion, _ := semver.Make(version)

	return madeVersion.GTE(minimalForWithOpenSSLFlag)
}

func getOSXLineArguments(version string) []string {
	var (
		path      = variables.Path("python", version)
		prefix    = "--prefix=" + path
		ensurepip = "--with-ensurepip=upgrade"
	)

	result := []string{
		prefix, ensurepip,
	}

	if isWithOpenSSLFlag(version) {
		result = append(result, "--with-openssl=/usr/local/opt/openssl")
	}

	return result
}

func getOSXEnvs(version string, original []string) []string {
	externals := []string{"readline"}

	if isWithOpenSSLFlag(version) == false {
		externals = append(externals, "openssl")
	}

	includeFlags := ""
	libFlags := ""

	for _, name := range externals {
		opt := "/usr/local/opt/"
		libFlags += "-L" + filepath.Join(opt, name, "lib") + " "
		includeFlags += "-I" + filepath.Join(opt, name, "include") + " "
	}

	// For zlib
	output, _ := exec.Command("xcrun", "--show-sdk-path").CombinedOutput()
	out := strings.TrimSpace(string(output))
	includeFlags += " -I" + filepath.Join(out, "/usr/include")

	original = append(original, "CFLAGS="+includeFlags)
	original = append(original, "LDFLAGS="+libFlags)

	// Since otherwise configure breaks for some versions :/
	original = append(original, "MACOSX_DEPLOYMENT_TARGET="+release.Version())

	return original
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
