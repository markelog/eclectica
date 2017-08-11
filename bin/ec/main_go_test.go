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

	BeforeEach(func() {
		fmt.Println()

		fmt.Println("Install tmp version")
		Execute("go", "run", path, "go@1.6.0")

		fmt.Println("Removing go@1.7.0")
		Execute("go", "run", path, "rm", "go@1.7.0")
		fmt.Println("Removed")
	})

	It("should list installed versions", func() {
		Execute("go", "run", path, "go@1.7.0")
		command, _ := Command("go", "run", path, "ls", "go").Output()

		Expect(strings.Contains(string(command), "♥ 1.7.0")).To(Equal(true))
	})

	It("should use local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".go-version")

		Execute("go", "run", path, "go@1.7.0")

		io.WriteFile(versionFile, "1.6.0")

		command, _ := Command("go", "run", path, "ls", "go").Output()

		Expect(strings.Contains(string(command), "♥ 1.6.0")).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())
	})

	It("should list remote versions", func() {
		Expect(checkRemoteList("go", "1.7.x", 10)).To(Equal(true))
	})

	It("should remove go version", func() {
		result := true

		Execute("go", "run", path, "go@1.7.0")
		Execute("go", "run", path, "go@1.6.0")
		Command("go", "run", path, "rm", "go@1.7.0").Output()

		plugin := plugins.New("go")
		versions := plugin.List()

		for _, version := range versions {
			if version == "1.7.0" {
				result = false
			}
		}

		Expect(result).To(Equal(true))
	})
})
