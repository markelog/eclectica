package golang

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/rc"
	"github.com/markelog/eclectica/variables"
)

var (
	goRootCommand = `
export GOROOT=%s
`
	goPathCommand = `
export GOPATH=%s
`
	gopath = filepath.Join(os.Getenv("HOME"), "gocode")
	goroot = filepath.Join(variables.Prefix("go"), "go")
)

func dealWithRc() (bool, error) {
	goRootCommand := fmt.Sprintf(goRootCommand, goroot)
	goPathCommand := fmt.Sprintf(goPathCommand, gopath)

	err := rc.New(goRootCommand).Add()
	if err != nil {
		return false, err
	}

	gopath := rc.New("export GOPATH")
	if os.Getenv("GOPATH") == "" || gopath.Exists() == false {
		err := rc.New(goPathCommand).Add()
		if err != nil {
			return false, err
		}
	}

	return dealWithShell()
}

func dealWithShell() (bool, error) {
	if os.Getenv("GOROOT") == goroot {
		return true, nil
	}

	partStart := `Go has been installed, but it requires to restart your shell,`

	if os.Getenv("GOPATH") == "" {
		partStart = `Go has been installed, but it requires you to set $GOPATH environment variable.
Eclectica preset it for you to %s path,`
		partStart = fmt.Sprintf(partStart, gopath)
	}

	start := partStart + ` for this to take affect you need
to execute following command (you would need to do it only`
	middle := " once"
	end := "):"
	command := "exec " + variables.GetShellName()

	print.PostInstall(start, middle, end, command)

	return false, nil
}
