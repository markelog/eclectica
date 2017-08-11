package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins"
)

var _ = Describe("ruby", func() {
	if shouldRun("ruby") == false {
		return
	}

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

	FIt("should install ruby 2.2.1", func() {
		Execute("go", "run", path, "ruby@2.2.1")

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), "2.2.1")).To(Equal(true))
		Expect(strings.Contains(string(ec), "♥ 2.2.1")).To(Equal(true))
	})

	It("should install ruby 2.1.5", func() {
		Execute("go", "run", path, "ruby@2.1.5")

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), "2.1.5")).To(Equal(true))
		Expect(strings.Contains(string(ec), "♥ 2.1.5")).To(Equal(true))
	})

	It("should install ruby 2.4.1", func() {
		Execute("go", "run", path, "ruby@2.4.1")

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), "2.4.1")).To(Equal(true))
		Expect(strings.Contains(string(ec), "♥ 2.4.1")).To(Equal(true))
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
		versions := plugin.List()

		for _, version := range versions {
			if version == "2.2.1" {
				result = false
			}
		}

		Expect(result).To(Equal(true))
	})
})
