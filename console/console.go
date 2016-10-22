package console

import (
	"os"
	"os/exec"
	"reflect"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/variables"
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

// Start Shell
func Shell() {
	var procAttr os.ProcAttr

	procAttr.Files = []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}

	args := []string{
		variables.GetShellName(),
	}

	proc, err := os.StartProcess(os.Getenv("SHELL"), args, &procAttr)
	print.Error(err)

	_, err = proc.Wait()
	print.Error(err)
}
