package rc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/markelog/eclectica/variables"
)

var (
	rcs = []string{
		".bash_profile", ".bashrc", ".profile",
		".zshrc", ".fishrc",
	}
)

type Rc struct {
	command string
	path    string
}

func New(command string) *Rc {
	rc := &Rc{
		command: command,
		path:    "",
	}

	rc.path = rc.Find()

	return rc
}

func (rc *Rc) Add() error {
	shell := variables.GetShellName()

	if shell == "bash" && runtime.GOOS == "linux" {
		bashrc := &Rc{
			command: rc.command,
			path:    filepath.Join(os.Getenv("HOME"), ".bashrc"),
		}

		bashProfile := &Rc{
			command: rc.command,
			path:    filepath.Join(os.Getenv("HOME"), ".bash_profile"),
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

	return rc.add()
}

func (rc *Rc) add() error {
	if rc.Exists() {
		return nil
	}

	file, err := os.OpenFile(rc.path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = file.WriteString(rc.command)
	return err
}

func (rc *Rc) Exists() bool {
	contents, err := ioutil.ReadFile(rc.path)
	if err != nil {
		return false
	}

	str := string(contents)

	return strings.Contains(str, rc.command)
}

func (rc *Rc) Find() string {
	home := os.Getenv("HOME")

	files, _ := ioutil.ReadDir(home)

	for _, possibility := range rcs {
		for _, file := range files {
			if file.Name() == possibility {
				return filepath.Join(home, possibility)
			}
		}
	}

	return ""
}
