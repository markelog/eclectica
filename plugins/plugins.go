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
	"github.com/markelog/eclectica/versions"

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
	Install() error
	Environment() (string, error)
	PostInstall() error
	ListRemote() ([]string, error)
	Info() (map[string]string, error)
	Bins() []string
	Dots() []string
	Current() string
}

type Plugin struct {
	name    string
	Version string
	Pkg     Pkg
	info    map[string]string
}

func New(args ...string) *Plugin {
	var (
		version string
		name    = args[0]
	)

	if len(args) == 2 {
		version = args[1]
	} else {
		version = ""
	}

	plugin := &Plugin{
		name:    name,
		Version: version,
	}

	switch {
	case name == "node":
		plugin.Pkg = &nodejs.Node{
			Version: version,
		}
	case name == "rust":
		plugin.Pkg = &rust.Rust{
			Version: version,
		}
	case name == "ruby":
		plugin.Pkg = &ruby.Ruby{
			Version: version,
		}
	case name == "go":
		plugin.Pkg = &golang.Golang{
			Version: version,
		}
	}

	if len(args) == 2 {
		plugin.info, _ = plugin.Info()
	}

	return plugin
}

func (plugin *Plugin) LocalInstall() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(pwd, fmt.Sprintf(".%s-version", plugin.name))

	err = io.WriteFile(path, plugin.Version)
	if err != nil {
		return err
	}

	return plugin.Install()
}

func (plugin *Plugin) Install() error {
	if plugin.Version == "" {
		return errors.New("Version was not defined")
	}

	err := Initiate()
	if err != nil {
		return err
	}

	// If this is already a current version we don't need to do anything
	if plugin.Version == plugin.Current() {
		return plugin.PostInstall()
	}

	bin := variables.GetBin(plugin.name, plugin.Version)

	// If binary for this plugin already exist then we can assume it was installed before;
	// which means we can bail at this point
	if _, err := os.Stat(bin); err == nil {
		return plugin.PostInstall()
	}

	err = plugin.Pkg.Install()
	if err != nil {
		return err
	}

	return plugin.PostInstall()
}

func (plugin *Plugin) PostInstall() (err error) {
	var (
		base    = variables.Path(plugin.name, plugin.Version)
		current = variables.Path(plugin.name)
	)

	err = symlink(current, base)
	if err != nil {
		return err
	}

	err = plugin.Proxy()
	if err != nil {
		return err
	}

	err = plugin.Pkg.PostInstall()
	if err != nil {
		return
	}

	// Start new shell from eclectica if needed
	StartShell()

	return
}

func (plugin *Plugin) Environment() (string, error) {
	return plugin.Pkg.Environment()
}

func (plugin *Plugin) Info() (map[string]string, error) {
	if plugin.Version == "" {
		return nil, errors.New("Version was not defined")
	}

	info, err := plugin.Pkg.Info()
	if err != nil {
		return nil, err
	}

	// I am crying over here :/
	tmpDir := os.TempDir()
	if runtime.GOOS == "linux" {
		tmpDir += "/"
	}

	if _, ok := info["name"]; ok == false {
		info["name"] = plugin.name
	}

	if _, ok := info["version"]; ok == false {
		info["version"] = plugin.Version
	}

	if _, ok := info["extension"]; ok == false {
		info["extension"] = "tar.gz"
	}

	if _, ok := info["unarchive-filename"]; ok == false {

		// Notice different value
		info["unarchive-filename"] = info["filename"]
	}

	if _, ok := info["destination-folder"]; ok == false {
		info["destination-folder"] = filepath.Join(variables.Home(), plugin.name, plugin.Version)
	}

	if _, ok := info["archive-folder"]; ok == false {
		info["archive-folder"] = tmpDir
	}

	if _, ok := info["archive-path"]; ok == false {
		info["archive-path"] = fmt.Sprintf("%s%s.%s", info["archive-folder"], info["filename"], info["extension"])
	}

	return info, nil
}

func (plugin *Plugin) Current() string {
	return plugin.Pkg.Current()
}

func (plugin *Plugin) List() (vers []string, err error) {
	vers = []string{}
	path := filepath.Join(variables.Home(), plugin.name)

	folders, _ := ioutil.ReadDir(path)
	for _, folder := range folders {
		name := folder.Name()

		if name == "current" {
			continue
		}

		vers = append(vers, name)
	}

	if len(vers) == 0 {
		err = errors.New("There is no installed versions")
	}

	return
}

func (plugin *Plugin) ListRemote() (map[string][]string, error) {
	vers, err := plugin.Pkg.ListRemote()

	if err != nil {
		return nil, err
	}

	return versions.Compose(vers), nil
}

func (plugin *Plugin) Remove(version string) error {
	var err error
	home := filepath.Join(variables.Home(), plugin.name)
	base := filepath.Join(home, version)

	err = os.RemoveAll(base)
	if err != nil {
		return err
	}

	_, err = plugin.List()
	if err == nil {
		return nil
	}

	if err.Error() != "There is no installed versions" {
		return err
	}

	err = plugin.removeProxy()
	if err != nil {
		return err
	}

	return os.RemoveAll(variables.Prefix(plugin.name))
}

func (plugin *Plugin) Download() (*grab.Response, error) {
	if plugin.Version == "" {
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

		return resp, errors.New("Incorrect version " + plugin.Version)
	}

	return resp, nil
}

func (plugin *Plugin) Extract() error {
	if plugin.Version == "" {
		return errors.New("Version was not defined")
	}

	// Create language folder with path like this – /home/user/.eclectica/versions/go
	extractionPlace, err := io.CreateDir(variables.Prefix(plugin.name))
	if err != nil {
		return err
	}

	// Just in case archive was downloaded, but not extracted
	// i.e. this issue comes up at the second run.
	// Which means we will delete folder with path like this –
	// /home/user/.eclectica/versions/go1.7.1.linux-amd64
	os.RemoveAll(filepath.Join(extractionPlace, plugin.info["filename"]))

	// Now we will extract archive from
	// path like this - 					 /tmp/go1.7.1.linux-amd64.tar.gz
	// inside folder like this -   /home/user/.eclectica/versions/go
	//
	// Which will give us path like this –
	// /home/user/.eclectica/versions/go/go1.7.1.linux-amd64
	//
	// or like this – /home/user/.eclectica/versions/go/go
	//
	// Depends under what name language devs archived their dist
	err = archive.Extract(plugin.info["archive-path"], extractionPlace)
	if err != nil {
		return err
	}

	// Now we will need get path - /home/user/.eclectica/versions/go/go1.7.1.linux-amd64
	tmpPath := filepath.Join(extractionPlace, plugin.info["unarchive-filename"])

	// And get path like this – /home/user/.eclectica/versions/go/1.7.1
	extractionPath := plugin.info["destination-folder"]

	// Clean up in case user extracts already extracted version.
	// That might happen if this is the second pass and in first time we errored somewhere above
	os.RemoveAll(extractionPath)

	// Then rename such tmp path to what we expected to see
	err = os.Rename(tmpPath, extractionPath)
	if err != nil {
		return err
	}

	return nil
}

func (plugin *Plugin) Bins() []string {
	return plugin.Pkg.Bins()
}

func (plugin *Plugin) Dots() []string {
	return plugin.Pkg.Dots()
}

func (plugin *Plugin) Proxy() (err error) {
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
			err = errors.New("Can't find ec-proxy binary")
		}

		return err
	}

	bins := plugin.Bins()

	for _, bin := range bins {
		err = cprf.Copy(executable, variables.DefaultInstall)
		if err != nil {
			return
		}

		fullProxy := filepath.Join(variables.DefaultInstall, "ec-proxy")
		fullBin := filepath.Join(variables.DefaultInstall, bin)

		err = os.Rename(fullProxy, fullBin)
		if err != nil {
			return
		}
	}

	return nil
}

func (plugin *Plugin) removeProxy() (err error) {
	bins := plugin.Bins()

	for _, bin := range bins {
		proxy := filepath.Join(variables.DefaultInstall, bin)

		err = os.RemoveAll(proxy)
		if err != nil {
			return
		}
	}

	return nil
}

func (plugin *Plugin) SetFullVersion(vers []string) error {
	if plugin.Version == "" {
		return errors.New("Version was not defined")
	}

	if versions.IsPartialVersion(plugin.Version) == false {
		return nil
	}

	// This shouldn't happen
	if len(vers) == 0 {
		return errors.New("No versions available")
	}

	last, err := versions.GetLatest(plugin.Version, vers)

	if err != nil {
		return err
	}

	plugin.Version = last
	plugin.Pkg = New(plugin.name, last).Pkg
	plugin.info, _ = plugin.Info()

	return nil
}

func SearchBin(name string) string {
	bins := map[string][]string{}

	for _, language := range Plugins {
		bins[language] = New(language).Bins()
	}

	for index, _ := range bins {
		for _, bin := range bins[index] {
			if name == bin {
				return index
			}
		}
	}

	return ""
}

// This one exists only to support nvm's `.nvmrc`
func Dots(name string) []string {
	return New(name).Dots()
}

func symlink(current, base string) (err error) {
	// Remove current@ symlink if it existed for previous version
	err = os.RemoveAll(current)
	if err != nil {
		return err
	}

	// Set up current@ symlink
	err = os.Symlink(base, current)
	if err != nil {
		return err
	}

	return nil
}
