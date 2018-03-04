// Package shell works with shell for eclectica
package shell

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/rc"
	"github.com/markelog/eclectica/variables"
)

const (
	command = `

#eclectica start
export PATH="$(ec path)"
#eclectica end

`
)

// Shell essential structure
type Shell struct {
	command       string
	shouldRestart bool
	plugins       []string
	rc            *rc.Rc
}

// New cretates new Shell struct
func New(plugins []string) *Shell {
	shell := &Shell{
		command:       command,
		shouldRestart: false,
		plugins:       plugins,
		rc:            rc.New(command),
	}

	return shell
}

// Check if we need to restart the shell
func (shell *Shell) Check() {
	shell.shouldRestart = shell.checkStatus()
}

// Initiate the shell if needed
func (shell *Shell) Initiate() (err error) {
	_, err = io.CreateDir(variables.DefaultInstall)
	if err != nil {
		return
	}

	err = shell.rc.Add()
	if err != nil {
		return
	}

	return
}

func (shell *Shell) Remove() (err error) {
	err = shell.rc.Remove()
	if err != nil {
		return
	}

	return
}

// Start starts the shell if needed
func (shell *Shell) Start() bool {
	if shell.shouldRestart {
		return Start()
	}

	return false
}

// checkStatus checks the status of the shell
func (shell *Shell) checkStatus() bool {
	cmd := Compose(shell.plugins)
	if strings.Contains(os.Getenv("PATH"), cmd) == false {
		return true
	}

	return false
}

// Compose returns $PATH paths for all provided languages
func Compose(plugins []string) (result string) {
	// First eclectica binaries
	result = ":" + variables.DefaultInstall

	for _, language := range plugins {
		result += ":" + filepath.Join(variables.Home(), language, "current/bin")
	}

	return
}

// Name name of the current shell
func Name() string {
	path := Path()
	parts := strings.Split(path, "/")

	return parts[len(parts)-1]
}

// Path gets path to current shell binary
func Path() string {
	path := os.Getenv("SHELL")

	if len(path) == 0 {
		return "/bin/bash"
	}

	return path
}

// Start starts a shell
// This is the only place beside cmd modules where we
// might output stuff to std(out | err)
func Start() bool {

	// If input is not a terminal - get out
	if terminal.IsTerminal(int(os.Stdout.Fd())) == false {
		return false
	}

	var procAttr os.ProcAttr

	procAttr.Files = []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}

	args := []string{
		Name(),
	}

	proc, err := os.StartProcess(Path(), args, &procAttr)
	print.Error(err)

	print.Green(
		"First time executing eclectica - had to restart the shell",
	)

	_, err = proc.Wait()
	print.Error(err)

	return true
}
