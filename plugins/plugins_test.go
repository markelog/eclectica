package plugins_test

import (
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

	"github.com/markelog/monkey"

	. "github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/shell"

	eIO "github.com/markelog/eclectica/io"
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
		It("find all bins", func() {
			for _, language := range Plugins {
				bins := New(&Args{
					Language: language,
				}).Bins()

				for _, elem := range bins {
					Expect(SearchBin(elem)).To(Equal(language))
				}
			}
		})
	})

	Describe("Remove", func() {
		var (
			list        = false
			current     = false
			osRemoveAll = false
		)

		resCurrent := "6.7.0"

		type Empty struct{}

		var resOsRemoveAll error

		resOsRemoveAll = nil

		osRemoveCount := 0

		var guardCurrent *monkey.PatchGuard
		var guardList *monkey.PatchGuard

		BeforeEach(func() {
			var d *Plugin

			pType := reflect.TypeOf(d)

			guardCurrent = monkey.PatchInstanceMethod(pType, "Current",
				func(plugin *Plugin) string {
					current = true
					return resCurrent
				},
			)

			guardList = monkey.PatchInstanceMethod(pType, "List",
				func(plugin *Plugin) []string {
					list = true
					return []string{"test"}
				},
			)

			monkey.Patch(os.RemoveAll, func(path string) error {
				osRemoveAll = true
				osRemoveCount = osRemoveCount + 1
				return resOsRemoveAll
			})

		})

		AfterEach(func() {
			current = false
			osRemoveAll = false

			resCurrent = ""
			resOsRemoveAll = nil

			osRemoveCount = 0

			monkey.Unpatch(os.RemoveAll)

			guardCurrent.Unpatch()
			guardList.Unpatch()
		})

		It("remove sequence", func() {
			resCurrent = ""

			New(&Args{
				Language: "node",
				Version:  "6.8.0",
			}).Remove()

			Expect(current).To(Equal(true))
			Expect(list).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))

			Expect(osRemoveCount).To(Equal(1))
		})

		It("remove current version", func() {
			resCurrent = "6.8.0"

			New(&Args{
				Language: "node",
				Version:  "6.8.0",
			}).Remove()

			Expect(current).To(Equal(true))
			Expect(list).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))

			// Since removeProxy might have a lot of os.RemoveAll's
			Expect(osRemoveCount).Should(BeNumerically(">", 1))
		})

		It("shouldn't remove without version", func() {
			err := New(&Args{
				Language: "node",
			}).Remove()

			Expect(err.Error()).To(Equal("version was not defined"))
		})
	})

	Describe("Install", func() {
		var (
			pluginSwitch = false
			current      = false
			postInstall  = false
			pkgInstall   = false
			osRemoveAll  = false
			osSymlink    = false
			osStat       = false
			isInstalled  = false
		)

		resCurrent := "6.7.0"

		type Empty struct{}

		var (
			resPluginsSwitch error
			resInstall       error
			resPostInstall   error
			resOsRemoveAll   error
			resOsSymlink     error
			resIsInstalled   bool
			resOsStat        os.FileInfo
		)

		resPluginsSwitch = nil
		resInstall = nil
		resPostInstall = nil
		resOsRemoveAll = nil
		resOsSymlink = nil
		resOsStat = nil
		resIsInstalled = false

		var (
			guardPluginSwitch *monkey.PatchGuard
			guardCurrent      *monkey.PatchGuard
			guardPostInstall  *monkey.PatchGuard
			guardPkgInstall   *monkey.PatchGuard
			guardIsInstalled  *monkey.PatchGuard
			guardInitiate     *monkey.PatchGuard
			guardCheckShell   *monkey.PatchGuard
		)

		BeforeEach(func() {
			var d *Plugin
			var n *nodejs.Node
			var i *shell.Shell

			pType := reflect.TypeOf(d)
			nodejsType := reflect.TypeOf(n)
			initType := reflect.TypeOf(i)

			guardInitiate = monkey.PatchInstanceMethod(initType, "Initiate",
				func(init *shell.Shell) error {
					return nil
				},
			)

			guardCheckShell = monkey.PatchInstanceMethod(initType, "Check",
				func(init *shell.Shell) {},
			)

			guardCurrent = monkey.PatchInstanceMethod(pType, "Current",
				func(plugin *Plugin) string {
					current = true
					return resCurrent
				},
			)

			guardIsInstalled = monkey.PatchInstanceMethod(pType, "IsInstalled",
				func(plugin *Plugin) bool {
					isInstalled = true
					return resIsInstalled
				},
			)

			guardPostInstall = monkey.PatchInstanceMethod(pType, "PostInstall",
				func(plugin *Plugin) error {
					postInstall = true
					return resPostInstall
				},
			)

			guardPluginSwitch = monkey.PatchInstanceMethod(nodejsType, "Switch",
				func(plugin *nodejs.Node) error {
					pluginSwitch = true
					return resPluginsSwitch
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

			// Nothing to check for this one
			monkey.Patch(eIO.WriteFile, func(path, content string) error {
				return nil
			})
		})

		AfterEach(func() {
			pluginSwitch = false
			current = false
			postInstall = false
			pkgInstall = false
			osRemoveAll = false
			osSymlink = false
			osStat = false
			isInstalled = false

			resCurrent = ""
			resPostInstall = nil
			resOsRemoveAll = nil
			resOsSymlink = nil
			resOsStat = nil
			resIsInstalled = false

			monkey.Unpatch(os.Symlink)
			monkey.Unpatch(os.Stat)
			monkey.Unpatch(os.RemoveAll)
			monkey.Unpatch(eIO.WriteFile)

			guardInitiate.Unpatch()
			guardCheckShell.Unpatch()

			guardPluginSwitch.Unpatch()
			guardPostInstall.Unpatch()
			guardCurrent.Unpatch()
			guardPkgInstall.Unpatch()
			guardIsInstalled.Unpatch()
		})

		It("install sequence for not installed version", func() {
			New(&Args{
				Language: "node",
				Version:  "6.8.0",
			}).Install()

			Expect(current).To(Equal(true))
			Expect(isInstalled).To(Equal(true))
			Expect(pluginSwitch).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))
			Expect(postInstall).To(Equal(true))
			Expect(pkgInstall).To(Equal(true))
			Expect(osSymlink).To(Equal(true))
		})

		It("install sequence for installed version", func() {
			resIsInstalled = true

			New(&Args{
				Language: "node",
				Version:  "6.8.0",
			}).Install()

			Expect(current).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))
			Expect(pluginSwitch).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))
			Expect(osSymlink).To(Equal(true))
			Expect(postInstall).To(Equal(false))
			Expect(pkgInstall).To(Equal(false))
		})

		It("install sequence for current version", func() {
			resCurrent = "6.8.0"

			New(&Args{
				Language: "node",
				Version:  "6.8.0",
			}).Install()

			Expect(current).To(Equal(true))
			Expect(pluginSwitch).To(Equal(false))
			Expect(isInstalled).To(Equal(true))
			Expect(osRemoveAll).To(Equal(false))
			Expect(osSymlink).To(Equal(false))
			Expect(postInstall).To(Equal(false))
			Expect(pkgInstall).To(Equal(false))
		})

		It("local install sequence", func() {
			New(&Args{
				Language: "node",
				Version:  "6.8.0",
			}).LocalInstall()

			Expect(current).To(Equal(true))
			Expect(isInstalled).To(Equal(true))
			Expect(pluginSwitch).To(Equal(true))
			Expect(osRemoveAll).To(Equal(false))
			Expect(postInstall).To(Equal(true))
			Expect(pkgInstall).To(Equal(true))
			Expect(osSymlink).To(Equal(false))
		})

		It("should not return early", func() {
			resOsStat, _ = os.Stat(os.Getenv("HOME"))

			New(&Args{
				Language: "node",
				Version:  "6.8.0",
			}).Install()

			Expect(current).To(Equal(true))
			Expect(osRemoveAll).To(Equal(true))
			Expect(osSymlink).To(Equal(true))
			Expect(osStat).To(Equal(true))
			Expect(postInstall).To(Equal(true))
			Expect(pkgInstall).To(Equal(true))
		})

		It("returns error if version was not defined", func() {
			plugin := New(&Args{
				Language: "node",
			})

			Expect(plugin.Install()).Should(MatchError("version was not defined"))
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

			plugin = New(&Args{
				Language: "node",
				Version:  "5.0.0",
			})
		})

		AfterEach(func() {
			monkey.Unpatch(variables.Home)
			os.RemoveAll(filepath.Join(versionsFolder, name))
		})

		It("returns error if version was not defined", func() {
			plugin := New(&Args{
				Language: "node",
			})
			Expect(plugin.Extract()).Should(MatchError("version was not defined"))
		})

		It("should extract language", func() {
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

			plugin = New(&Args{
				Language: "node",
				Version:  "5.0.0",
			})
		})

		AfterEach(func() {
			defer ts.Close()
			guard.Unpatch()
			os.RemoveAll(archivePath)
			os.RemoveAll(destFolder)
		})

		It("returns error if version was not defined", func() {
			plugin := New(&Args{
				Language: "node",
			})
			_, err := plugin.Download()
			Expect(err).Should(MatchError("version was not defined"))
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
				_, err := New(&Args{
					Language: "node",
					Version:  "5.0.0",
				}).Download()

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

			plugin = New(&Args{
				Language: "node",
				Version:  "5.0.0",
			})
		})

		AfterEach(func() {
			guard.Unpatch()
		})

		It("returns error if there is no installed versions", func() {
			vers := plugin.List()

			Expect(len(vers)).To(Equal(0))
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
				func(*nodejs.Node) map[string]string {
					return info
				},
			)

			plugin = New(&Args{
				Language: "node",
				Version:  "5.0.0",
			})
		})

		AfterEach(func() {
			guard.Unpatch()
		})

		It("returns error if version was not defined", func() {
			plugin := New(&Args{
				Language: "node",
			})
			_, err := plugin.Info()

			Expect(err).Should(MatchError("version was not defined"))
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
			info, _ := New(&Args{
				Language: "node",
				Version:  "5.0.0",
			}).Info()

			Expect(info["extension"]).To(Equal("test"))
		})

		It("should add extension if it was not defined by the plugin", func() {
			Expect(info["extension"]).To(Equal("tar.gz"))
		})

		It("should not add extension if it was defined by the plugin", func() {
			info["extension"] = "test"
			info, _ := New(&Args{
				Language: "node",
				Version:  "5.0.0",
			}).Info()

			Expect(info["extension"]).To(Equal("test"))
		})

		It("should add extension if it was not defined by the plugin", func() {
			Expect(info["extension"]).To(Equal("tar.gz"))
		})

		It("should not add extension if it was defined by the plugin", func() {
			info["extension"] = "test"
			info, _ := New(&Args{
				Language: "node",
				Version:  "5.0.0",
			}).Info()

			Expect(info["extension"]).To(Equal("test"))
		})
	})

	Describe("Rollback", func() {
		var (
			bins              = false
			lastRemoveAllPath = ""
			osRemoveCount     = 0

			d         *Plugin
			guardList *monkey.PatchGuard
			guardBins *monkey.PatchGuard
			pType     = reflect.TypeOf(d)
		)

		BeforeEach(func() {

			guardBins = monkey.PatchInstanceMethod(pType, "Bins",
				func(plugin *Plugin) []string {
					bins = true
					return []string{"test"}
				},
			)

			monkey.Patch(os.RemoveAll, func(path string) error {
				lastRemoveAllPath = path
				osRemoveCount = osRemoveCount + 1
				return nil
			})

		})

		AfterEach(func() {
			bins = false
			lastRemoveAllPath = ""
			osRemoveCount = 0

			monkey.Unpatch(os.RemoveAll)

			guardList.Unpatch()
			guardBins.Unpatch()
			guardList.Unpatch()
		})

		It("doesn't try to remove binaries when there is a version to fallback on", func() {
			version := "6.8.0"
			eventCalled := false

			guardList = monkey.PatchInstanceMethod(pType, "List",
				func(plugin *Plugin) []string {
					return []string{"test"}
				},
			)

			p := New(&Args{
				Language: "node",
				Version:  version,
			})
			p.Events().On("done", func() {
				eventCalled = true
			})

			p.Rollback()

			Expect(eventCalled).To(Equal(true))
			Expect(bins).To(Equal(false))
			Expect(osRemoveCount).To(Equal(1))
			Expect(lastRemoveAllPath).Should(ContainSubstring("versions/node/" + version))
		})

		It("removes binaries when there is no version to fallback on", func() {
			version := "6.8.0"
			eventCalled := false

			guardList = monkey.PatchInstanceMethod(pType, "List",
				func(plugin *Plugin) []string {
					return []string{}
				},
			)

			p := New(&Args{
				Language: "node",
				Version:  version,
			})
			p.Events().On("done", func() {
				eventCalled = true
			})

			p.Rollback()

			Expect(eventCalled).To(Equal(true))
			Expect(bins).To(Equal(true))
			Expect(osRemoveCount).To(Equal(2))
			Expect(lastRemoveAllPath).Should(ContainSubstring(".eclectica/bin/test"))
		})
	})
})
