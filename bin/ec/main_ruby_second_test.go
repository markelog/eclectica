package main_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/io"
)

var _ = Describe("ruby", func() {
	if shouldRun("ruby-second") == false {
		return
	}

	var (
		mainVersion      = "2.1.5"
		secondaryVersion = "2.2.1"
	)

	BeforeEach(func() {
		fmt.Println()

		fmt.Println("Install " + mainVersion + " version")
		Execute("go", "run", path, "ruby@"+mainVersion)

		fmt.Println("Remove ruby@" + secondaryVersion)
		Execute("go", "run", path, "rm", "ruby@"+secondaryVersion)
	})

	AfterSuite(func() {
		Execute("go", "run", path, "rm", "ruby@"+mainVersion)
		Execute("go", "run", path, "rm", "ruby@"+secondaryVersion)
	})

	It("should use local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".ruby-version")

		Execute("go", "run", path, "ruby@"+mainVersion)
		Execute("go", "run", path, "ruby@"+secondaryVersion)

		io.WriteFile(versionFile, mainVersion)

		command, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(command), "♥ "+mainVersion)).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())
	})

	It("should install ruby "+secondaryVersion, func() {
		Execute("go", "run", path, "ruby@"+secondaryVersion)

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), secondaryVersion)).To(Equal(true))
		Expect(strings.Contains(
			string(ec), "♥ "+secondaryVersion),
		).To(Equal(true))
	})

	It("should install ruby "+mainVersion, func() {
		Execute("go", "run", path, "ruby@"+mainVersion)

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), mainVersion)).To(Equal(true))
		Expect(strings.Contains(string(ec), "♥ "+mainVersion)).To(Equal(true))
	})

	It("should install ruby 2.4.1", func() {
		Execute("go", "run", path, "ruby@2.4.1")

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), "2.4.1")).To(Equal(true))
		Expect(strings.Contains(string(ec), "♥ 2.4.1")).To(Equal(true))

		Execute("go", "run", path, "rm", "2.4.1")
	})

	It("should install bundler", func() {
		Execute("gem", "install", "bundler")

		gems, _ := Command("gem", "list").Output()

		Expect(strings.Contains(string(gems), "bundler")).To(Equal(true))
	})

	It("should list installed ruby versions", func() {
		Execute("go", "run", path, "ruby@"+secondaryVersion)
		command, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(command), "♥ 2.2.1")).To(Equal(true))
		Expect(strings.Contains(string(command), "ruby-v2.2.1")).To(Equal(false))
	})

	It("should list remote ruby versions", func() {
		Expect(checkRemoteList("ruby", "2.x", 5)).To(Equal(true))
	})
})
