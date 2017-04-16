package variables

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	NonInstallCommands = []string{"ls", "rm", "version", "init", "--help", "-h"}
	DefaultInstall     = filepath.Join(Base(), "bin")
)

func Prefix(name string) string {
	return filepath.Join(Home(), name)
}

// TempDir gets OS consistent folder path
// I am crying over here :/
func TempDir() (tmp string) {
	tmp = os.TempDir()
	if runtime.GOOS == "linux" {
		tmp += "/"
	}

	return
}

// IsDebug checks if eclectica in the debug state
// i.e. will print more info when executing commands
func IsDebug() bool {
	return os.Getenv("DEBUG") == "true"
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

// Path gives full path to parent of "bin" folder
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
	path := GetShellPath()
	parts := strings.Split(path, "/")

	return parts[len(parts)-1]
}

func GetShellPath() string {
	path := os.Getenv("SHELL")

	if len(path) == 0 {
		return "/bin/bash"
	}

	return path
}

func Base() string {
	return filepath.Join(os.Getenv("HOME"), ".eclectica")
}

func Home() string {
	return filepath.Join(Base(), "versions")
}
