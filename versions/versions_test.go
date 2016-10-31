package versions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/versions"
)

var _ = Describe("versions", func() {
	Describe("ComposeVersions", func() {
		It("should compose major versions", func() {
			compose := Compose([]string{"0.8.2", "4.4.7", "6.3.0", "6.4.2"})

			Expect(compose["0.x"]).To(Equal([]string{"0.8.2"}))
			Expect(compose["4.x"]).To(Equal([]string{"4.4.7"}))
			Expect(compose["6.x"]).To(Equal([]string{"6.3.0", "6.4.2"}))
		})

		It("should compose minor versions", func() {
			compose := Compose([]string{"2.1.1", "2.2.1", "2.3.1"})

			Expect(compose["2.1.x"]).To(Equal([]string{"2.1.1"}))
			Expect(compose["2.2.x"]).To(Equal([]string{"2.2.1"}))
			Expect(compose["2.3.x"]).To(Equal([]string{"2.3.1"}))
		})

		It("should compose peculiar versions", func() {
			compose := Compose([]string{"1.4.3", "1.5beta1", "1.5beta2", "1.5rc1"})

			Expect(compose["1.4.x"]).To(Equal([]string{"1.4.3"}))
			Expect(compose["1.5.x"]).To(Equal([]string{"1.5beta1", "1.5beta2", "1.5rc1"}))
		})
	})

	Describe("GetKeys", func() {
		It("should get version keys", func() {
			list := map[string][]string{"4.x": []string{}, "0.x": []string{"0.8.2"}}
			keys := GetKeys(list)

			Expect(keys[1]).To(Equal("0.x"))
			Expect(keys[0]).To(Equal("4.x"))
		})
	})

	Describe("GetElements", func() {
		It("should get version elements", func() {
			list := Compose([]string{"0.8.2", "4.4.7", "6.3.0", "6.4.2"})
			elements := GetElements("6.x", list)

			Expect(elements[1]).To(Equal("6.3.0"))
			Expect(elements[0]).To(Equal("6.4.2"))
		})

		It("should get sorted version elements", func() {
			list := Compose([]string{"1.12.3", "1.12.0", "1.12.1", "1.12.2"})
			elements := GetElements("1.x", list)

			Expect(elements[3]).To(Equal("1.12.0"))
			Expect(elements[2]).To(Equal("1.12.1"))
			Expect(elements[1]).To(Equal("1.12.2"))
			Expect(elements[0]).To(Equal("1.12.3"))
		})
	})

	Describe("IsPartialVersion", func() {
		It("Should return false for full version", func() {
			Expect(IsPartialVersion("6.8.1")).To(Equal(false))
		})

		It("Should return true for full version without patch", func() {
			Expect(IsPartialVersion("6.8")).To(Equal(true))
		})

		It("Should return true for full version without minor", func() {
			Expect(IsPartialVersion("6")).To(Equal(true))
		})
	})

	Describe("Semverify", func() {
		It("Shouldn't do anything for valid version", func() {
			Expect(Semverify("6.8.1")).To(Equal("6.8.1"))
		})

		It("Should add `.0` to the end there", func() {
			Expect(Semverify("6.8")).To(Equal("6.8.0"))
		})

		It("Should add `.0` to between text and numbers", func() {
			Expect(Semverify("1.7beta")).To(Equal("1.7.0-beta"))
		})

		It("Should add `.0` to between text with numbers and numbers", func() {
			Expect(Semverify("1.7beta1")).To(Equal("1.7.0-beta1"))
		})
	})

	Describe("Unsemverify", func() {
		It("Shouldn't do anything for versions without nil patch version", func() {
			Expect(Unsemverify("6.8.1")).To(Equal("6.8.1"))
		})

		It("Should remove `.0` from the end there", func() {
			Expect(Unsemverify("6.8.0")).To(Equal("6.8"))
		})

		It("Should remove `.0` from between text and numbers", func() {
			Expect(Unsemverify("1.7.0-beta")).To(Equal("1.7beta"))
		})

		It("Should remove `.0` from between text with numbers and numbers", func() {
			Expect(Unsemverify("1.7.0-beta1")).To(Equal("1.7beta1"))
		})
	})
})
