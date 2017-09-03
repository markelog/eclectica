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
	"github.com/markelog/eclectica/variables"
)

var _ = Describe("node", func() {
	mainVersion := "5.1.0"

	if shouldRun("node") == false {
		return
	}

	BeforeEach(func() {
		fmt.Println()

		fmt.Println("Install " + mainVersion + " version")
		Execute("go", "run", path, "node@"+mainVersion)

		fmt.Println("Removing node@6.4.0")
		Execute("go", "run", path, "rm", "node@6.4.0")
		fmt.Println("Removed")
	})

	Describe("preserve globally installed modules", func() {
		It("between major versions (ojm module)", func() {
			Execute("npm", "install", "--global", "ojm")

			Execute("go", "run", path, "node@6.0.0")

			command, _ := Command("ojm").CombinedOutput()

			expected := "Check if site is down through isup.com"

			Expect(string(command)).Should(ContainSubstring(expected))
		})

		It("between minor versions (ojm module)", func() {
			Execute("npm", "install", "--global", "ojm")

			Execute("go", "run", path, "node@5.0.0")

			command, _ := Command("ojm").CombinedOutput()

			expected := "Check if site is down through isup.com"

			Expect(string(command)).Should(ContainSubstring(expected))
		})
	})

	It("should use local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".node-version")

		Execute("go", "run", path, "node@6.4.0")

		io.WriteFile(versionFile, mainVersion)

		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 5.1.0")).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())
	})

	It("should install node 6.4.0", func() {
		Execute("go", "run", path, "node@6.4.0")
		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 6.4.0")).To(Equal(true))
	})

	It("test presence of the npmrc config", func() {
		npmrcPath := filepath.Join(variables.Path("node", mainVersion), "/etc/npmrc")

		data := io.Read(npmrcPath)

		Expect(data).To(Equal("scripts-prepend-node-path=false"))
	})

	It("should install node 6.4.0", func() {
		Execute("go", "run", path, "node@6.4.0")
		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 6.4.0")).To(Equal(true))
	})

	It("should list installed node versions", func() {
		Execute("go", "run", path, "node@6.4.0")
		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 6.4.0")).To(Equal(true))
		Expect(strings.Contains(string(command), "node-v6.4.0-darwin-x64")).To(Equal(false))
	})

	It("should list remote node versions", func() {
		Expect(checkRemoteList("node", "6.x", 5)).To(Equal(true))
	})

	It("should remove node version", func() {
		result := true

		Execute("go", "run", path, "node@6.4.0")
		Execute("go", "run", path, "node@"+mainVersion)
		Command("go", "run", path, "rm", "node@6.4.0").CombinedOutput()

		plugin := plugins.New("node")
		versions := plugin.List()

		for _, version := range versions {
			if version == "6.4.0" {
				result = false
			}
		}

		Expect(result).To(Equal(true))
	})

	It("should install yarn", func() {
		yarnBin := filepath.Join(variables.Path("node", mainVersion), "bin/yarn")

		bytes, _ := Command(yarnBin, "--help").CombinedOutput()

		Expect(string(bytes)).To(ContainSubstring("Usage: yarn [command] [flags]"))
	})

	It("should still have yarn after additional installation", func() {
		Execute("go", "run", path, "node@6.4.0")

		yarnBin := filepath.Join(variables.Path("node", "6.4.0"), "bin/yarn")

		bytes, _ := Command(yarnBin, "--help").CombinedOutput()

		Expect(string(bytes)).To(ContainSubstring("Usage: yarn [command] [flags]"))
	})
})
