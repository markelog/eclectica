package main_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/io"
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

var _ = Describe("main", func() {
	if os.Getenv("INT") != "true" {
		return
	}

	Describe("main logic", func() {
		It("should output version", func() {
			regVersion := "\\d+\\.\\d+\\.\\d+$"

			command, _ := Command("go", "run", path, "version").Output()
			strCommand := strings.TrimSpace(string(command))

			Expect(strCommand).To(MatchRegexp(regVersion))
		})

		It("should show list without language", func() {
			output := checkRemoteUse()

			Expect(strings.Contains(output, "Language")).To(Equal(true))
			Expect(strings.Contains(output, "node")).To(Equal(true))
		})

		It("should show list with language", func() {
			output := checkRemoteUseWithLanguage("node")

			Expect(strings.Contains(output, "Mask")).To(Equal(true))
			Expect(strings.Contains(output, "6.x")).To(Equal(true))
		})

		Describe("ec rm", func() {
			It("should not allow removal of current version", func() {
				Execute("go", "run", path, "node@6.5.0")

				output, _ := Command("go", "run", path, "rm", "node@6.5.0").CombinedOutput()
				out := string(output)

				Expect(strings.Contains(out, "Cannot remove active version")).To(Equal(true))
			})

			It("should remove version", func() {
				Execute("go", "run", path, "node@6.5.0")
				Execute("go", "run", path, "node@6.4.0")

				Execute("go", "run", path, "rm", "node@6.5.0")

				command, _ := Command("go", "run", path, "ls", "node").Output()

				Expect(strings.Contains(string(command), "6.5.0")).To(Equal(false))
			})
		})
	})

	Describe("Rust", func() {
		BeforeEach(func() {
			fmt.Println()

			fmt.Println("Install tmp version")
			Execute("go", "run", path, "rust@1.8.0")

			fmt.Println("Removing rust@1.9.0")
			Execute("go", "run", path, "rm", "rust@1.9.0")
			fmt.Println("Removed")
		})

		It("should install rust 1.9.0", func() {
			Execute("go", "run", path, "rust@1.9.0")
			command, _ := Command("go", "run", path, "ls", "rust").Output()

			fmt.Println()

			Expect(strings.Contains(string(command), "♥ 1.9.0")).To(Equal(true))
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".rust-version")

			Execute("go", "run", path, "rust@1.9.0")

			io.WriteFile(versionFile, "1.8.0")

			command, _ := Command("go", "run", path, "ls", "rust").Output()

			Expect(strings.Contains(string(command), "♥ 1.8.0")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
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
			Execute("go", "run", path, "rust@1.8.0")
			Command("go", "run", path, "rm", "rust@1.9.0").Output()

			plugin := plugins.New("rust")
			versions, _ := plugin.List()

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

			fmt.Println("Install tmp version")
			Execute("go", "run", path, "node@5.1.0")

			Execute("go", "run", path, "rm", "node@6.4.0")
			fmt.Println("Removed")
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".node-version")

			Execute("go", "run", path, "node@6.4.0")

			io.WriteFile(versionFile, "5.1.0")

			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 5.1.0")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
		})

		It("should install node 6.4.0", func() {
			Execute("go", "run", path, "node@6.4.0")
			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 6.4.0")).To(Equal(true))
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
			Execute("go", "run", path, "node@5.1.0")
			Command("go", "run", path, "rm", "node@6.4.0").Output()

			plugin := plugins.New("node")
			versions, _ := plugin.List()

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

			fmt.Println("Install tmp version")
			Execute("go", "run", path, "ruby@2.1.5")

			fmt.Println("Removing ruby@2.2.1")
			Execute("go", "run", path, "rm", "ruby@2.2.1")
			fmt.Println("Removed")
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".ruby-version")

			Execute("go", "run", path, "ruby@2.2.1")

			io.WriteFile(versionFile, "2.1.5")

			command, _ := Command("go", "run", path, "ls", "ruby").Output()

			Expect(strings.Contains(string(command), "♥ 2.1.5")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
		})

		It("should install ruby 2.2.1", func() {
			Execute("go", "run", path, "ruby@2.2.1")
			command, _ := Command("go", "run", path, "ls", "ruby").Output()

			Expect(strings.Contains(string(command), "♥ 2.2.1")).To(Equal(true))
		})

		It("should install bundler", func() {
			tempDir := os.TempDir()
			gems := filepath.Join(tempDir, "gems")

			Execute("go", "run", path, "ruby@2.2.1")
			os.Setenv("GEM_HOME", tempDir)
			Command("gem", "install", "bundler").Output()

			folders, _ := ioutil.ReadDir(gems)
			Expect(strings.Contains(folders[0].Name(), "bundler-")).To(Equal(true))

			os.RemoveAll(gems)
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
			Execute("go", "run", path, "ruby@2.1.0")

			Command("go", "run", path, "rm", "ruby@2.2.1").Output()

			plugin := plugins.New("ruby")
			versions, _ := plugin.List()

			for _, version := range versions {
				if version == "2.2.1" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})

	Describe("go", func() {
		BeforeEach(func() {
			fmt.Println()

			fmt.Println("Install tmp version")
			Execute("go", "run", path, "go@1.6")

			fmt.Println("Removing go@1.7")
			Execute("go", "run", path, "rm", "go@1.7")
			fmt.Println("Removed")
		})

		It("should list installed go versions", func() {
			Execute("go", "run", path, "go@1.7")
			command, _ := Command("go", "run", path, "ls", "go").Output()

			Expect(strings.Contains(string(command), "1.7")).To(Equal(true))
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".go-version")

			Execute("go", "run", path, "go@1.7")

			io.WriteFile(versionFile, "1.6")

			command, _ := Command("go", "run", path, "ls", "go").Output()

			Expect(strings.Contains(string(command), "♥ 1.6")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
		})

		It("should list remote go versions", func() {
			Expect(checkRemoteList("go", "1.7.x", 5)).To(Equal(true))
		})

		It("should remove go version", func() {
			result := true

			Execute("go", "run", path, "go@1.7")
			Execute("go", "run", path, "go@1.6")
			Command("go", "run", path, "rm", "go@1.7").Output()

			plugin := plugins.New("go")
			versions, _ := plugin.List()

			for _, version := range versions {
				if version == "1.7" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})
})
