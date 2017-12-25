// Package rc provides a bit of logic for .*.rc files
package rc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"

	"github.com/markelog/eclectica/variables"
)

var (
	rcs = map[string][]string{
		"bash": {".bash_profile", ".bashrc", ".profile"},
		"zsh":  {".zshrc"},
	}
)

// Rc essential structure
type Rc struct {
	command string
	path    string
}

// New returns new Rc struct
func New(command string) *Rc {
	rc := &Rc{}

	rc.command = command
	rc.path = rc.Find()

	return rc
}

// Add bash configs on Linux system
// .bashrc works when you open new bash session (open terminal)
// .bash_profile is executed when you login
//
// So in order for our env variables to be consistently exposed we need to modify both of them
// Note: on Mac, .bash_profile is executed when new bash session is opened,
// so we don't need to this in there
func (rc *Rc) Add() error {
	shell := variables.GetShellName()

	if shell != "bash" {
		return rc.add()
	}

	pathsRc := filepath.Join(os.Getenv("HOME"), ".bashrc")
	pathsProfile := filepath.Join(os.Getenv("HOME"), ".bash_profile")

	// Make sure we have those files
	if _, err := os.Stat(pathsRc); err != nil {
		return rc.add()
	}
	if _, err := os.Stat(pathsProfile); err != nil {
		return rc.add()
	}

	bashrc := &Rc{
		command: rc.command,
		path:    pathsRc,
	}

	bashProfile := &Rc{
		command: rc.command,
		path:    pathsProfile,
	}

	err := bashrc.add()
	if err != nil {
		return err
	}

	err = bashProfile.add()
	if err != nil {
		return err
	}

	return nil
}

// add helper method for Add()
func (rc *Rc) add() (err error) {
	if rc.Exists() {
		return
	}

	file, err := os.OpenFile(rc.path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()

	if err != nil {
		return errors.New(err)
	}

	_, err = file.WriteString(rc.command)
	if err != nil {
		err = errors.New(err)
	}

	return
}

// Exists checks if rc data already present
func (rc *Rc) Exists() bool {
	contents, err := ioutil.ReadFile(rc.path)
	if err != nil {
		return false
	}

	str := string(contents)

	return strings.Contains(str, rc.command)
}

// Find finds proper rc file
func (rc *Rc) Find() string {
	home := os.Getenv("HOME")
	shell := variables.GetShellName()

	files, _ := ioutil.ReadDir(home)

	for _, possibility := range rcs[shell] {
		for _, file := range files {
			if file.Name() == possibility {
				return filepath.Join(home, possibility)
			}
		}
	}

	return ""
}
