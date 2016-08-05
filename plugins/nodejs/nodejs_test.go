package nodejs_test

import (
  "regexp"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  ."github.com/markelog/eclectica/plugins/nodejs"
)

var _ = Describe("nodejs", func() {
  var (
    remotes []string
    err error
  )

  Describe("ListVersions", func() {
    BeforeEach(func() {
      remotes, err = ListVersions()
    })

    It("should not return an error", func() {
      Expect(err).To(BeNil())
    })

    It("should have correct version values", func() {
      rp := regexp.MustCompile("[[:digit:]]+\\.[[:digit:]]+\\.[[:digit:]]+$")

      for _, element := range remotes {
        Expect(rp.MatchString(element)).To(Equal(true))
      }
    })
  })

  Describe("ComposeVersions", func() {
    It("should compose versions", func() {
      compose := ComposeVersions([]string{"0.8.2", "4.4.7", "6.3.0", "6.4.2"})

      Expect(compose["0.x"]).To(Equal([]string{"0.8.2"}))
      Expect(compose["4.x"]).To(Equal([]string{"4.4.7"}))
      Expect(compose["6.x"]).To(Equal([]string{"6.3.0", "6.4.2"}))
    })
  })

  Describe("GetKeys", func() {
    It("should get version keys", func() {
      list := map[string][]string{"4.x": []string{}, "0.x": []string{"0.8.2"}}
      keys := GetKeys(list)

      Expect(keys[0]).To(Equal("0.x"))
      Expect(keys[1]).To(Equal("4.x"))
    })
  })

  Describe("GetElements", func() {
    FIt("should get version elements", func() {
      list := ComposeVersions([]string{"0.8.2", "4.4.7", "6.3.0", "6.4.2"})
      elements := GetElements("6.x", list)

      Expect(elements[0]).To(Equal("6.3.0"))
      Expect(elements[1]).To(Equal("6.4.2"))
    })
  })
})
