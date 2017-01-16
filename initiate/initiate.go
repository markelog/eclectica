package initiate

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

type Init struct {
	language           string
	plugins            []string
	shouldRestartShell bool
}

func New(language string, plugins []string) *Init {
	init := &Init{
		language:           language,
		plugins:            plugins,
		shouldRestartShell: false,
	}

	return init
}

func (init *Init) CheckShell() {
	init.shouldRestartShell = init.needRestartShell()
}

func (init *Init) Initiate() (err error) {
	_, err = io.CreateDir(variables.DefaultInstall)
	if err != nil {
		return
	}

	command := init.composeCommand()

	err = rc.New(command).Add()
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

func (init *Init) needRestartShell() bool {
	if strings.Contains(os.Getenv("PATH"), variables.DefaultInstall) == false {
		return true
	}

	ecPath := fmt.Sprintf("/%s/%s", ".eclectica/bin", init.language)
	output, _ := exec.Command("type", init.language).Output()
	out := string(output)

	if strings.Contains(out, ecPath) == false {
		return true
	}

	return false
}

func (init *Init) composeCommand() string {
	result := "# Eclectic stuff\n"

	// For golang
	result += "export PATH=" + variables.Path("go") + "/bin:$PATH\n"

	// For shared modules
	shared := filepath.Join(variables.Home(), "shared")
	result += "export PATH=" + shared + "/bin:$PATH\n"

	for _, language := range init.plugins {
		result += "export PATH=" +
			filepath.Join(variables.Home(), language, "current/bin") + ":$PATH\n"
	}

	result += "export PATH=" + variables.DefaultInstall + ":$PATH\n"

	return result
}
