package console

import (
	"os"
	"os/exec"
	"reflect"

	"github.com/markelog/eclectica/cmd/print"
)

// Pass array to exec.Command
func Get(args []string) *exec.Cmd {
	fn := reflect.ValueOf(exec.Command)
	rargs := make([]reflect.Value, len(args))

	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}

	cmd := fn.Call(rargs)[0].Interface().(*exec.Cmd)

	return cmd
}

// Start new shell
func Shell() {

	// Get the current working directory.
	cwd, err := os.Getwd()
	print.Error(err)

	// Transfer stdin, stdout, and stderr to the new process
	// and also set target directory for the shell to start in.
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   cwd,
	}

	proc, err := os.StartProcess(os.Getenv("SHELL"), []string{}, &pa)
	print.Error(err)

	// Wait until user exits the shell
	_, err = proc.Wait()
	print.Error(err)
}
