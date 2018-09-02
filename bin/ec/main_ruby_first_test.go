package main_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/plugins"
)

var _ = Describe("ruby", func() {
	if shouldRun("ruby-first") == false {
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

	It("should install two versions of ruby", func() {
		Execute("go", "run", path, "ruby@2.4.2")
		Execute("go", "run", path, "ruby@2.4.1")

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), "2.4.1")).To(Equal(true))
		Expect(strings.Contains(string(ec), "♥ 2.4.1")).To(Equal(true))
	})

	It("should remove ruby version", func() {
		result := true

		Execute("go", "run", path, "ruby@2.2.1")
		Execute("go", "run", path, "ruby@2.1.0")

		Command("go", "run", path, "rm", "ruby@2.2.1").Output()

		plugin := plugins.New(&plugins.Args{
			Language: "ruby",
		})
		versions := plugin.List()

		for _, version := range versions {
			if version == "2.2.1" {
				result = false
			}
		}

		Expect(result).To(Equal(true))
	})

	It("should install two versions of ruby", func() {
		Execute("go", "run", path, "ruby@2.4.2")
		Execute("go", "run", path, "ruby@2.4.1")

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), "2.4.1")).To(Equal(true))
		Expect(strings.Contains(string(ec), "♥ 2.4.1")).To(Equal(true))
	})
})
