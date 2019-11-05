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

var _ = Describe("rust", func() {
	if shouldRun("rust") == false {
		return
	}

	var (
		mainVersion      = "1.38.0"
		secondaryVersion = "1.37.0"
	)

	BeforeEach(func() {
		fmt.Println()

		fmt.Println("Install " + mainVersion + " version")
		Execute("go", "run", path, "rust@"+mainVersion)

		fmt.Println("Remove rust@" + secondaryVersion)
		Execute("go", "run", path, "rm", "rust@"+secondaryVersion)
	})

	AfterSuite(func() {
		Execute("go", "run", path, "rm", "rust@"+mainVersion)
		Execute("go", "run", path, "rm", "rust@"+secondaryVersion)
	})

	It("should install rust "+secondaryVersion, func() {
		Execute("go", "run", path, "rust@"+secondaryVersion)
		command, _ := Command("go", "run", path, "ls", "rust").Output()
		result := string(command)

		fmt.Println()

		Expect(strings.Contains(result, "♥ "+secondaryVersion)).To(Equal(true))
	})

	It("should use local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".rust-version")

		Execute("go", "run", path, "rust@"+secondaryVersion)

		io.WriteFile(versionFile, ""+mainVersion)

		command, _ := Command("go", "run", path, "ls", "rust").Output()

		Expect(strings.Contains(string(command), "♥ "+mainVersion)).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())
	})

	It("should list installed rust versions", func() {
		Execute("go", "run", path, "rust@"+secondaryVersion)
		command, _ := Command("go", "run", path, "ls", "rust").Output()

		Expect(strings.Contains(string(command), ""+secondaryVersion)).To(Equal(true))
	})

	It("should list remote rust versions", func() {
		Expect(checkRemoteList("rust", "1.x", 120)).To(Equal(true))
	})

	It("should remove rust version", func() {
		result := true

		Execute("go", "run", path, "rust@"+secondaryVersion)
		Execute("go", "run", path, "rust@"+mainVersion)
		Command("go", "run", path, "rm", "rust@"+secondaryVersion).Output()

		plugin := plugins.New(&plugins.Args{
			Language: "rust",
		})
		versions := plugin.List()

		for _, version := range versions {
			if version == ""+secondaryVersion {
				result = false
			}
		}

		Expect(result).To(Equal(true))
	})
})
