package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/variables"
)

var _ = Describe("main", func() {
	if os.Getenv("TEST_ALL") == "" && os.Getenv("TEST_LANGUAGE") == "" {
		return
	}

	Describe("main logic", func() {
		if shouldRun("main") == false {
			return
		}

		It("should output version", func() {
			regVersion := "\\d+\\.\\d+\\.\\d+$"

			command, _ := Command("go", "run", path, "version").Output()
			strCommand := strings.TrimSpace(string(command))

			Expect(strCommand).To(MatchRegexp(regVersion))
		})

		It("should show list without language", func() {
			output := checkRemoteUse()

			Expect(strings.Contains(output, "Language")).To(Equal(true))
			Expect(strings.Contains(output, "node")).To(Equal(true))
		})

		It("should show list with language", func() {
			output := checkRemoteUseWithLanguage("node")

			Expect(strings.Contains(output, "Mask")).To(Equal(true))
			Expect(strings.Contains(output, "6.x")).To(Equal(true))
		})

		It("should support install with partial versions (major)", func() {
			Execute("go", "run", path, "rm", "node@5.12.0")
			Execute("go", "run", path, "node@5")

			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 5.12.0")).To(Equal(true))
		})

		It("should support install with partial versions (minor)", func() {
			Execute("go", "run", path, "rm", "node@5.12.0")
			Execute("go", "run", path, "node@5.12")

			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 5.12.0")).To(Equal(true))
		})

		Describe("ec rm", func() {
			It("should remove version", func() {
				Execute("go", "run", path, "node@6.5.0")
				Execute("go", "run", path, "node@6.4.0")

				Execute("go", "run", path, "rm", "node@6.5.0")

				command, _ := Command("go", "run", path, "ls", "node").Output()

				Expect(strings.Contains(string(command), "6.5.0")).To(Equal(false))
			})
		})

		Describe("local", func() {
			It("should install version but don't switch to it globally", func() {
				current := Command("go", "run", path, "ls", "node")

				pwd, _ := os.Getwd()
				upper := Command("go", "run", path, "ls", "node")
				upper.Dir = filepath.Join(pwd, "..")

				versionFile := filepath.Join(pwd, ".node-version")

				Execute("go", "run", path, "node@6.5.0")
				Execute("go", "run", path, "node@6.4.0", "-l")

				upperRes, _ := upper.CombinedOutput()
				currentRes, _ := current.CombinedOutput()

				Expect(strings.Contains(string(currentRes), "6.5.0")).To(Equal(true))
				Expect(strings.Contains(string(currentRes), "♥ 6.4.0")).To(Equal(true))

				if _, err := os.Stat(versionFile); err != nil {
					Expect(true).To(Equal(false))
				} else {
					os.RemoveAll(versionFile)
				}

				Expect(strings.Contains(string(upperRes), "♥ 6.5.0")).To(Equal(true))
				Expect(strings.Contains(string(upperRes), "6.4.0")).To(Equal(true))
			})

			It("should provide more or less error message if local install is not there", func() {
				pwd, _ := os.Getwd()
				versionFile := filepath.Join(pwd, ".node-version")

				Command("go", "run", path, "rm", "node@6.3.0")

				io.WriteFile(versionFile, "6.3.0")

				result, _ := Command("node", "-v").CombinedOutput()

				expected := string(result)
				actual := "Version \"6.3.0\" was defined on \"./ec/.node-version\" path but this version is not installed"

				Expect(expected).To(ContainSubstring(actual))

				os.RemoveAll(versionFile)
			})
		})
	})

	Describe("rust", func() {
		if shouldRun("rust") == false {
			return
		}

		BeforeEach(func() {
			fmt.Println()

			fmt.Println("Install tmp version")
			Execute("go", "run", path, "rust@1.8.0")

			fmt.Println("Removing rust@1.9.0")
			Execute("go", "run", path, "rm", "rust@1.9.0")
			fmt.Println("Removed")
		})

		It("should install rust 1.9.0", func() {
			Execute("go", "run", path, "rust@1.9.0")
			command, _ := Command("go", "run", path, "ls", "rust").Output()

			fmt.Println()

			Expect(strings.Contains(string(command), "♥ 1.9.0")).To(Equal(true))
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".rust-version")

			Execute("go", "run", path, "rust@1.9.0")

			io.WriteFile(versionFile, "1.8.0")

			command, _ := Command("go", "run", path, "ls", "rust").Output()

			Expect(strings.Contains(string(command), "♥ 1.8.0")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
		})

		It("should list installed rust versions", func() {
			Execute("go", "run", path, "rust@1.9.0")
			command, _ := Command("go", "run", path, "ls", "rust").Output()

			Expect(strings.Contains(string(command), "1.9.0")).To(Equal(true))
		})

		It("should list remote rust versions", func() {
			Expect(checkRemoteList("rust", "1.x", 120)).To(Equal(true))
		})

		It("should remove rust version", func() {
			result := true

			Execute("go", "run", path, "rust@1.9.0")
			Execute("go", "run", path, "rust@1.8.0")
			Command("go", "run", path, "rm", "rust@1.9.0").Output()

			plugin := plugins.New("rust")
			versions, _ := plugin.List()

			for _, version := range versions {
				if version == "1.9.0" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})

	Describe("node", func() {
		mainVersion := "5.1.0"

		if shouldRun("node") == false {
			return
		}

		BeforeEach(func() {
			fmt.Println()

			fmt.Println("Removing node@6.4.0")

			fmt.Println("Install" + mainVersion + " version")
			Execute("go", "run", path, "node@"+mainVersion)

			Execute("go", "run", path, "rm", "node@6.4.0")
			fmt.Println("Removed")
		})

		Describe("preserve globally installed modules", func() {
			It("between major versions", func() {
				Execute("npm", "install", "--global", "ojm")

				Execute("go", "run", path, "node@6.0.0")

				command, _ := Command("ojm").Output()

				expected := "Check if site is down through isup.com"

				Expect(string(command)).Should(ContainSubstring(expected))
			})

			It("between minor versions", func() {
				Execute("npm", "install", "--global", "ojm")

				Execute("go", "run", path, "node@5.0.0")

				command, _ := Command("ojm").Output()

				expected := "Check if site is down through isup.com"

				Expect(string(command)).Should(ContainSubstring(expected))
			})

			It("doesn't try to install global module if its already exist", func() {
				Execute("npm", "install", "--global", "ojm")

				Execute("go", "run", path, "node@5.0.0")
				Execute("go", "run", path, "node@"+mainVersion)

				command, _ := Command("ojm").Output()

				expected := "Check if site is down through isup.com"

				Expect(string(command)).Should(ContainSubstring(expected))
			})
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".node-version")

			Execute("go", "run", path, "node@6.4.0")

			io.WriteFile(versionFile, mainVersion)

			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 5.1.0")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
		})

		It("should install node 6.4.0", func() {
			Execute("go", "run", path, "node@6.4.0")
			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 6.4.0")).To(Equal(true))
		})

		It("test presence of the npmrc config", func() {
			npmrcPath := filepath.Join(variables.Path("node", mainVersion), "/etc/npmrc")

			data := io.Read(npmrcPath)

			Expect(data).To(Equal("scripts-prepend-node-path=false"))
		})

		It("should install node 6.4.0", func() {
			Execute("go", "run", path, "node@6.4.0")
			command, _ := Command("go", "run", path, "ls", "node").Output()

			Expect(strings.Contains(string(command), "♥ 6.4.0")).To(Equal(true))
		})

		It("should list installed node versions", func() {
			Execute("go", "run", path, "node@6.4.0")
			command, _ := Command("go", "run", path, "ls", "node").Output()

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
			Command("go", "run", path, "rm", "node@6.4.0").Output()

			plugin := plugins.New("node")
			versions, _ := plugin.List()

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
	})

	Describe("elm", func() {
		if shouldRun("elm") == false {
			return
		}

		BeforeEach(func() {
		})

		It("should install 0.18.0 version", func() {
			Execute("go", "run", path, "elm@0.18.0")

			command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

			Expect(strings.Contains(string(command), "♥ 0.18.0")).To(Equal(true))
			Expect(err).To(BeNil())
		})

		It("should install 0.17.1 version", func() {
			Execute("go", "run", path, "elm@0.17.1")

			command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

			Expect(strings.Contains(string(command), "♥ 0.17.1")).To(Equal(true))
			Expect(err).To(BeNil())
		})

		It("should install 0.17.0 version", func() {
			Execute("go", "run", path, "elm@0.17.0")

			command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

			Expect(strings.Contains(string(command), "♥ 0.17.0")).To(Equal(true))
			Expect(err).To(BeNil())
		})

		It("should install 0.16.0 version", func() {
			Execute("go", "run", path, "elm@0.16.0")

			command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

			Expect(strings.Contains(string(command), "♥ 0.16.0")).To(Equal(true))
			Expect(err).To(BeNil())
		})

		It("should install 0.15.1 version", func() {
			Execute("go", "run", path, "elm@0.15.1")

			command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

			Expect(strings.Contains(string(command), "♥ 0.15.1")).To(Equal(true))
			Expect(err).To(BeNil())
		})

		It("should install one version after another", func() {
			Execute("go", "run", path, "elm@0.17.0")
			Execute("go", "run", path, "elm@0.18.0")

			command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

			Expect(strings.Contains(string(command), "♥ 0.18.0")).To(Equal(true))
			Expect(err).To(BeNil())
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".elm-version")

			Execute("go", "run", path, "elm@0.18.0")

			io.WriteFile(versionFile, "0.17.0")

			command, _ := Command("go", "run", path, "ls", "elm").Output()

			Expect(strings.Contains(string(command), "♥ 0.17.0")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
		})
	})

	Describe("ruby", func() {
		if shouldRun("ruby") == false {
			return
		}

		BeforeEach(func() {
			fmt.Println()

			fmt.Println("Install tmp version")
			Execute("go", "run", path, "ruby@2.1.5")

			fmt.Println("Removing ruby@2.2.1")
			Execute("go", "run", path, "rm", "ruby@2.2.1")
			fmt.Println("Removed")
		})

		It("should use local version", func() {
			pwd, _ := os.Getwd()
			versionFile := filepath.Join(filepath.Dir(pwd), ".ruby-version")

			Execute("go", "run", path, "ruby@2.2.1")

			io.WriteFile(versionFile, "2.1.5")

			command, _ := Command("go", "run", path, "ls", "ruby").Output()

			Expect(strings.Contains(string(command), "♥ 2.1.5")).To(Equal(true))

			err := os.RemoveAll(versionFile)

			Expect(err).To(BeNil())
		})

		It("should install ruby 2.2.1", func() {
			Execute("go", "run", path, "ruby@2.2.1")
			command, _ := Command("go", "run", path, "ls", "ruby").Output()

			Expect(strings.Contains(string(command), "♥ 2.2.1")).To(Equal(true))
		})

		It("should install bundler", func() {
			tempDir := os.TempDir()
			gems := filepath.Join(tempDir, "gems")

			Execute("go", "run", path, "ruby@2.2.1")
			os.Setenv("GEM_HOME", tempDir)
			Command("gem", "install", "bundler").Output()

			folders, _ := ioutil.ReadDir(gems)
			Expect(strings.Contains(folders[0].Name(), "bundler-")).To(Equal(true))

			os.RemoveAll(gems)
		})

		It("should list installed ruby versions", func() {
			Execute("go", "run", path, "ruby@2.2.1")
			command, _ := Command("go", "run", path, "ls", "ruby").Output()

			Expect(strings.Contains(string(command), "♥ 2.2.1")).To(Equal(true))
			Expect(strings.Contains(string(command), "ruby-v2.2.1")).To(Equal(false))
		})

		It("should list remote ruby versions", func() {
			if runtime.GOOS == "darwin" {
				Expect(checkRemoteList("ruby", "2.1.x", 5)).To(Equal(true))
			}

			if runtime.GOOS == "linux" {
				Expect(checkRemoteList("ruby", "2.x", 5)).To(Equal(true))
			}
		})

		It("should remove ruby version", func() {
			result := true

			Execute("go", "run", path, "ruby@2.2.1")
			Execute("go", "run", path, "ruby@2.1.0")

			Command("go", "run", path, "rm", "ruby@2.2.1").Output()

			plugin := plugins.New("ruby")
			versions, _ := plugin.List()

			for _, version := range versions {
				if version == "2.2.1" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})

	Describe("go", func() {
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
			versions, _ := plugin.List()

			for _, version := range versions {
				if version == "1.7.0" {
					result = false
				}
			}

			Expect(result).To(Equal(true))
		})
	})

	Describe("python", func() {
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
					versions, _ := plugin.List()

					for _, version := range versions {
						if version == "2.7.12" {
							result = false
						}
					}

					Expect(result).To(Equal(true))

					Execute("go", "run", path, "python@2.7.12")
				})

				It("should have pip installed for case when it delivered with binaries", func() {
					command, err := Command(pipBin).CombinedOutput()

					Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
					Expect(err).To(BeNil())
				})

				It("should have easy_install installed for case when it delivered with binaries", func() {
					command, err := Command(eIBin).CombinedOutput()

					Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
					Expect(err).To(BeNil())
				})

				Describe("2.7.8 version", func() {
					BeforeEach(func() {
						Execute("go", "run", path, "python@2.7.8")
					})

					It("should have pip installed for case when downloaded", func() {
						command, err := Command(pipBin).CombinedOutput()

						Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
						Expect(err).To(BeNil())
					})

					It("should have easy_install installed for case when downloaded", func() {
						command, err := Command(eIBin).CombinedOutput()

						Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
						Expect(err).To(BeNil())
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
				versions, _ := plugin.List()

				for _, version := range versions {
					if version == "3.5.2" {
						result = false
					}
				}

				Expect(result).To(Equal(true))
				Execute("go", "run", path, "python@3.5.2")
			})

			It("should have pip installed for case when it delivered with binaries", func() {
				command, err := Command(pipBin).CombinedOutput()

				Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
				Expect(err).To(BeNil())
			})

			It("should have easy_install installed for case when it delivered with binaries", func() {
				command, err := Command(eIBin).CombinedOutput()

				Expect(strings.Contains(string(command), "has not been established")).To(Equal(false))
				Expect(err).To(BeNil())
			})
		})
	})
})
