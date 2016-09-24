package plugins

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/cavaliercoder/grab"
	"github.com/kardianos/osext"
	"github.com/markelog/archive"
	"github.com/markelog/cprf"

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
	Install(string) error
	PostInstall() (bool, error)
	ListRemote() ([]string, error)
	Info(string) (map[string]string, error)
	Current() string
}

type Plugin struct {
	name    string
	version string
	pkg     Pkg
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

	plugin.pkg = pkg

	if len(args) == 2 {
		info, _ := plugin.Info()
		plugin.info = info
	}

	return plugin
}

func (plugin *Plugin) CreateProxy() (err error) {
	folder, err := osext.ExecutableFolder()
	if err != nil {
		return
	}

	executable := filepath.Join(folder, "ec-proxy")

	_, err = os.Stat(executable)
	if err != nil {
		return
	}

	languageExecutable := filepath.Join(folder, plugin.name)

	err = cprf.Copy(executable, languageExecutable)
	if err != nil {
		return
	}

	return nil
}

func (plugin *Plugin) LocalInstall() error {
	err := Initiate(plugin.name)
	if err != nil {
		return err
	}

	if plugin.version == "" {
		return errors.New("Version was not defined")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(pwd, fmt.Sprintf(".%s-version", plugin.name))

	err = io.WriteFile(path, plugin.version)
	if err != nil {
		return err
	}

	return nil
}

func (plugin *Plugin) Install() error {
	err := Initiate(plugin.name)
	if err != nil {
		return err
	}

	if plugin.version == "" {
		return errors.New("Version was not defined")
	}

	base := variables.Prefix(plugin.name)

	_, err = io.CreateDir(base)
	if err != nil {
		return err
	}

	for _, file := range variables.Files {
		_, err := io.CreateDir(filepath.Join(base, file))
		if err != nil {
			return err
		}
	}

	err = plugin.pkg.Install(plugin.version)
	if err != nil {
		return err
	}

	return plugin.PostInstall()
}

func (plugin *Plugin) PostInstall() (err error) {
	showMessage, err := plugin.pkg.PostInstall()
	if err != nil {
		return
	}

	err = plugin.CreateProxy()
	if err != nil {
		return err
	}

	if showMessage && variables.NeedToRestartShell(plugin.name) {
		printShellMessage(plugin.name)
	}

	return
}

func (plugin *Plugin) Info() (map[string]string, error) {
	if plugin.version == "" {
		return nil, errors.New("Version was not defined")
	}

	info, err := plugin.pkg.Info(plugin.version)
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
	return plugin.pkg.Current()
}

func (plugin *Plugin) List() (versions []string, err error) {
	versions = []string{}
	path := filepath.Join(variables.Home(), plugin.name)

	folders, _ := ioutil.ReadDir(path)
	for _, folder := range folders {
		versions = append(versions, folder.Name())
	}

	if len(versions) == 0 {
		err = errors.New("There is no installed versions")
	}

	return
}

func (plugin *Plugin) ListRemote() (map[string][]string, error) {
	versions, err := plugin.pkg.ListRemote()

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

	extractionPlace, err := io.CreateDir(filepath.Join(variables.Home(), plugin.name))
	if err != nil {
		return err
	}

	// Just in case archive was downloaded, but not extracted
	// i.e. this issue comes up in the second run
	os.RemoveAll(filepath.Join(extractionPlace, plugin.info["filename"]))

	err = archive.Extract(plugin.info["archive-path"], extractionPlace)
	if err != nil {
		return err
	}

	downloadPath := filepath.Join(extractionPlace, plugin.info["unarchive-filename"])
	extractionPath := plugin.info["destination-folder"]

	err = os.Rename(downloadPath, extractionPath)
	if err != nil {
		return err
	}

	return nil
}
