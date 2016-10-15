package plugins_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"syscall"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bouk/monkey"

	. "github.com/markelog/eclectica/plugins"

	"github.com/markelog/eclectica/plugins/nodejs"
	"github.com/markelog/eclectica/variables"
)

var _ = Describe("plugins", func() {
	var (
		name           string
		path           string
		version        string
		archivePath    string
		destFolder     string
		versionsFolder string
		url            string
		filename       string
		info           map[string]string
		plugin         *Plugin
	)

	Describe("SearchBin", func() {
		It("find all ruby bins", func() {
			bins := New("ruby").Bins()

			for _, elem := range bins {
				Expect(SearchBin(elem)).To(Equal("ruby"))
			}
		})

		It("find all node bins", func() {
			bins := New("node").Bins()

			for _, elem := range bins {
				Expect(SearchBin(elem)).To(Equal("node"))
			}
		})

		It("find all go bins", func() {
			bins := New("go").Bins()

			for _, elem := range bins {
				Expect(SearchBin(elem)).To(Equal("go"))
			}
		})

		It("find all rust bins", func() {
			bins := New("rust").Bins()

			for _, elem := range bins {
				Expect(SearchBin(elem)).To(Equal("rust"))
			}
		})
	})

	Describe("Install", func() {
		var (
			initiate    = false
			current     = false
			postInstall = false
			pkgInstall  = false
			osRemoveAll = false
			osSymlink   = false
			osStat      = false
		)

		resCurrent := "6.7.0"

		type Empty struct{}

		var resInitiate error
		var resInstall error
		var resPostInstall error
		var resPkgInstall error
		var resOsRemoveAll error
		var resOsSymlink error
		var resOsStat os.FileInfo

		resInitiate = nil
		resInstall = nil
		resPostInstall = nil
		resPkgInstall = nil
		resOsRemoveAll = nil
		resOsSymlink = nil
		resOsStat = nil

		var guardCurrent *monkey.PatchGuard
		var guardPostInstall *monkey.PatchGuard
		var guardPkgInstall *monkey.PatchGuard

		BeforeEach(func() {
			var d *Plugin
			var n *nodejs.Node

			pType := reflect.TypeOf(d)
			nodejsType := reflect.TypeOf(n)

			monkey.Patch(Initiate, func() error {
				initiate = true
				return resInitiate
			})

			guardCurrent = monkey.PatchInstanceMethod(pType, "Current",
				func(plugin *Plugin) string {
					current = true
					return resCurrent
				},
			)

			guardPostInstall = monkey.PatchInstanceMethod(pType, "PostInstall",
				func(plugin *Plugin) error {
					postInstall = true
					return resPostInstall
				},
			)

			guardPkgInstall = monkey.PatchInstanceMethod(nodejsType, "Install",
				func(plugin *nodejs.Node) error {
					pkgInstall = true
					return resInstall
				},
			)

			monkey.Patch(os.RemoveAll, func(path string) error {
				osRemoveAll = true
				return resOsRemoveAll
			})

			monkey.Patch(os.Symlink, func(from, to string) error {
				osSymlink = true
				return resOsSymlink
			})

			monkey.Patch(os.Stat, func(path string) (os.FileInfo, error) {
				osStat = true
				return resOsStat, nil
			})
		})

		AfterEach(func() {
			initiate = false
			current = false
			postInstall = false
			pkgInstall = false
			osRemoveAll = false
			osSymlink = false
			osStat = false

			resInitiate = nil
			resCurrent = ""
			resPostInstall = nil
			resPkgInstall = nil
			resOsRemoveAll = nil
			resOsSymlink = nil
			resOsStat = nil

			monkey.Unpatch(Initiate)
			monkey.Unpatch(os.Symlink)
			monkey.Unpatch(os.Stat)
			monkey.Unpatch(os.RemoveAll)

			guardPostInstall.Unpatch()
			guardCurrent.Unpatch()
			guardPkgInstall.Unpatch()
		})

		It("install sequence", func() {
			New("node", "6.8.0").Install()

			Expect(initiate).To(Equal(true))
			Expect(current).To(Equal(true))
			Expect(postInstall).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))
			Expect(osSymlink).To(Equal(true))
			Expect(osStat).To(Equal(true))
			Expect(pkgInstall).To(Equal(false))
		})

		It("should return early if current version is installed one", func() {
			resCurrent = "6.8.0"

			New("node", "6.8.0").Install()

			Expect(initiate).To(Equal(true))
			Expect(current).To(Equal(true))
			Expect(postInstall).To(Equal(true))
			Expect(osRemoveAll).To(Equal(false))
			Expect(osSymlink).To(Equal(false))
			Expect(osStat).To(Equal(false))
			Expect(pkgInstall).To(Equal(false))
		})

		It("should return early if it couldn't RemoveAll", func() {
			resOsRemoveAll = errors.New("something")

			New("node", "6.8.0").Install()

			Expect(initiate).To(Equal(true))
			Expect(current).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))
			Expect(postInstall).To(Equal(false))
			Expect(osSymlink).To(Equal(false))
			Expect(osStat).To(Equal(false))
			Expect(pkgInstall).To(Equal(false))
		})

		It("should not return early", func() {
			resOsStat, _ = os.Stat(os.Getenv("HOME"))

			New("node", "6.8.0").Install()

			Expect(initiate).To(Equal(true))
			Expect(current).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))
			Expect(osSymlink).To(Equal(true))
			Expect(osStat).To(Equal(true))
			Expect(postInstall).To(Equal(true))
			Expect(pkgInstall).To(Equal(false))
		})

		It("returns error if version was not defined", func() {
			plugin := New("node")

			Expect(plugin.Install()).Should(MatchError("Version was not defined"))
		})
	})

	Describe("Extract", func() {
		BeforeEach(func() {
			path, _ = filepath.Abs("../testdata/plugins")
			versionsFolder, _ = filepath.Abs("../testdata/plugins/versions")
			name = "node"
			version = "5.0.0"
			filename = "node-arch"
			destFolder, _ = filepath.Abs("../testdata/plugins/versions/" + name + "/" + version)
			archivePath = filepath.Join(path, filename+".tar.gz")

			info = map[string]string{
				"name":               name,
				"version":            version,
				"archive-path":       archivePath,
				"destination-folder": destFolder,
				"filename":           filename,
				"unarchive-filename": filename,
			}

			monkey.Patch(variables.Home, func() string {
				return versionsFolder
			})

			var d *Plugin
			ptype := reflect.TypeOf(d)

			var guard *monkey.PatchGuard
			guard = monkey.PatchInstanceMethod(ptype, "Info",
				func(plugin *Plugin) (map[string]string, error) {
					guard.Unpatch()
					return info, nil
				})

			plugin = New("node", "5.0.0")
		})

		AfterEach(func() {
			monkey.Unpatch(variables.Home)
			os.RemoveAll(filepath.Join(versionsFolder, name))
		})

		It("returns error if version was not defined", func() {
			plugin := New("node")
			Expect(plugin.Extract()).Should(MatchError("Version was not defined"))
		})

		It("should extract langauge", func() {
			plugin.Extract()

			_, err := os.Stat(filepath.Join(destFolder, "/test.txt"))
			Expect(err).To(BeNil())
		})

		It("should extract even if previous archive was downloaded, but not extracted", func() {
			failedAttempt := filepath.Join(versionsFolder, name, filename)

			os.MkdirAll(failedAttempt, 0700)

			plugin.Extract()

			_, err := os.Stat(destFolder + "/test.txt")
			Expect(err).To(BeNil())

			_, err = os.Stat(failedAttempt)
			Expect(err).ShouldNot(BeNil())
		})
	})

	Describe("Download", func() {
		var (
			guard *monkey.PatchGuard
			ts    *httptest.Server
		)

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				status := 200

				if _, ok := r.URL.Query()["status"]; ok {
					fmt.Sscanf(r.URL.Query().Get("status"), "%d", &status)
				}

				w.WriteHeader(status)
				io.WriteString(w, "test")
			}))

			path, _ = filepath.Abs("../testdata/plugins")
			filename = "node-v5.0.0-darwin-x64.tar.gz"
			archivePath = filepath.Join(path, filename)
			destFolder = filepath.Join(path, filename)
			url = ts.URL + "/" + filename

			info = map[string]string{
				"name":               "node",
				"version":            "5.0.0",
				"destination-folder": destFolder,
				"archive-folder":     path,
				"archive-path":       archivePath,
				"url":                url,
			}

			monkey.Patch(variables.Home, func() string {
				return versionsFolder
			})

			var d *Plugin
			ptype := reflect.TypeOf(d)

			guard = monkey.PatchInstanceMethod(ptype, "Info",
				func(*Plugin) (map[string]string, error) {
					return info, nil
				},
			)

			plugin = New("node", "5.0.0")
		})

		AfterEach(func() {
			defer ts.Close()
			guard.Unpatch()
			os.RemoveAll(archivePath)
			os.RemoveAll(destFolder)
		})

		It("returns error if version was not defined", func() {
			plugin := New("node")
			_, err := plugin.Download()
			Expect(err).Should(MatchError("Version was not defined"))
		})

		Describe("200 response", func() {
			It("should download tar", func() {
				plugin.Download()

				_, err := os.Stat(archivePath)
				Expect(err).To(BeNil())
			})

			It("should not download anything if file already exist", func() {
				os.MkdirAll(destFolder, 0777)

				response, _ := plugin.Download()
				Expect(response).To(BeNil())
			})
		})

		Describe("404 response", func() {
			It("should return error", func() {
				info["url"] += "?status=404"
				_, err := New("node", "5.0.0").Download()

				Expect(err).Should(MatchError("Incorrect version 5.0.0"))
			})
		})
	})

	Describe("List", func() {
		var guard *monkey.PatchGuard

		BeforeEach(func() {
			guard = monkey.Patch(os.Stat, func(path string) (os.FileInfo, error) {
				return nil, syscall.ENOENT
			})

			plugin = New("node", "5.0.0")
		})

		AfterEach(func() {
			guard.Unpatch()
		})

		It("returns error if there is no installed versions", func() {
			_, err := plugin.List()

			Expect(err).Should(MatchError("There is no installed versions"))
		})
	})

	Describe("Info", func() {
		var guard *monkey.PatchGuard

		BeforeEach(func() {
			info = map[string]string{
				"filename": "node-arch",
				"url":      url,
			}

			var d *nodejs.Node
			ptype := reflect.TypeOf(d)

			guard = monkey.PatchInstanceMethod(ptype, "Info",
				func(*nodejs.Node) (map[string]string, error) {
					return info, nil
				},
			)

			plugin = New("node", "5.0.0")
		})

		AfterEach(func() {
			guard.Unpatch()
		})

		It("returns error if version was not defined", func() {
			plugin := New("node")
			_, err := plugin.Info()

			Expect(err).Should(MatchError("Version was not defined"))
		})

		It("should augment output from plugin `Info` method", func() {
			info, _ := plugin.Info()

			tmpDir := os.TempDir()
			if runtime.GOOS == "linux" {
				tmpDir += "/"
			}

			Expect(info["name"]).To(Equal("node"))
			Expect(info["version"]).To(Equal("5.0.0"))
			Expect(info["archive-folder"]).To(Equal(tmpDir))
			Expect(info["archive-path"]).To(Equal(tmpDir + "node-arch.tar.gz"))
			Expect(info["destination-folder"]).To(Equal(variables.Home() + "/node/5.0.0"))
		})

		It("should not add extension if it was defined by the plugin", func() {
			info["extension"] = "test"
			info, _ := New("node", "5.0.0").Info()

			Expect(info["extension"]).To(Equal("test"))
		})

		It("should add extension if it was not defined by the plugin", func() {
			Expect(info["extension"]).To(Equal("tar.gz"))
		})

		It("should not add extension if it was defined by the plugin", func() {
			info["extension"] = "test"
			info, _ := New("node", "5.0.0").Info()

			Expect(info["extension"]).To(Equal("test"))
		})

		It("should add extension if it was not defined by the plugin", func() {
			Expect(info["extension"]).To(Equal("tar.gz"))
		})

		It("should not add extension if it was defined by the plugin", func() {
			info["extension"] = "test"
			info, _ := New("node", "5.0.0").Info()

			Expect(info["extension"]).To(Equal("test"))
		})
	})

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

			Expect(keys[0]).To(Equal("0.x"))
			Expect(keys[1]).To(Equal("4.x"))
		})
	})

	Describe("GetElements", func() {
		It("should get version elements", func() {
			list := Compose([]string{"0.8.2", "4.4.7", "6.3.0", "6.4.2"})
			elements := GetElements("6.x", list)

			Expect(elements[0]).To(Equal("6.3.0"))
			Expect(elements[1]).To(Equal("6.4.2"))
		})
	})
})
