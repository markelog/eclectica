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

func Path(name, version string) string {
	if version == "" {
		version = "current"
	}

	return filepath.Join(Home(), name, version)
}

func GetBin(name, version string) string {

	// TODO: fix
	if name == "rust" {
		name = "rustc"
	}

	base := Path(name, version)

	return filepath.Join(base, "bin", name)
}

func GetShellName() string {
	name := os.Getenv("SHELL")
	path := strings.Split(name, "/")

	return path[len(path)-1]
}

func Home() string {
	return filepath.Join(os.Getenv("HOME"), ".eclectica/versions")
}
