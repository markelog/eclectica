package plugins

import (
	"os"
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

func StartShell() {
	if strings.Contains(os.Getenv("PATH"), variables.DefaultInstall) == false {
		console.Shell()
	}
}
