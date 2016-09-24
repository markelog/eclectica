package variables

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	NonInstallCommands = []string{"ls", "rm", "version", "init", "--help", "-h"}
	Files              = [4]string{"bin", "lib", "include", "share"}
	DefaultInstall     = filepath.Join(os.Getenv("HOME"), ".eclectica/install")
	DefaultBins        = "/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:/usr/local/sbin"
)

func Prefix(name string) string {
	if ShouldBeLocalBin(name) {
		return filepath.Join(Home(), name)
	}

	return filepath.Join(DefaultInstall, name)
}

func GetBin(name, version string) string {
	if version == "" {
		version = "current"
	}

	return filepath.Join(Home(), name, version, "bin", name)
}

func GetShellName() string {
	name := os.Getenv("SHELL")
	path := strings.Split(name, "/")

	return path[len(path)-1]
}

func Home() string {
	return filepath.Join(os.Getenv("HOME"), ".eclectica/versions")
}

func NeedToRestartShell(name string) bool {
	return ShouldBeLocalBin(name) == false && os.Getenv("ECLECTICA") != "true"
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
