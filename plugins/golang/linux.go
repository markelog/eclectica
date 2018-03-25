// Package golang provides additional golang methods for linux installation
package golang

import (
	"os/exec"
	"runtime"

	"github.com/mgutz/ansi"

	"github.com/markelog/eclectica/cmd/print"
)

func checkGit() (has bool) {
	text, _ := exec.Command("which", "git").Output()

	// So hard to extract exit code :/
	if len(text) == 0 {
		return false
	}

	return true
}

func dealWithShell() (err error) {
	var (
		command string
		message string
		has     = checkGit()
	)

	if has {
		return
	}

	message = `Golang has been installed, but it requires for git to be` +
		` installed also,
  (in order for ` + "`" + ansi.Color("go get", "blue") + "` to work) " +
		"see " + ansi.Color("https://golang.org/s/gogetcmd", "green") +
		` for more info.

  Run following command to resolve the issue (you need to do it only ` + ansi.Color("once", "red") + "):"

	if runtime.GOOS == "linux" {
		command = "sudo apt-get update && sudo apt-get install -y git"
	}

	if runtime.GOOS == "darwin" {
		command = "brew update && brew install git"
	}

	print.Warning(message, command)
	print.LastPrint()

	return
}
