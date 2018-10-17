package main_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/markelog/eclectica/io"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ruby", func() {
	if shouldRun("ruby-second") == false {
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

		Execute("go", "run", path, "ruby@2.1.5")
		Execute("go", "run", path, "ruby@2.2.1")

		io.WriteFile(versionFile, "2.1.5")

		command, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(command), "♥ 2.1.5")).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())
	})

	It("should install ruby 2.2.1", func() {
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
		Execute("gem", "install", "bundler")

		gems, _ := Command("gem", "list").Output()

		Expect(strings.Contains(string(gems), "bundler")).To(Equal(true))
	})

	It("should list installed ruby versions", func() {
		Execute("go", "run", path, "ruby@2.2.1")
		command, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(command), "♥ 2.2.1")).To(Equal(true))
		Expect(strings.Contains(string(command), "ruby-v2.2.1")).To(Equal(false))
	})

	It("should list remote ruby versions", func() {
		Expect(checkRemoteList("ruby", "2.x", 5)).To(Equal(true))
	})
})
