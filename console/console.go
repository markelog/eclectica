package console

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"reflect"
	"regexp"

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

// GetError is just facade to handling a console errors
// Is stdOut param redudant?
func GetError(err error, stdErr, stdOut *bytes.Buffer) error {
	strErr := stdErr.String()
	strOut := stdOut.String()

	if len(strErr) != 0 {
		return errors.New(strErr)
	}

	if len(strOut) != 0 {
		return errors.New(strOut)
	}

	// "Exit status" is just silly (not that following is much better)
	r, _ := regexp.Compile("^exit status")
	if r.MatchString(err.Error()) {
		err = errors.New("Unknown error :/")
	}

	return err
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
