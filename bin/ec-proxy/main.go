package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/console"
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

func getVersion(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "current"
	}

	file, err := os.Open(path)
	print.Error(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		return scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		print.Error(err)
	}

	return "current"
}

func main() {
	_, name := path.Split(os.Args[0])

	language := plugins.SearchBin(name)
	base := variables.Home()

	pwd, err := os.Getwd()
	print.Error(err)

	versionPath := filepath.Join(pwd, fmt.Sprintf(".%s-version", language))
	version := getVersion(versionPath)

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
