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
			list := map[string][]string{
				"10.x": {},
				"4.x":  {},
				"0.x":  {"0.8.2"},
			}
			keys := GetKeys(list)

			Expect(keys[0]).To(Equal("10.x"))
			Expect(keys[1]).To(Equal("4.x"))
			Expect(keys[2]).To(Equal("0.x"))
		})

		It("should get version keys with 0.10.x in it", func() {
			list := map[string][]string{
				"0.9.x":  {},
				"0.10.x": {},
				"0.1.x":  {"0.8.2"},
			}
			keys := GetKeys(list)

			Expect(keys[0]).To(Equal("0.10.x"))
			Expect(keys[1]).To(Equal("0.9.x"))
			Expect(keys[2]).To(Equal("0.1.x"))
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

	Describe("Complete", func() {
		It("support for 'latest' keyword for semver version structure", func() {
			version := "latest"
			versions := []string{
				"6.1.0", "5.2.0", "6.2.0", "6.8.3", "7.7.0", "7.3.0",
			}

			test, _ := Complete(version, versions)

			Expect(test).To(Equal("7.7.0"))
		})

		It("support for 'latest' keyword for incomplete version structure", func() {
			version := "latest"
			versions := []string{
				"1.4.3", "1.5.1", "1.5.2", "1.5.3",
				"1.5.4", "1.5", "1.5beta1", "1.5beta2",
				"1.5beta3", "1.5rc1", "1.6.1", "1.6.2",
				"1.6.3", "1.6.4", "1.6", "1.6beta1",
				"1.6beta2", "1.6rc1", "1.6rc2",
				"1.7.1", "1.7.2", "1.7.3", "1.7.4",
				"1.7.5", "1.7", "1.7beta1", "1.7beta2",
				"1.7rc1", "1.7rc2", "1.7rc3", "1.7rc4",
				"1.7rc5", "1.7rc6", "1.8", "1.8beta1",
				"1.8beta2", "1.8rc1", "1.8rc2",
			}

			test, _ := Complete(version, versions)

			Expect(test).To(Equal("1.8.0"))
		})

		It("support for partial major version", func() {
			version := "6"
			versions := []string{
				"6.1.0", "5.2.0", "6.2.0", "6.8.3",
			}

			test, _ := Complete(version, versions)

			Expect(test).To(Equal("6.8.3"))
		})

		It("support for partial minor version", func() {
			version := "6.4"
			versions := []string{
				"6.1.0", "5.2.0", "6.2.0", "6.8.3", "6.4.2", "6.4.0",
			}

			test, _ := Complete(version, versions)

			Expect(test).To(Equal("6.4.2"))
		})

		It("shouldn't do anything for full version", func() {
			version := "6.1.1"
			versions := []string{
				"6.1.0", "5.2.0", "6.2.0", "6.8.3", "6.4.2", "6.4.0",
			}

			test, err := Complete(version, versions)

			Expect(err).To(BeNil())
			Expect(test).To(Equal("6.1.1"))
		})
	})

	Describe("IsPartial", func() {
		It("Should return true for 'latest' keyword", func() {
			Expect(IsPartial("latest")).To(Equal(true))
		})

		It("Should return false for full version", func() {
			Expect(IsPartial("6.8.1")).To(Equal(false))
		})

		It("Should return true for full version without patch", func() {
			Expect(IsPartial("6.8")).To(Equal(true))
		})

		It("Should return true for full version without minor", func() {
			Expect(IsPartial("6")).To(Equal(true))
		})
	})

	Describe("Semverify", func() {
		It("Shouldn't do anything for valid version", func() {
			Expect(Semverify("6.8.1")).To(Equal("6.8.1"))
		})

		It("Should add `.0` to the end there", func() {
			Expect(Semverify("6.8")).To(Equal("6.8.0"))
		})

		It("Should add two `0` to the end there", func() {
			Expect(Semverify("6")).To(Equal("6.0.0"))
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
