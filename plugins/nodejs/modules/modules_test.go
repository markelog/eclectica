package modules_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/plugins/nodejs/modules"
)

var _ = Describe("modules", func() {
	Describe("SameMajors", func() {
		It("should return true for same majors", func() {
			result := New("6.12.0", "6.11.2").SameMajors()

			Expect(result).To(Equal(true))
		})

		It("should return false for different majors", func() {
			result := New("5.12.0", "6.11.2").SameMajors()

			Expect(result).To(Equal(false))
		})
	})

	Describe("Path", func() {
		It("should return correct path for node_modules", func() {
			result := New("6.12.0", "6.11.2").Path("6.12.0")

			Expect(result).To(ContainSubstring(".eclectica/versions/node/6.12.0/lib/node_modules"))
		})
	})
})
