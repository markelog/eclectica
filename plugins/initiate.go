package plugins

import (
	"path/filepath"

	// "github.com/markelog/eclectica/console"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/rc"
	"github.com/markelog/eclectica/variables"
)

func composeCommand(languages []string) string {
	result := "# Eclectic stuff\n"

	for _, language := range languages {
		result += "export PATH=" +
			filepath.Join(variables.Home(), language, "current/bin") + ":$PATH\n"
	}

	result += "export PATH=" + variables.DefaultInstall + ":$PATH\n"

	return result
}

func Initiate(languages []string) (err error) {
	_, err = io.CreateDir(variables.DefaultInstall)
	if err != nil {
		return err
	}

	command := composeCommand(languages)

	return rc.New(command).Add()
}
