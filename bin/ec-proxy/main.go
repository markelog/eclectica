package main

import (
	"errors"
	"fmt"
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
	"github.com/markelog/eclectica/versions"
)

// Pipe results of command execution to parent and
// pass environment variables from language plugin
func setCmd(cmd *exec.Cmd, language, version string) {
	environment, err := plugins.New(&plugins.Args{
		Language: language,
		Version:  version,
	}).Environment()
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

func getVersion(language string) (version, dotPath string) {
	plugin := plugins.New(&plugins.Args{
		Language: language,
	})
	dotFiles := plugin.Dots()

	version, dotPath, err := io.GetVersion(dotFiles)
	print.Error(err)

	if version == "current" {
		return
	}

	if versions.IsPartial(version) == false {
		return version, dotPath
	}

	vers := plugin.List()
	if len(vers) == 0 {
		notInstalled(version, dotPath)
	}

	found, err := versions.Latest(version, vers)
	if err != nil {
		notInstalled(version, dotPath)
	}

	return found, dotPath
}

func notInstalled(version, dotPath string) {
	var (
		start  = "version: \"" + version + "\" "
		ending = "path but this version is not installed"
	)

	// Different error message for the partial version
	if versions.IsPartial(version) {
		start = "mask: \"" + version + "\" "
		ending = "path but none of these versions were installed"
	}

	err := errors.New(start + "was defined on \"" + getRelativePath(dotPath) + "\" " + ending)

	print.Error(err)
}

func main() {
	_, name := path.Split(os.Args[0])

	language := plugins.SearchBin(name)
	version, dotPath := getVersion(language)
	base := variables.Home()

	pathPart := filepath.Join(base, language, version)
	binPath := filepath.Join(pathPart, "bin", name)

	if variables.IsDebug() {
		fmt.Println("bin path: " + binPath)
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		notInstalled(version, dotPath)
	}

	args := []string{binPath}
	args = append(args, os.Args[1:]...)

	cmd := console.Get(args)

	setCmd(cmd, language, version)

	err := cmd.Run()

	// Pass the exit code back
	if sysErr, ok := err.(*exec.ExitError); ok {
		if status, ok := sysErr.Sys().(syscall.WaitStatus); ok {
			os.Exit(status.ExitStatus())
		} else {
			os.Exit(1)
		}
	}
}
