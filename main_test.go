package main_test

import (
	"path/filepath"
	"time"
	"os/exec"
	"reflect"
	"strings"
	"syscall"
	"bytes"
	"os"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/plugins/rust"
	"github.com/markelog/eclectica/plugins/nodejs"
)

var (
	path string
)

func init() {
	path, _ = filepath.Abs("./main.go")
}

func Command(args... interface{}) *exec.Cmd {
	fn := reflect.ValueOf(exec.Command)
	rargs := make([]reflect.Value, len(args))
    for i, a := range args {
        rargs[i] = reflect.ValueOf(a)
    }

	cmd := fn.Call(rargs)[0].Interface().(*exec.Cmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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

	rust := &rust.Rust{}
	node := &nodejs.Node{}

	Describe("Rust", func() {
		It("should install rust 1.0.0", func() {
			Command("go", "run", path, "rust@1.0.0").Output()

			Expect(rust.Current()).To(Equal("1.0.0"))
		})

		It("should list installed rust versions", func() {
			Command("go", "run", path, "rust@1.0.0").Output()
			command, _ := Command("go", "run", path, "ls", "rust").Output()
			Expect(strings.Contains(string(command), "1.0.0")).To(Equal(true))
		})

		It("should list installed node versions", func() {
			Expect(checkRemoteList("rust", "1.x", 20)).To(Equal(true))
		})

		It("should remove rust version", func() {
			result := true

			Command("go", "run", path, "rust@1.0.0").Output()
			cmd, _ := Command("go", "run", path, "rm", "rust@1.0.0").Output()

			plugin := plugins.New("rust")
			versions := plugin.List()

			for _, version := range versions {
				if version == "1.0.0" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})

	Describe("node", func() {
		It("should install node 6.4.0", func() {
			Command("go", "run", path, "node@6.4.0").Output()

			Expect(node.Current()).To(Equal("6.4.0"))
		})

		It("should list installed node versions", func() {
			Command("go", "run", path, "node@6.4.0")
			command, _ := Command("go", "run", path, "ls", "node").Output()
			Expect(strings.Contains(string(command), "6.4.0")).To(Equal(true))
		})

		It("should list installed node versions", func() {
			Expect(checkRemoteList("node", "6.x", 5)).To(Equal(true))
		})

		It("should remove node version", func() {
			result := true

			Command("go", "run", path, "node@6.4.0").Output()
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
})
