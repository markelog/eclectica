package main_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/io"
)

var _ = Describe("main logic", func() {
	if shouldRun("main") == false {
		return
	}

	It("should output version", func() {
		regVersion := "\\d+\\.\\d+\\.\\d+$"

		command, _ := Command("go", "run", path, "version").Output()
		strCommand := strings.TrimSpace(string(command))

		Expect(strCommand).To(MatchRegexp(regVersion))
	})

	It("should show proper help", func() {
		command, _ := Command("go", "run", path, "--help").Output()
		strCommand := strings.TrimSpace(string(command))

		Expect(strCommand).To(ContainSubstring(`ec [command] [flags] [<language>@<version>]`))
	})

	It("should show error for language typo", func() {
		command, _ := Command("go", "run", path, "noda").CombinedOutput()
		strCommand := strings.TrimSpace(string(command))

		Expect(strCommand).To(
			ContainSubstring(
				`Eclectica does not support`,
			),
		)

		Expect(strCommand).To(
			ContainSubstring(
				`node`,
			),
		)

		Expect(strCommand).To(
			ContainSubstring(
				`noda`,
			),
		)
	})

	It("should show list without language", func() {
		output := checkRemoteUse()

		Expect(strings.Contains(output, "langauge:")).To(Equal(true))
		Expect(strings.Contains(output, "node")).To(Equal(true))
	})

	It("should show list with language", func() {
		output := checkRemoteUseWithLanguage("node")

		Expect(strings.Contains(output, "    mask:")).To(Equal(true))
		Expect(strings.Contains(output, "6.x")).To(Equal(true))
	})

	Describe("partial install", func() {
		BeforeEach(func() {
			Execute("go", "run", path, "rm", "node@5.12.0")
		})

		AfterEach(func() {
			Execute("go", "run", path, "rm", "node@5.12.0")
		})

		It("should support install with partial versions (major)", func() {
			Execute("go", "run", path, "node@5")

			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 5.12.0")).To(Equal(true))
		})

		It("should support install with partial versions (minor)", func() {
			Execute("go", "run", path, "node@5.12")

			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 5.12.0")).To(Equal(true))
		})
	})

	Describe("ec rm", func() {
		BeforeEach(func() {
			Execute("go", "run", path, "node@6.5.0")
		})

		It("should remove version", func() {
			Execute("go", "run", path, "rm", "node@6.5.0")
			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "6.5.0")).To(Equal(false))
		})
	})

	Describe("ec ls", func() {
		BeforeEach(func() {
			Execute("go", "run", path, "rm", "node@5.12.0")
			Execute("go", "run", path, "node@5")
		})

		AfterEach(func() {
			Execute("go", "run", path, "rm", "node@5.12.0")
		})

		It("should not double list the package", func() {
			Execute("go", "run", path, "rm", "node@5.12.0")
			Execute("go", "run", path, "node@5")

			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(string(command)).To(ContainSubstring("♥ 5.12.0"))
			Expect(string(command)).ToNot(ContainSubstring("  5.12.0"))
		})
	})

	Describe("local", func() {
		Describe("vs global for 6.x versions", func() {
			var (
				pwd         string
				versionFile string
			)

			BeforeEach(func() {
				pwd, _ = os.Getwd()
				versionFile = filepath.Join(pwd, ".node-version")
			})

			AfterEach(func() {
				Execute("go", "run", path, "rm", "node@6.5.0")
				Execute("go", "run", path, "rm", "node@6.4.0")

				if _, err := os.Stat(versionFile); err == nil {
					os.RemoveAll(versionFile)
				}
			})

			It("should install version but don't switch to it globally", func() {
				current := Command("go", "run", path, "ls", "node")

				upper := Command("go", "run", path, "ls", "node")
				upper.Dir = filepath.Join(pwd, "..")

				Execute("go", "run", path, "node@6.5.0")
				Execute("go", "run", path, "node@6.4.0", "-l")

				upperRes, _ := upper.CombinedOutput()
				currentRes, _ := current.CombinedOutput()

				_, err := os.Stat(versionFile)
				Expect(err).To(BeNil())

				Expect(strings.Contains(string(currentRes), "6.5.0")).To(Equal(true))
				Expect(strings.Contains(string(currentRes), "♥ 6.4.0")).To(Equal(true))

				Expect(strings.Contains(string(upperRes), "♥ 6.5.0")).To(Equal(true))
				Expect(strings.Contains(string(upperRes), "6.4.0")).To(Equal(true))
			})
		})

		Describe("useful error messages", func() {
			var (
				pwd         string
				versionFile string
			)

			BeforeEach(func() {
				pwd, _ = os.Getwd()
				versionFile = filepath.Join(pwd, ".node-version")

				Execute("go", "run", path, "node@4.0.0")
			})

			AfterEach(func() {
				Execute("go", "run", path, "rm", "node@4.0.0")

				os.RemoveAll(versionFile)
			})

			It("should provide more or less complete error message if local install is not there", func() {
				io.WriteFile(versionFile, "6.3.0")

				result, _ := Command("node", "-v").CombinedOutput()

				actual := string(result)

				expected := `version "6.3.0" was defined on "./ec/.node-version" path but this version is not installed`

				Expect(actual).To(ContainSubstring(expected))
			})

			It("should provide more or less complete error message if mask is not satisfactory", func() {
				pwd, _ := os.Getwd()
				versionFile := filepath.Join(pwd, ".node-version")

				io.WriteFile(versionFile, "5")

				result, _ := Command("node", "-v").CombinedOutput()

				actual := string(result)
				expected := `mask "5" was defined on "./ec/.node-version" path but none of these versions were installed`

				Expect(actual).To(ContainSubstring(expected))
			})

		})

		Describe("mask: of 5 versions", func() {
			BeforeEach(func() {
				Execute("go", "run", path, "node@5.11.0")
				Execute("go", "run", path, "node@5.12.0")
			})

			AfterEach(func() {
				Execute("go", "run", path, "rm", "node@5.11.0")
				Execute("go", "run", path, "rm", "node@5.12.0")
			})

			It("should choose latest available version", func() {
				pwd, _ := os.Getwd()
				versionFile := filepath.Join(pwd, ".node-version")

				io.WriteFile(versionFile, "5")

				result, _ := Command("node", "-v").CombinedOutput()

				actual := string(result)
				expected := `v5.12.0`

				Expect(actual).To(ContainSubstring(expected))

				os.RemoveAll(versionFile)
			})
		})
	})
})
