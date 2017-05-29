package initiate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/markelog/eclectica/console"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/rc"
	"github.com/markelog/eclectica/variables"
)

type Init struct {
	language           string
	plugins            []string
	command            string
	shouldRestartShell bool
	rc                 *rc.Rc
}

func New(language string, plugins []string) *Init {
	init := &Init{
		language: language,
		plugins:  plugins,
		command: `

#eclectica start
export PATH="$(ec path)"
#eclectica end

`,
		shouldRestartShell: false,
		rc:                 nil,
	}

	return init
}

func (init *Init) CheckShell() {
	init.shouldRestartShell = init.isShellRestart()
}

func (init *Init) Initiate() (err error) {
	_, err = io.CreateDir(variables.DefaultInstall)
	if err != nil {
		return
	}

	init.rc = rc.New(init.command)

	err = init.rc.Add()
	if err != nil {
		return
	}

	return
}

func (init *Init) RestartShell() {
	if init.shouldRestartShell {
		console.Shell()
	}
}

func (init *Init) isShellRestart() bool {
	if strings.Contains(os.Getenv("PATH"), Compose(init.plugins)) == false {
		return true
	}

	return false
}

func Compose(plugins []string) string {
	result := ""

	for _, language := range plugins {
		result += ":" + filepath.Join(variables.Home(), language, "current/bin")
	}

	// Eclectica binaries
	result += ":" + variables.DefaultInstall

	// For shared modules
	shared := filepath.Join(variables.Base(), "shared")
	result += ":" + filepath.Join(shared, "bin")

	return result
}
