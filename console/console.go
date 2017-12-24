// Package console provides helpful console related methods
package console

import (
	"io"
	"io/ioutil"
	"os/exec"
	"reflect"
	"regexp"

	"github.com/go-errors/errors"
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

// Error is just facade for handling console pipes errors
func Error(err error, stdout, stderr io.ReadCloser) error {
	if stdout == nil || stderr == nil {
		return nil
	}

	strErr, errRead := ioutil.ReadAll(stdout)
	if err != nil {
		return errors.New(errRead)
	}
	strOut, errOut := ioutil.ReadAll(stderr)
	if err != nil {
		return errors.New(errOut)
	}

	if len(strErr) != 0 {
		str := string(strErr)
		return errors.New(str)
	}

	if len(strOut) != 0 {
		str := string(strOut)
		return errors.New(str)
	}

	// "Exit status" is just silly (not that following is much better)
	r, _ := regexp.Compile("^exit status")
	if r.MatchString(err.Error()) {
		err = errors.New("Unknown error :/")
	}

	return err
}
