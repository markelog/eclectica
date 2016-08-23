package main_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"time"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/plugins"
)

var (
	path string
)

func init() {
	path, _ = filepath.Abs("./main.go")
}

func getCmd(args []interface{}) *exec.Cmd {
	fn := reflect.ValueOf(exec.Command)
	rargs := make([]reflect.Value, len(args))
	for i, a := range args {
		rargs[i] = reflect.ValueOf(a)
	}

	cmd := fn.Call(rargs)[0].Interface().(*exec.Cmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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
	pgid, _ := syscall.Getpgid(cmd.Process.Pid)
	syscall.Kill(-pgid, 15)
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

var _ = Describe("main", func() {
	if os.Getenv("TRAVIS") != "true" && os.Getenv("INT") != "true" {
		return
	}

	Describe("Rust", func() {
		BeforeEach(func() {
			fmt.Println()
			fmt.Println("Removing rust@1.9.0")
			Execute("go", "run", path, "rm", "rust@1.9.0")
		})

		It("should install rust 1.9.0", func() {
			Execute("go", "run", path, "rust@1.9.0")
			command, _ := Command("go", "run", path, "ls", "rust").Output()

			Expect(strings.Contains(string(command), "♥ 1.9.0")).To(Equal(true))
		})

		It("should list installed rust versions", func() {
			Execute("go", "run", path, "rust@1.9.0")
			command, _ := Command("go", "run", path, "ls", "rust").Output()

			Expect(strings.Contains(string(command), "1.9.0")).To(Equal(true))
		})

		It("should list remote rust versions", func() {
			Expect(checkRemoteList("rust", "1.x", 120)).To(Equal(true))
		})

		It("should remove rust version", func() {
			result := true

			Execute("go", "run", path, "rust@1.9.0")
			Command("go", "run", path, "rm", "rust@1.9.0").Output()

			plugin := plugins.New("rust")
			versions := plugin.List()

			for _, version := range versions {
				if version == "1.9.0" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})

	Describe("node", func() {
		BeforeEach(func() {
			fmt.Println()
			fmt.Println("Removing node@6.4.0")
			Execute("go", "run", path, "rm", "node@6.4.0")
		})

		It("should install node 6.4.0", func() {
			Execute("go", "run", path, "node@6.4.0")
			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 6.4.0")).To(Equal(true))
		})

		It("should list installed node versions", func() {
			Execute("go", "run", path, "node@6.4.0")
			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 6.4.0")).To(Equal(true))
			Expect(strings.Contains(string(command), "node-v6.4.0-darwin-x64")).To(Equal(false))
		})

		It("should list remote node versions", func() {
			Expect(checkRemoteList("node", "6.x", 5)).To(Equal(true))
		})

		It("should remove node version", func() {
			result := true

			Execute("go", "run", path, "node@6.4.0")
			Command("go", "run", path, "rm", "node@6.4.0").Output()

			plugin := plugins.New("node")
			versions := plugin.List()

			for _, version := range versions {
				if version == "6.4.0" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})

	Describe("ruby", func() {
		BeforeEach(func() {
			fmt.Println()
			fmt.Println("Removing ruby@2.2.1")
			Execute("go", "run", path, "rm", "ruby@2.2.1")
		})

		It("should install ruby 2.2.1", func() {
			Execute("go", "run", path, "ruby@2.2.1")
			command, _ := Command("go", "run", path, "ls", "ruby").Output()

			Expect(strings.Contains(string(command), "♥ 2.2.1")).To(Equal(true))
		})

		It("should list installed ruby versions", func() {
			Execute("go", "run", path, "ruby@2.2.1")
			command, _ := Command("go", "run", path, "ls", "ruby").Output()

			Expect(strings.Contains(string(command), "♥ 2.2.1")).To(Equal(true))
			Expect(strings.Contains(string(command), "ruby-v2.2.1")).To(Equal(false))
		})

		It("should list remote ruby versions", func() {
			if runtime.GOOS == "darwin" {
				Expect(checkRemoteList("ruby", "2.1.x", 5)).To(Equal(true))
			}

			if runtime.GOOS == "linux" {
				Expect(checkRemoteList("ruby", "2.x", 5)).To(Equal(true))
			}
		})

		It("should remove ruby version", func() {
			result := true

			Execute("go", "run", path, "ruby@2.2.1")
			Command("go", "run", path, "rm", "ruby@2.2.1").Output()

			plugin := plugins.New("ruby")
			versions := plugin.List()

			for _, version := range versions {
				if version == "2.2.1" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})
})
