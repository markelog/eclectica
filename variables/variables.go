package variables

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	NonInstallCommands = []string{"ls", "rm", "version", "init", "--help", "-h"}
	Files              = [4]string{"bin", "lib", "include", "share"}
	DefaultInstall     = filepath.Join(os.Getenv("HOME"), ".eclectica/bin")
)

func Prefix(name string) string {
	return filepath.Join(Home(), name)
}

func ExecutablePath(name string) string {
	return DefaultInstall
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

func Index(path string, value string) int {
	value = filepath.Join(value, "bin")
	bins := strings.Split(path, ":")

	for index, element := range bins {
		if element == value {
			return index
		}
	}

	return -1
}

func InLocalBin(path, local, name string) bool {
	bins := strings.Split(path, ":")
	localIndex := Index(path, local)

	// If local path is not present in the $PATH
	if localIndex == -1 {
		return false
	}

	for index, element := range bins {

		// Exclude local path
		if element == local {
			continue
		}

		bin := filepath.Join(element, name)

		// Compare executable positions in the $PATH, if specific binary is present
		if _, err := os.Stat(bin); err == nil {
			if index < localIndex {
				return false
			}
		}
	}

	return true
}

func ShouldBeLocalBin(name string) bool {
	return InLocalBin(os.Getenv("PATH"), os.Getenv("HOME"), name)
}
