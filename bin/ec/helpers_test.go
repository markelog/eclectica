package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/markelog/eclectica/variables"
)

var (
	path, _ = filepath.Abs("./main.go")
	bins    = variables.DefaultInstall
)

func shouldRun(langauge string) bool {
	if os.Getenv("TEST_ALL") == "true" {
		return true
	}

	if os.Getenv("TEST_LANGUAGE") == langauge {
		return true
	}

	return false
}

func getCmd(args []interface{}) *exec.Cmd {
	fn := reflect.ValueOf(exec.Command)
	rargs := make([]reflect.Value, len(args))

	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}

	cmd := fn.Call(rargs)[0].Interface().(*exec.Cmd)

	return cmd
}

func Command(args ...interface{}) *exec.Cmd {
	return getCmd(args)
}

func Execute(args ...interface{}) *exec.Cmd {
	cmd := getCmd(args)

	// Output result for testing purposes
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	return cmd
}

func Kill(cmd *exec.Cmd) {
	cmd.Process.Kill()
}

func checkRemoteList(name, mask string, timeout int) bool {
	cmd := Command("go", "run", path, "ls", "-r", name)
	output := &bytes.Buffer{}
	cmd.Stdout = output
	result := false
	proceed := true

	timer := time.AfterFunc(time.Duration(timeout)*time.Second, func() {
		Kill(cmd)
		proceed = false
	})

	go func() {
		for {
			out := string(output.Bytes())
			result = strings.Contains(out, mask)

			if result {
				timer.Stop()
				Kill(cmd)
				proceed = false
				return
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	cmd.Start()

	for proceed {
		time.Sleep(200 * time.Millisecond)
	}

	return result
}

func checkRemoteUse() (result string) {
	cmd := Command("go", "run", path, "-r")
	output := &bytes.Buffer{}
	cmd.Stdout = output
	proceed := true

	go func() {
		for {
			result = string(output.Bytes())

			if len(result) > 0 {
				Kill(cmd)
				proceed = false
				return
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	cmd.Start()

	for proceed {
		time.Sleep(200 * time.Millisecond)
	}

	return result
}

func checkRemoteUseWithLanguage(name string) (result string) {
	cmd := Command("go", "run", path, "-r", "go")
	output := &bytes.Buffer{}
	cmd.Stdout = output
	proceed := true

	go func() {
		for {
			result = string(output.Bytes())

			if len(result) > 0 && strings.Contains(result, "Mask") {
				Kill(cmd)
				proceed = false
				return
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	cmd.Start()

	for proceed {
		time.Sleep(200 * time.Millisecond)
	}

	return result
}
