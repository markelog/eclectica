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

func setCmd(cmd *exec.Cmd, name, version string) {
	plugin := plugins.New(name)

	environment, err := plugin.Pkg.Environment(version)
	print.Error(err)

	env := os.Environ()
	env = append(env, environment)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Env = env
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
