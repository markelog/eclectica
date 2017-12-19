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
	if shouldRun("node") == false {
		return
	}

	var (
		mainVersion      = "5.12.0"
		secondaryVersion = "6.4.0"
	)

	BeforeEach(func() {
		fmt.Println()

		Execute("go", "run", path, "rm", "node@"+mainVersion)
		fmt.Println("Install " + mainVersion + " version")
		Execute("go", "run", path, "node@"+mainVersion)

		Execute("go", "run", path, "rm", "node@"+secondaryVersion)
	})

	Describe("globally installed modules", func() {
		dir, _ := os.Getwd()

		It("preserves between major versions (ojm module)", func() {
			Execute("npm", "install", "--global", "ojm")

			Execute("go", "run", path, "node@"+secondaryVersion, "-w")

			expected := "Check if site is down through isup.com"
			command, _ := Command("ojm").CombinedOutput()
			output := string(command)

			Expect(output).Should(ContainSubstring(expected))
		})

		It("preserves between major versions (node-sass module)", func() {
			testdata := filepath.Join(dir, "../../testdata/plugins/nodejs/example.scss")
			Execute("npm", "install", "--global", "node-sass")

			Execute("go", "run", path, "node@"+secondaryVersion, "-w")

			expected := "background: #eeffcc;"
			command, _ := Command("node-sass", testdata).CombinedOutput()

			Expect(string(command)).Should(ContainSubstring(expected))
		})

		It("preserves between minor versions (ojm module)", func() {
			Execute("npm", "install", "--global", "ojm")

			Execute("go", "run", path, "node@5.0.0", "-w")

			command, _ := Command("ojm").CombinedOutput()
			expected := "Check if site is down through isup.com"

			Expect(string(command)).Should(ContainSubstring(expected))

			Execute("go", "run", path, "rm", "node@5.0.0")
		})

		It("preserves between minor versions (node-sass module)", func() {
			testdata := filepath.Join(dir, "../../testdata/plugins/nodejs/example.scss")
			Execute("npm", "install", "--global", "node-sass")

			Execute("go", "run", path, "node@5.0.0", "-w")

			command, _ := Command("node-sass", testdata).CombinedOutput()
			expected := "background: #eeffcc;"

			Expect(string(command)).Should(ContainSubstring(expected))

			Execute("go", "run", path, "rm", "node@5.0.0")
		})

		It("does not preserves modules", func() {
			Execute("npm", "install", "--global", "ojm")

			Execute("go", "run", path, "node@"+secondaryVersion)

			command, _ := Command("ojm").CombinedOutput()

			expected := "Check if site is down through isup.com"

			Expect(string(command)).ShouldNot(ContainSubstring(expected))
		})
	})

	Describe("installation of the local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".node-version")

		BeforeEach(func() {
			io.WriteFile(versionFile, mainVersion)
		})

		AfterEach(func() {
			os.RemoveAll(versionFile)
		})

		It("should use local version", func() {
			Execute("go", "run", path, "node@"+secondaryVersion)

			command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()
			output := string(command)

			Expect(output).Should(ContainSubstring("♥ " + mainVersion))
		})
	})

	It("should install node "+secondaryVersion, func() {
		Execute("go", "run", path, "node@"+secondaryVersion)

		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()
		output := string(command)

		println(output)

		Expect(strings.Contains(output, "♥ "+secondaryVersion)).To(Equal(true))
	})

	It("should list installed node versions", func() {
		Execute("go", "run", path, "node@"+secondaryVersion)
		command, _ := Command("go", "run", path, "ls", "node").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ "+secondaryVersion)).To(Equal(true))
		Expect(strings.Contains(string(command), mainVersion)).To(Equal(true))
		Expect(strings.Contains(string(command), "node-v"+secondaryVersion+"-darwin-x64")).To(Equal(false))
	})

	It("should list remote node versions", func() {
		Expect(checkRemoteList("node", "6.x", 15)).To(Equal(true))
	})

	It("should remove node version", func() {
		success := true

		Execute("go", "run", path, "node@"+secondaryVersion)
		Execute("go", "run", path, "node@"+mainVersion)
		Command("go", "run", path, "rm", "node@"+secondaryVersion).CombinedOutput()

		plugin := plugins.New(&plugins.Args{
			Language: "node",
		})
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
