package variables

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	NonInstallCommands = []string{"ls", "rm", "version", "init", "--help", "-h"}
	DefaultInstall     = filepath.Join(os.Getenv("HOME"), ".eclectica/bin")
)

func Prefix(name string) string {
	return filepath.Join(Home(), name)
}

// IsDebug checks if eclectica in the debug state
// i.e. will print more info when executing commands
func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
}

// IsCI if we running in CI environment
func IsCI() bool {
	return os.Getenv("CI") == "true"
}

func nameAndVersion(args []interface{}) (string, string) {
	var (
		name    = args[0].(string)
		version string
	)

	if len(args) == 2 {
		version = args[1].(string)
	} else {
		version = "current"
	}

	return name, version
}

func Path(args ...interface{}) string {
	name, version := nameAndVersion(args)

	return filepath.Join(Home(), name, version)
}

func GetBin(args ...interface{}) string {
	name, version := nameAndVersion(args)

	base := Path(name, version)

	// TODO: fix
	if name == "rust" {
		name = "rustc"
	}

	return filepath.Join(base, "bin", name)
}

func GetShellName() string {
	name := os.Getenv("SHELL")
	path := strings.Split(name, "/")

	return path[len(path)-1]
}

func Base() string {
	return filepath.Join(os.Getenv("HOME"), ".eclectica")
}

func Home() string {
	return filepath.Join(Base(), "versions")
}
