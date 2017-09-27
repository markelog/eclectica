package console

import (
	"bytes"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/variables"
)

// Get gets cmd instance by passing array to exec.Command
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

	// If shell is not output - get out
	if terminal.IsTerminal(int(os.Stdout.Fd())) == false {
		return
	}

	var procAttr os.ProcAttr

	procAttr.Files = []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}

	args := []string{
		variables.GetShellName(),
	}

	proc, err := os.StartProcess(variables.GetShellPath(), args, &procAttr)
	print.Error(err)

	_, err = proc.Wait()
	print.Error(err)
}

// GetError is just facade to handling console errors
// Is stdOut param redudant?
func GetError(err error, stdErr, stdOut *bytes.Buffer) error {
	strErr := stdErr.String()
	strOut := stdOut.String()

	if len(strErr) != 0 {
		return errors.New(strErr)
	}

	if len(strOut) != 0 {
		return errors.New(trimMessage(strOut))
	}

	// "Exit status" is just silly (not that following is much better)
	r, _ := regexp.Compile("^exit status")
	if r.MatchString(err.Error()) {
		err = errors.New("Unknown error :/")
	}

	return err
}

func trimMessage(message string) (result string) {
	result = strings.TrimSpace(message)

	messageLength := 30

	messages := strings.Split(result, "\n")
	last := len(messages) - 1

	result = messages[last]

	if len(result) < messageLength {
		return
	}

	result = "..." + result[messageLength-3:]

	return
}
