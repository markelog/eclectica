package plugins

import (
	"fmt"

	"github.com/markelog/eclectica/rc"
	"github.com/markelog/eclectica/variables"
)

var (
	command = `
# Eclectic stuff
export PATH=%s/bin:$PATH
export ECLECTICA=true
`
)

func Initiate(name string) (err error) {
	command := fmt.Sprintf(command, variables.DefaultInstall)

	if name == "rust" {
		name = "rustc"
	}

	if variables.ShouldBeLocalBin(name) {
		return
	}

	return rc.New(command).Add()
}
