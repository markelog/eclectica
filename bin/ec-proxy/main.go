package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"reflect"
	"syscall"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/variables"
)

func getCmd(args []string) *exec.Cmd {
	fn := reflect.ValueOf(exec.Command)
	rargs := make([]reflect.Value, len(args))

	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}

	cmd := fn.Call(rargs)[0].Interface().(*exec.Cmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd
}

func getVersion(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ""
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

	return ""
}

func main() {
	_, name := path.Split(os.Args[0])

	pwd, err := os.Getwd()
	print.Error(err)

	versionPath := filepath.Join(pwd, fmt.Sprintf(".%s-version", name))
	version := getVersion(versionPath)

	path := variables.GetBin(name, version)

	args := []string{path}
	args = append(args, os.Args[1:]...)

	cmd := getCmd(args)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}
