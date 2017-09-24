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
	secondaryVersion := "6.4.0"

	if shouldRun("node") == false {
		return
	}

	BeforeEach(func() {
		fmt.Println()

		Execute("go", "run", path, "rm", "node@"+mainVersion)
		fmt.Println("Install " + mainVersion + " version")
		Execute("go", "run", path, "node@"+mainVersion)

		Execute("go", "run", path, "rm", "node@"+secondaryVersion)
	})

	Describe("preserve globally installed modules", func() {
		dir, _ := os.Getwd()
		testdata := filepath.Join(dir, "../../testdata/plugins/nodejs/example.scss")

		It("between major versions (ojm module)", func() {
			Execute("npm", "install", "--global", "ojm")

			Execute("go", "run", path, "node@"+secondaryVersion)

			command, _ := Command("ojm").CombinedOutput()

			expected := "Check if site is down through isup.com"

			Expect(string(command)).Should(ContainSubstring(expected))
		})

		It("between major versions (node-sass module)", func() {
			Execute("npm", "install", "--global", "node-sass")

			Execute("go", "run", path, "node@"+secondaryVersion)

			command, _ := Command("node-sass", testdata).CombinedOutput()
			expected := "background: #eeffcc;"

			Expect(string(command)).Should(ContainSubstring(expected))
		})

		It("between minor versions (ojm module)", func() {
			Execute("npm", "install", "--global", "ojm")

			Execute("go", "run", path, "node@5.0.0")

			command, _ := Command("ojm").CombinedOutput()
			expected := "Check if site is down through isup.com"

			Expect(string(command)).Should(ContainSubstring(expected))

			Execute("go", "run", path, "rm", "node@5.0.0")
		})

		It("between minor versions (node-sass module)", func() {
			Execute("npm", "install", "--global", "node-sass")

			Execute("go", "run", path, "node@5.0.0")

			command, _ := Command("node-sass", testdata).CombinedOutput()
			expected := "background: #eeffcc;"

			Expect(string(command)).Should(ContainSubstring(expected))

			Execute("go", "run", path, "rm", "node@5.0.0")
		})
	})

	It("should use local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".node-version")

		Execute("go", "run", path, "node@"+secondaryVersion)

		io.WriteFile(versionFile, mainVersion)

		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 5.1.0")).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())
	})

	It("should install node "+secondaryVersion, func() {
		Execute("go", "run", path, "node@"+secondaryVersion)
		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ "+secondaryVersion)).To(Equal(true))
	})

	It("test presence of the npmrc config", func() {
		npmrcPath := filepath.Join(variables.Path("node", mainVersion), "/etc/npmrc")

		data := io.Read(npmrcPath)

		Expect(data).To(Equal("scripts-prepend-node-path=false"))
	})

	It("should list installed node versions", func() {
		Execute("go", "run", path, "node@"+secondaryVersion)
		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ "+secondaryVersion)).To(Equal(true))
		Expect(strings.Contains(string(command), mainVersion)).To(Equal(true))
		Expect(strings.Contains(string(command), "node-v"+secondaryVersion+"-darwin-x64")).To(Equal(false))
	})

	It("should list remote node versions", func() {
		Expect(checkRemoteList("node", "6.x", 5)).To(Equal(true))
	})

	It("should remove node version", func() {
		success := true

		Execute("go", "run", path, "node@"+secondaryVersion)
		Execute("go", "run", path, "node@"+mainVersion)
		Command("go", "run", path, "rm", "node@"+secondaryVersion).CombinedOutput()

		plugin := plugins.New("node")
		versions := plugin.List()

		for _, version := range versions {
			if version == secondaryVersion {
				success = false
			}
		}

		Expect(success).To(Equal(true))
	})

	It("should install yarn", func() {
		yarnBin := filepath.Join(variables.Path("node", mainVersion), "bin/yarn")

		bytes, _ := Command(yarnBin, "--help").CombinedOutput()

		Expect(string(bytes)).To(ContainSubstring("Usage: yarn [command] [flags]"))
	})

	It("should still have yarn after additional installation", func() {
		Execute("go", "run", path, "node@"+secondaryVersion)

		yarnBin := filepath.Join(variables.Path("node", secondaryVersion), "bin/yarn")

		bytes, _ := Command(yarnBin, "--help").CombinedOutput()

		Expect(string(bytes)).To(ContainSubstring("Usage: yarn [command] [flags]"))
	})
})
