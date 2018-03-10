// Package rc provides a bit of logic for .*.rc files
package rc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-errors/errors"

	"github.com/markelog/eclectica/variables"
)

const (
	begin   = `#eclectica start`
	end     = `#eclectica end`
	command = `
export PATH="$(ec path)"
`
)

var (
	reg = regexp.MustCompile(begin + "(?:[^\n]*\n+)+" + end)
	rcs = map[string][]string{
		"bash": {".bash_profile", ".bashrc", ".profile"},
		"zsh":  {".zshrc"},
	}
)

// Rc essential structure
type Rc struct {
	path string
}

// New returns new Rc struct
func New() *Rc {
	rc := &Rc{}
	rc.path = rc.Find()

	return rc
}

// getRcs gets rc instances
func (rc *Rc) getRcs() (err error, bashrc *Rc, bashProfile *Rc) {
	pathsRc := filepath.Join(os.Getenv("HOME"), ".bashrc")
	pathsProfile := filepath.Join(os.Getenv("HOME"), ".bash_profile")

	// Make sure we have those files
	if _, err := os.Stat(pathsRc); err != nil {
		return errors.New("couldn't find the " + pathsRc), nil, nil
		// return rc.remove()
	}
	if _, err := os.Stat(pathsProfile); err != nil {
		return errors.New("couldn't find the " + pathsProfile), nil, nil
	}

	bashrc = &Rc{
		path: pathsRc,
	}

	bashProfile = &Rc{
		path: pathsProfile,
	}

	return nil, bashrc, bashProfile
}

// Add bash configs on Unix system
// .bashrc works when you open new bash session (open terminal)
// .bash_profile is executed when you login
//
// So in order for our env variables to be
// consistently exposed we need to modify both of them
func (rc *Rc) Add() error {
	shell := variables.GetShellName()

	if shell != "bash" {
		return rc.add()
	}

	err, bashrc, bashProfile := rc.getRcs()
	if err != nil {
		return rc.add()
	}

	err = bashrc.add()
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

	_, err = file.WriteString(begin + command + end)
	if err != nil {
		err = errors.New(err)
	}

	return
}

// Remove bash configs on Unix system
func (rc *Rc) Remove() error {
	shell := variables.GetShellName()

	if shell != "bash" {
		return rc.remove()
	}

	err, bashrc, bashProfile := rc.getRcs()
	if err != nil {
		return rc.remove()
	}

	err = bashrc.remove()
	if err != nil {
		return err
	}

	err = bashProfile.remove()
	if err != nil {
		return err
	}

	return nil
}

// add helper method for Remove()
func (rc *Rc) remove() (err error) {
	if rc.Exists() == false {
		return
	}

	read, err := ioutil.ReadFile(rc.path)
	if err != nil {
		return errors.New(err)
	}

	replaced := reg.ReplaceAll(read, []byte(""))

	err = ioutil.WriteFile(rc.path, replaced, 0)
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

	return reg.MatchString(string(contents))
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
