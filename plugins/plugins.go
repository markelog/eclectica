package plugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/cavaliercoder/grab"
	"github.com/kardianos/osext"
	"github.com/markelog/archive"
	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/console"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"

	// plugins
	"github.com/markelog/eclectica/plugins/golang"
	"github.com/markelog/eclectica/plugins/nodejs"
	"github.com/markelog/eclectica/plugins/ruby"
	"github.com/markelog/eclectica/plugins/rust"
)

var (
	Plugins = []string{
		"node",
		"rust",
		"ruby",
		"go",
	}
)

type Pkg interface {
	Bins() []string
	Install(string) error
	Environment(string) (string, error)
	PostInstall(string) (bool, error)
	ListRemote() ([]string, error)
	Info(string) (map[string]string, error)
	Current() string
}

type Plugin struct {
	name    string
	version string
	Pkg     Pkg
	info    map[string]string
}

func New(args ...string) *Plugin {
	var pkg Pkg
	var version string
	name := args[0]

	if len(args) == 2 {
		version = args[1]
	} else {
		version = ""
	}

	plugin := &Plugin{
		name:    name,
		version: version,
	}

	switch {
	case name == "node":
		pkg = &nodejs.Node{}
	case name == "rust":
		pkg = &rust.Rust{}
	case name == "ruby":
		pkg = &ruby.Ruby{}
	case name == "go":
		pkg = &golang.Golang{}
	}

	plugin.Pkg = pkg

	if len(args) == 2 {
		info, _ := plugin.Info()
		plugin.info = info
	}

	return plugin
}

func SearchBin(name string) string {
	bins := map[string][]string{
		"rust": rust.Bins,
		"go":   golang.Bins,
		"node": nodejs.Bins,
		"ruby": ruby.Bins,
	}

	for language, _ := range bins {
		for _, bin := range bins[language] {
			if name == bin {
				return language
			}
		}
	}

	return ""
}

func (plugin *Plugin) CreateProxy() (err error) {
	ecProxyFolder := os.Getenv("EC_PROXY_PLACE")

	if ecProxyFolder == "" {
		ecProxyFolder, err = osext.ExecutableFolder()
		if err != nil {
			return
		}
	}

	executable := filepath.Join(ecProxyFolder, "ec-proxy")

	_, err = os.Stat(executable)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Can't find ec-proxy binary")
		}

		return err
	}

	// TODO: fix, hack for rust
	name := plugin.name
	if name == "rust" {
		name = "rustc"
	}

	bins := plugin.Pkg.Bins()

	for _, bin := range bins {
		languageExecutable := filepath.Join(variables.ExecutablePath(name), bin)

		err = cprf.Copy(executable, languageExecutable)
		if err != nil {
			return
		}
	}

	return nil
}

func (plugin *Plugin) LocalInstall() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(pwd, fmt.Sprintf(".%s-version", plugin.name))

	err = io.WriteFile(path, plugin.version)
	if err != nil {
		return err
	}

	return plugin.Install()
}

func (plugin *Plugin) Install() error {
	if plugin.version == "" {
		return errors.New("Version was not defined")
	}

	err := Initiate(Plugins)
	if err != nil {
		return err
	}

	base := variables.Path(plugin.name, plugin.version)

	_, err = io.CreateDir(base)
	if err != nil {
		return err
	}

	currentPath := filepath.Join(variables.Prefix(plugin.name), "current")
	os.RemoveAll(currentPath)

	err = plugin.Pkg.Install(plugin.version)
	if err != nil {
		return err
	}

	return plugin.PostInstall()
}

func (plugin *Plugin) PostInstall() (err error) {
	currentPath := filepath.Join(variables.Prefix(plugin.name), "current")

	os.Symlink(variables.Path(plugin.name, plugin.version), currentPath)

	err = plugin.CreateProxy()
	if err != nil {
		return err
	}

	showMessage, err := plugin.Pkg.PostInstall(plugin.version)
	if err != nil {
		return
	}

	if showMessage && variables.NeedToRestartShell(plugin.name) {
		printShellMessage(plugin.name)
	}

	if strings.Contains(os.Getenv("PATH"), variables.DefaultInstall) == false {
		console.Shell()
	}

	return
}

func (plugin *Plugin) Info() (map[string]string, error) {
	if plugin.version == "" {
		return nil, errors.New("Version was not defined")
	}

	info, err := plugin.Pkg.Info(plugin.version)
	if err != nil {
		return nil, err
	}

	// I am crying over here :/
	tmpDir := os.TempDir()
	if runtime.GOOS == "linux" {
		tmpDir += "/"
	}

	if _, ok := info["extension"]; ok == false {
		info["extension"] = "tar.gz"
	}

	if _, ok := info["unarchive-filename"]; ok == false {
		info["unarchive-filename"] = info["filename"]
	}

	info["name"] = plugin.name
	info["archive-folder"] = tmpDir
	info["archive-path"] = fmt.Sprintf("%s%s.%s", info["archive-folder"], info["filename"], info["extension"])

	info["destination-folder"] = filepath.Join(variables.Home(), plugin.name, plugin.version)

	return info, nil
}

func (plugin *Plugin) Current() string {
	return plugin.Pkg.Current()
}

func (plugin *Plugin) List() (versions []string, err error) {
	versions = []string{}
	path := filepath.Join(variables.Home(), plugin.name)

	folders, _ := ioutil.ReadDir(path)
	for _, folder := range folders {
		name := folder.Name()

		if name == "current" {
			continue
		}

		versions = append(versions, name)
	}

	if len(versions) == 0 {
		err = errors.New("There is no installed versions")
	}

	return
}

func (plugin *Plugin) ListRemote() (map[string][]string, error) {
	versions, err := plugin.Pkg.ListRemote()

	if err != nil {
		return nil, err
	}

	return Compose(versions), nil
}

func (plugin *Plugin) Remove(version string) error {
	var err error
	home := filepath.Join(variables.Home(), plugin.name)
	base := filepath.Join(home, version)

	err = os.RemoveAll(base)
	if err != nil {
		return err
	}

	return nil
}

func (plugin *Plugin) Download() (*grab.Response, error) {
	if plugin.version == "" {
		return nil, errors.New("Version was not defined")
	}

	// If already downloaded
	if _, err := os.Stat(plugin.info["destination-folder"]); err == nil {
		return nil, nil
	}

	response, err := grab.GetAsync(plugin.info["archive-folder"], plugin.info["url"])
	if err != nil {
		return nil, err
	}

	resp := <-response

	if resp.HTTPResponse.StatusCode == 404 {
		grab.NewClient().CancelRequest(resp.Request)

		return resp, errors.New("Incorrect version " + plugin.version)
	}

	return resp, nil
}

func (plugin *Plugin) Extract() error {
	if plugin.version == "" {
		return errors.New("Version was not defined")
	}

	// Create something folder with path something like this – /home/user/.eclectica/versions/go
	extractionPlace, err := io.CreateDir(variables.Prefix(plugin.name))
	if err != nil {
		return err
	}

	// Just in case archive was downloaded, but not extracted
	// i.e. this issue comes up in the second run.
	// Which means we will delete folder with path like this –
	// /home/user/.eclectica/versions/go1.7.1.linux-amd64
	os.RemoveAll(filepath.Join(extractionPlace, plugin.info["filename"]))

	// Now we will extract archive from path something
	// like this – /tmp/go1.7.1.linux-amd64.tar.gz
	// to 				 /home/user/.eclectica/versions/go
	//
	// Which will give us folder with path like this –
	// /home/user/.eclectica/versions/go/go1.7.1.linux-amd64.tar.gz
	//
	// or like this – /home/user/.eclectica/versions/go/go
	//
	// Depends under what name language devs archived their dist
	err = archive.Extract(plugin.info["archive-path"], extractionPlace)
	if err != nil {
		return err
	}

	// Now we will need get path /home/user/.eclectica/versions/go/go1.7.1.linux-amd64.tar.gz
	tmpPath := filepath.Join(extractionPlace, plugin.info["unarchive-filename"])

	// And path like this – /home/user/.eclectica/versions/go/1.7.1
	extractionPath := plugin.info["destination-folder"]

	// Then rename that `tmpPath` expected path
	err = os.Rename(tmpPath, extractionPath)
	if err != nil {
		return err
	}

	return nil
}
