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

	var (
		mainVersion      = "2.4.2"
		secondaryVersion = "2.2.1"
	)

	BeforeEach(func() {
		fmt.Println()

		fmt.Println("Install " + mainVersion + " version")
		Execute("go", "run", path, "ruby@"+mainVersion)

		fmt.Println("Remove ruby@" + secondaryVersion)
		Execute("go", "run", path, "rm", "ruby@"+secondaryVersion)
	})

	teardown := func() {
		Execute("go", "run", path, "rm", "ruby@"+mainVersion)
		Execute("go", "run", path, "rm", "ruby@"+secondaryVersion)
	}

	It("should install two versions of ruby", func() {
		Execute("go", "run", path, "ruby@"+mainVersion)
		Execute("go", "run", path, "ruby@"+secondaryVersion)

		ruby, _ := Command("ruby", "--version").Output()
		ec, _ := Command("go", "run", path, "ls", "ruby").Output()

		Expect(strings.Contains(string(ruby), secondaryVersion)).To(Equal(true))
		Expect(strings.Contains(
			string(ec), "♥ "+secondaryVersion),
		).To(Equal(true))
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

		Command("go", "run", path, "rm", "ruby@2.1.0").Output()
	})

	teardown()
})
