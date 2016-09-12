package variables

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	Commands       = []string{"ls", "rm", "version", "init", "--help", "-h"}
	Files          = [4]string{"bin", "lib", "include", "share"}
	DefaultInstall = fmt.Sprintf("%s/.eclectica/install/go", os.Getenv("HOME"))
	DefaultBins    = "/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin:/usr/local/sbin"
)

func Prefix(name string) string {
	if ShouldBeLocalBin(name) {
		return os.Getenv("HOME")
	}

	return DefaultInstall
}

func GetShellName() string {
	name := os.Getenv("SHELL")
	path := strings.Split(name, "/")

	return path[len(path)-1]
}

func Home() string {
	return fmt.Sprintf("%s/.eclectica/versions", os.Getenv("HOME"))
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

		bin := element + "/" + name

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
