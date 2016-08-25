package variables

import (
	"fmt"
	"os"
	"strings"
)

var (
	Commands       = []string{"ls", "rm", "init", "--help", "-h"}
	Files          = [4]string{"bin", "lib", "include", "share"}
	DefaultInstall = fmt.Sprintf("%s/.eclectica/install", os.Getenv("HOME"))
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

	for _, element := range bins {

		// Exclude local path
		if element == local {
			continue
		}

		bin := element + "/" + name

		// Compare executable positions in the $PATH, if specific binary is present
		if _, err := os.Stat(bin); err == nil {
			if Index(path, element) < localIndex {
				return false
			}
		}
	}

	return true
}

func ShouldBeLocalBin(name string) bool {
	return InLocalBin(os.Getenv("PATH"), os.Getenv("HOME"), name)
}
