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
)

func Prefix(name string) string {
	if HasLocalBin() {
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

func NeedToRestartShell() bool {
	return HasLocalBin() == false && os.Getenv("ECLECTICA") != "true"
}

func HasLocalBin() bool {
	path := os.Getenv("PATH")
	home := os.Getenv("HOME")
	localPath := home + "/bin"

	return strings.Contains(path, localPath)
}
