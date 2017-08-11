package main_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins"
)

var _ = Describe("python", func() {
	var (
		pipBin = filepath.Join(bins, "pip")
		eIBin  = filepath.Join(bins, "easy_install")
	)

	Describe("2.x", func() {
		Describe("old", func() {
			if shouldRun("python2-old") == false {
				return
			}

			It(`should install "old" 2.6.9 version`, func() {
				Execute("go", "run", path, "python@2.6.9")

				command, _ := Command("go", "run", path, "ls", "python").Output()

				Expect(strings.Contains(string(command), "♥ 2.6.9")).To(Equal(true))
			})

			It(`should install "old" 2.7.0 version`, func() {
				Execute("go", "run", path, "python@2.7.0")

				command, _ := Command("go", "run", path, "ls", "python").Output()

				Expect(strings.Contains(string(command), "♥ 2.7.0")).To(Equal(true))
			})
		})

		Describe("2.7.x versions", func() {
			if shouldRun("python2.7") == false {
				return
			}

			BeforeEach(func() {
				Execute("go", "run", path, "python@2.7.10")
				Execute("go", "run", path, "python@2.7.12")
			})

			It("should list installed versions", func() {
				command, _ := Command("go", "run", path, "ls", "python").Output()

				Expect(strings.Contains(string(command), "♥ 2.7.12")).To(Equal(true))
			})

			It("should use local version", func() {
				pwd, _ := os.Getwd()
				versionFile := filepath.Join(filepath.Dir(pwd), ".python-version")

				io.WriteFile(versionFile, "2.7.10")

				command, _ := Command("go", "run", path, "ls", "python").Output()
				Expect(strings.Contains(string(command), "♥ 2.7.10")).To(Equal(true))

				err := os.RemoveAll(versionFile)
				Expect(err).To(BeNil())
			})

			It("should list remote versions", func() {
				Expect(checkRemoteList("python", "2.x", 50)).To(Equal(true))
			})

			It("should remove version", func() {
				result := true

				Command("go", "run", path, "rm", "python@2.7.12").Output()

				plugin := plugins.New("python")
				versions := plugin.List()

				for _, version := range versions {
					if version == "2.7.12" {
						result = false
					}
				}

				Expect(result).To(Equal(true))

				Execute("go", "run", path, "python@2.7.12")
			})

			It("should have pip installed when it delivered with binaries", func() {
				command, err := Command(pipBin).CombinedOutput()

				Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
				Expect(err).To(BeNil())
			})

			It("should have easy_install installed when it delivered with binaries", func() {
				command, _ := Command(eIBin).CombinedOutput()

				expected := "error: No urls, filenames, or requirements specified (see --help)"
				actual := string(command)

				Expect(actual).ToNot(ContainSubstring("has not been established"))
				Expect(actual).To(ContainSubstring(expected))
			})

			Describe("2.7.8 version", func() {
				BeforeEach(func() {
					Execute("go", "run", path, "python@2.7.8")
				})

				It("should have pip installed when downloaded", func() {
					command, err := Command(pipBin).CombinedOutput()

					Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
					Expect(err).To(BeNil())
				})

				It("should have easy_install installed when downloaded", func() {
					command, _ := Command(eIBin).CombinedOutput()

					expected := "error: No urls, filenames, or requirements specified (see --help)"
					actual := string(command)

					Expect(actual).ToNot(ContainSubstring("has not been established"))
					Expect(actual).To(ContainSubstring(expected))
				})
			})
		})
	})

	Describe("3.x", func() {
		if shouldRun("python3") == false {
			return
		}

		BeforeEach(func() {
			Execute("go", "run", path, "python@3.5.1")
			Execute("go", "run", path, "python@3.5.2")
		})

		It("should list installed versions", func() {
			command, _ := Command("go", "run", path, "ls", "python").Output()

			Expect(strings.Contains(string(command), "♥ 3.5.2")).To(Equal(true))
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".python-version")

			io.WriteFile(versionFile, "3.5.1")

			command, _ := Command("go", "run", path, "ls", "python").Output()

			Expect(strings.Contains(string(command), "♥ 3.5.1")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
		})

		It("should list remote versions", func() {
			Expect(checkRemoteList("python", "3.x", 50)).To(Equal(true))
		})

		It("should remove version", func() {
			result := true

			Command("go", "run", path, "rm", "python@3.5.2").Output()

			plugin := plugins.New("python")
			versions := plugin.List()

			for _, version := range versions {
				if version == "3.5.2" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
			Execute("go", "run", path, "python@3.5.2")
		})

		It("should have pip installed when it delivered with binaries", func() {
			command, err := Command(pipBin).CombinedOutput()

			Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
			Expect(err).To(BeNil())
		})

		It("should have easy_install installed when it delivered with binaries", func() {
			command, _ := Command(eIBin).CombinedOutput()

			expected := "error: No urls, filenames, or requirements specified (see --help)"
			actual := string(command)

			Expect(actual).ToNot(ContainSubstring("has not been established"))
			Expect(actual).To(ContainSubstring(expected))
		})
	})
})
