package main_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/io"
)

var _ = Describe("elm", func() {
	if shouldRun("elm") == false {
		return
	}

	It("should install 0.18.0 version", func() {
		Execute("go", "run", path, "elm@0.18.0")

		command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 0.18.0")).To(Equal(true))
		Expect(err).To(BeNil())

		Execute("go", "run", path, "rm", "elm@0.18.0")
	})

	It("should install 0.17.1 version", func() {
		Execute("go", "run", path, "elm@0.17.1")

		command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 0.17.1")).To(Equal(true))
		Expect(err).To(BeNil())

		Execute("go", "run", path, "rm", "elm@0.17.1")
	})

	It("should install 0.17.0 version", func() {
		Execute("go", "run", path, "elm@0.17.0")

		command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 0.17.0")).To(Equal(true))
		Expect(err).To(BeNil())

		Execute("go", "run", path, "rm", "elm@0.17.0")
	})

	It("should install 0.16.0 version", func() {
		Execute("go", "run", path, "elm@0.16.0")

		command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 0.16.0")).To(Equal(true))
		Expect(err).To(BeNil())

		Execute("go", "run", path, "rm", "elm@0.16.0")
	})

	It("should install 0.15.1 version", func() {
		Execute("go", "run", path, "elm@0.15.1")

		command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 0.15.1")).To(Equal(true))
		Expect(err).To(BeNil())

		Execute("go", "run", path, "rm", "elm@0.15.1")
	})

	It("should install one version after another", func() {
		Execute("go", "run", path, "elm@0.17.0")
		Execute("go", "run", path, "elm@0.18.0")

		command, err := Command("go", "run", path, "ls", "elm").CombinedOutput()

		Expect(strings.Contains(string(command), "♥ 0.18.0")).To(Equal(true))
		Expect(err).To(BeNil())

		Execute("go", "run", path, "rm", "elm@0.17.0")
		Execute("go", "run", path, "rm", "elm@0.18.0")
	})

	It("should use local version", func() {
		pwd, _ := os.Getwd()
		versionFile := filepath.Join(filepath.Dir(pwd), ".elm-version")

		Execute("go", "run", path, "elm@0.17.0")
		Execute("go", "run", path, "elm@0.18.0")

		io.WriteFile(versionFile, "0.17.0")

		command, _ := Command("go", "run", path, "ls", "elm").Output()

		Expect(strings.Contains(string(command), "♥ 0.17.0")).To(Equal(true))

		err := os.RemoveAll(versionFile)

		Expect(err).To(BeNil())

		Execute("go", "run", path, "rm", "elm@0.17.0")
		Execute("go", "run", path, "rm", "elm@0.18.0")
	})
})
