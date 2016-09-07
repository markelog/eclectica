package request_test

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/jarcoal/httpmock"

	"github.com/bouk/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/request"
)

var _ = Describe("request", func() {
	Describe("Body", func() {
		var (
			body string
			err  error
		)

		Describe("fail", func() {
			Describe("fail misarably", func() {
				var guardExtract *monkey.PatchGuard
				client := &http.Client{}
				htype := reflect.TypeOf(client)

				BeforeEach(func() {
					guardExtract = monkey.PatchInstanceMethod(htype, "Get", func(*http.Client, string) (*http.Response, error) {

						return nil, errors.New("Weird error")
					})
				})

				AfterEach(func() {
					guardExtract.Unpatch()
				})

				It("should respond with error and don't panic", func() {
					_, err = Body("https://somewhere")

					Expect(err).Should(MatchError("Weird error"))
				})
			})

			BeforeEach(func() {
				httpmock.Activate()
			})

			AfterEach(func() {
				defer httpmock.DeactivateAndReset()
			})

			It("should return an error", func() {
				httpmock.RegisterResponder(
					"GET",
					"https://somewhere",
					httpmock.NewStringResponder(500, ""),
				)

				_, err = Body("https://somewhere")

				Expect(err).Should(MatchError("Can't establish connection"))
			})
		})

		Describe("success", func() {
			BeforeEach(func() {
				httpmock.Activate()

				httpmock.RegisterResponder(
					"GET",
					"https://somewhere",
					httpmock.NewStringResponder(200, "yey"),
				)
			})

			AfterEach(func() {
				defer httpmock.DeactivateAndReset()
			})

			BeforeEach(func() {
				body, err = Body("https://somewhere")
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("gets response", func() {
				Expect(body).To(Equal("yey"))
			})
		})
	})
})
