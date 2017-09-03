package modules_test

import (
	"os/exec"
	"reflect"

	"github.com/bouk/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/plugins/nodejs/modules"
)

var _ = Describe("modules", func() {
	Describe("List", func() {
		It("should list modules", func() {
			cmd := &exec.Cmd{}

			monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
				return cmd
			})

			guard := monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Output", func(*exec.Cmd) ([]uint8, error) {
				output := `├── nodemon@1.11.0
├── npm@3.10.10
├── pm2@2.6.1
└── yarn@0.27.5`

				return []uint8(output), nil
			})

			modules := New("5.12.0", "6.11.2")
			packages, err := modules.List()

			monkey.Unpatch(exec.Command)
			guard.Unpatch()

			Expect(err).ShouldNot(HaveOccurred())

			Expect(packages).To(ContainElement("nodemon"))
			Expect(packages).To(ContainElement("pm2"))
		})
	})

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
