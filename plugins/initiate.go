package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/markelog/eclectica/console"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/rc"
	"github.com/markelog/eclectica/variables"
)

func composeCommand() string {
	result := "# Eclectic stuff\n"

	for _, language := range Plugins {
		result += "export PATH=" +
			filepath.Join(variables.Home(), language, "current/bin") + ":$PATH\n"
	}

	result += "export PATH=" + variables.DefaultInstall + ":$PATH\n"

	return result
}

func Initiate() (err error) {
	_, err = io.CreateDir(variables.DefaultInstall)
	if err != nil {
		return err
	}

	command := composeCommand()

	return rc.New(command).Add()
}

func pathShell(language string) bool {
	if strings.Contains(os.Getenv("PATH"), variables.DefaultInstall) == false {
		console.Shell()
		return true
	}

	return false
}

func hashShell(language string) bool {
	output, _ := exec.Command("hash").Output()
	out := string(output)

	ecPath := fmt.Sprintf("/%s/%s", ".eclectica/bin", language)
	binPath := fmt.Sprintf("/%s/%s", "bin", language)

	if strings.Contains(out, binPath) && strings.Contains(out, ecPath) == false {
		console.Shell()
		return true
	}

	return false
}

func StartShell(language string) {
	pathResult := pathShell(language)

	if pathResult {
		return
	}

	hashShell()
}
