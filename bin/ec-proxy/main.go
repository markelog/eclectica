package main

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"

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
		env = append(env, environment...)

		cmd.Env = env
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
}

// Get relative path to the dot file
func getRelativePath(dotPath string) string {
	cwd, err := os.Getwd()
	print.Error(err)

	dir, err := filepath.Rel(filepath.Dir(cwd), cwd)
	print.Error(err)

	return "./" + filepath.Join(dir, filepath.Base(dotPath))
}

func main() {
	_, name := path.Split(os.Args[0])

	language := plugins.SearchBin(name)
	dotFiles := plugins.Dots(language)
	base := variables.Home()

	version, dotPath, err := io.GetVersion(dotFiles)
	print.Error(err)

	pathPart := filepath.Join(base, language, version)
	binPath := filepath.Join(pathPart, "bin", name)

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		err = errors.New("Version \"" + version + "\" was defined on \"" +
			getRelativePath(dotPath) + "\" path but this version is not installed")
		print.Error(err)
	}

	args := []string{binPath}
	args = append(args, os.Args[1:]...)

	cmd := console.Get(args)

	setCmd(cmd, language, version)

	err = cmd.Run()
	if sysErr, ok := err.(*exec.ExitError); ok {
		if status, ok := sysErr.Sys().(syscall.WaitStatus); ok {
			os.Exit(status.ExitStatus())
		} else {
			os.Exit(1)
		}
	}
}
