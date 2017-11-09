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

	BeforeEach(func() {
		fmt.Println()

		fmt.Println("Install tmp version")
		Execute("go", "run", path, "rust@1.20.0")

		fmt.Println("Removing rust@1.21.0")
		Execute("go", "run", path, "rm", "rust@1.21.0")
		fmt.Println("Removed")
	})

	It("should install rust 1.21.0", func() {
		Execute("go", "run", path, "rust@1.21.0")
		command, _ := Command("go", "run", path, "ls", "rust").Output()

		fmt.Println()

		Expect(strings.Contains(string(command), "♥ 1.21.0")).To(Equal(true))
	})

	It("should use local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".rust-version")

		Execute("go", "run", path, "rust@1.21.0")

		io.WriteFile(versionFile, "1.20.0")

		command, _ := Command("go", "run", path, "ls", "rust").Output()

		Expect(strings.Contains(string(command), "♥ 1.20.0")).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())
	})

	It("should list installed rust versions", func() {
		Execute("go", "run", path, "rust@1.21.0")
		command, _ := Command("go", "run", path, "ls", "rust").Output()

		Expect(strings.Contains(string(command), "1.21.0")).To(Equal(true))
	})

	It("should list remote rust versions", func() {
		Expect(checkRemoteList("rust", "1.x", 120)).To(Equal(true))
	})

	It("should remove rust version", func() {
		result := true

		Execute("go", "run", path, "rust@1.21.0")
		Execute("go", "run", path, "rust@1.20.0")
		Command("go", "run", path, "rm", "rust@1.21.0").Output()

		plugin := plugins.New("rust")
		versions := plugin.List()

		for _, version := range versions {
			if version == "1.21.0" {
				result = false
			}
		}

		Expect(result).To(Equal(true))
	})
})
