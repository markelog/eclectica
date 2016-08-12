package info

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/markelog/eclectica/cmd/info"
)

var _ = Describe("info", func() {
  Describe("GetLanguage", func() {
    It("should get language", func() {
      result, hasLanguage := info.GetLanguage([]string{"-r", "rust",})

      Expect(hasLanguage).To(Equal(true))
      Expect(result).To(Equal("rust"))
    })

    It("should get language in different sequence", func() {
      result, hasLanguage := info.GetLanguage([]string{"rust", "-r",})

      Expect(hasLanguage).To(Equal(true))
      Expect(result).To(Equal("rust"))
    })

    It("should not get non-existing language", func() {
      _, hasLanguage := info.GetLanguage([]string{"-r", "rustc@1.2.3",})

      Expect(hasLanguage).To(Equal(false))
    })
  })

  Describe("hasCommand", func() {
    It("should detect command", func() {
      hasCommand := info.HasCommand([]string{"-r", "ls",})

      Expect(hasCommand).To(Equal(true))
    })

    It("should detect command in different sequence", func() {
      hasCommand := info.HasCommand([]string{"ls", "-r"})

      Expect(hasCommand).To(Equal(true))
    })

    It("should not detect command", func() {
      hasCommand := info.HasCommand([]string{"-r", "rustc@1.2.3",})

      Expect(hasCommand).To(Equal(false))
    })
  })
})
