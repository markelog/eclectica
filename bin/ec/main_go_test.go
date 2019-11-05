package main_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins"
)

var _ = Describe("go", func() {
	if shouldRun("go") == false {
		return
	}

	var (
		mainVersion      = "1.13.0"
		secondaryVersion = "1.12.0"
	)

	BeforeEach(func() {
		fmt.Println()

		fmt.Println("Install " + secondaryVersion + " version")
		Execute("go", "run", path, "go@"+secondaryVersion)

		fmt.Println("Remove go@" + mainVersion)
		Execute("go", "run", path, "rm", "go@"+mainVersion)
	})

	AfterSuite(func() {
		Execute("go", "run", path, "rm", "go@"+mainVersion)
		Execute("go", "run", path, "rm", "go@"+secondaryVersion)
	})

	It("should list installed versions", func() {
		Execute("go", "run", path, "go@"+mainVersion)
		command, _ := Command("go", "run", path, "ls", "go").Output()

		Expect(strings.Contains(string(command), "♥ "+mainVersion)).To(Equal(true))
	})

	It("should use local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".go-version")

		Execute("go", "run", path, "go@"+mainVersion)

		io.WriteFile(versionFile, secondaryVersion)

		command, _ := Command("go", "run", path, "ls", "go").Output()

		Expect(strings.Contains(string(command), "♥ "+secondaryVersion)).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())
	})

	It("should list remote versions", func() {
		Expect(checkRemoteList("go", "1.9.x", 20)).To(Equal(true))
	})

	It("should remove go version", func() {
		result := true

		Execute("go", "run", path, "go@"+mainVersion)
		Execute("go", "run", path, "go@"+secondaryVersion)
		Command("go", "run", path, "rm", "go@"+mainVersion).Output()

		plugin := plugins.New(&plugins.Args{
			Language: "go",
		})
		versions := plugin.List()

		for _, version := range versions {
			if version == mainVersion {
				result = false
			}
		}

		Expect(result).To(Equal(true))
	})
})
