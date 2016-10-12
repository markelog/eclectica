package main

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/console"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/variables"
)

// Pipe results of command execution to parent and
// pass environment variables from language plugin
func setCmd(cmd *exec.Cmd, name, version string) {
	environment, err := plugins.New(name, version).Environment()
	print.Error(err)

	if len(environment) > 0 {
		env := os.Environ()
		env = append(env, environment)

		cmd.Env = env
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
}

func main() {
	_, name := path.Split(os.Args[0])

	language := plugins.SearchBin(name)
	base := variables.Home()

	version, err := io.GetVersion(language)
	print.Error(err)

	pathPart := filepath.Join(base, language, version)
	path := filepath.Join(pathPart, "bin", name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = errors.New("Version " + version + " has not been established")
		print.Error(err)
	}

	args := []string{path}
	args = append(args, os.Args[1:]...)

	cmd := console.Get(args)

	setCmd(cmd, language, version)

	cmd.Run()
}
