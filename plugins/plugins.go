package plugins

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/cavaliercoder/grab"
	"github.com/chuckpreslar/emission"
	"github.com/kardianos/osext"
	"github.com/markelog/archive"
	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/initiate"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"

	// plugins
	"github.com/markelog/eclectica/plugins/elm"
	"github.com/markelog/eclectica/plugins/golang"
	"github.com/markelog/eclectica/plugins/nodejs"
	"github.com/markelog/eclectica/plugins/python"
	"github.com/markelog/eclectica/plugins/ruby"
	"github.com/markelog/eclectica/plugins/rust"
)

var (
	Plugins = []string{
		"node",
		"rust",
		"ruby",
		"go",
		"python",
		"elm",
	}
)

type Pkg interface {
	PreDownload() error
	PreInstall() error
	Install() error
	PostInstall() error
	Switch() error
	Link() error
	Events() *emission.Emitter
	Environment() ([]string, error)
	ListRemote() ([]string, error)
	Info() map[string]string
	Bins() []string
	Dots() []string
}

type Plugin struct {
	name    string
	Version string
	Pkg     Pkg
	emitter *emission.Emitter
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
		emitter: emission.NewEmitter(),
	}

	switch {
	case name == "node":
		plugin.Pkg = &nodejs.Node{
			Version: version,
			Emitter: plugin.emitter,
		}
	case name == "rust":
		plugin.Pkg = &rust.Rust{
			Version: version,
			Emitter: plugin.emitter,
		}
	case name == "ruby":
		plugin.Pkg = &ruby.Ruby{
			Version: version,
			Emitter: plugin.emitter,
		}
	case name == "go":
		plugin.Pkg = &golang.Golang{
			Version: version,
			Emitter: plugin.emitter,
		}
	case name == "python":
		plugin.Pkg = &python.Python{
			Version: version,
			Emitter: plugin.emitter,
		}
	case name == "elm":
		plugin.Pkg = &elm.Elm{
			Version: version,
			Emitter: plugin.emitter,
		}
	}

	if len(args) == 2 {
		plugin.info, _ = plugin.Info()
	}

	return plugin
}

func (plugin *Plugin) PreDownload() error {
	return plugin.Pkg.PreDownload()
}

func (plugin *Plugin) PreInstall() error {
	return plugin.Pkg.PreInstall()
}

func (plugin *Plugin) LocalInstall() (err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	if plugin.Version == "" {
		return errors.New("Version was not defined")
	}

	var (
		version = fmt.Sprintf(".%s-version", plugin.name)
		path    = filepath.Join(pwd, version)
	)

	// Handle CTRL+C signal
	plugin.Interrupt()

	init := initiate.New(plugin.name, Plugins)
	init.CheckShell()

	err = init.Initiate()
	if err != nil {
		return
	}

	// If this is already a current version we can safely say this one is installed
	if plugin.Version == plugin.Current() {
		return nil
	}

	// If it was already installed, just switch and bail out
	if plugin.IsInstalled() {
		err = plugin.Switch()
		if err != nil {
			return
		}

		return nil
	}

	err = plugin.Finish()
	if err != nil {
		return
	}

	err = io.WriteFile(path, plugin.Version)
	if err != nil {
		plugin.Rollback()
		return
	}

	init.RestartShell()

	return
}

func (plugin *Plugin) Install() (err error) {
	err = plugin.Pkg.PreInstall()
	if err != nil {
		return
	}

	if plugin.Version == "" {
		return errors.New("Version was not defined")
	}

	// If this is already a current version we can safely say this one is installed
	if plugin.Version == plugin.Current() {
		return nil
	}

	// Handle CTRL+C signal
	plugin.Interrupt()

	init := initiate.New(plugin.name, Plugins)
	init.CheckShell()

	err = init.Initiate()
	if err != nil {
		return
	}

	// If it was already installed, just switch @current link if needed
	if plugin.IsInstalled() {
		err = plugin.Switch()
		if err != nil {
			return
		}

		return plugin.Link()
	}

	err = plugin.Finish()
	if err != nil {
		return
	}

	err = plugin.Switch()
	if err != nil {
		return
	}

	err = plugin.Link()
	if err != nil {
		return
	}

	plugin.emitter.Emit("done")

	// Start new shell from eclectica if needed
	// note: should be the last action
	init.RestartShell()

	return
}

func (plugin *Plugin) PostInstall() (err error) {
	err = plugin.Proxy()
	if err != nil {
		plugin.Rollback()
		return
	}

	err = plugin.Pkg.PostInstall()
	if err != nil {
		plugin.Rollback()
		return
	}

	variables.WriteVersion(plugin.name, plugin.Version)

	return
}

func (plugin *Plugin) Switch() (err error) {
	err = plugin.Pkg.Switch()
	if err != nil {
		plugin.Rollback()
		return
	}

	return
}

func (plugin *Plugin) Environment() ([]string, error) {
	return plugin.Pkg.Environment()
}

func (plugin *Plugin) Info() (map[string]string, error) {
	if plugin.Version == "" {
		return nil, errors.New("Version was not defined")
	}

	info := plugin.Pkg.Info()
	tmpDir := variables.TempDir()

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

// Current returns current used version
func (plugin *Plugin) Current() string {
	return variables.CurrentVersion(plugin.name)
}

// Rollback places everything back
func (plugin *Plugin) Rollback() {
	path := variables.Path(plugin.name, plugin.Version)
	os.RemoveAll(path)

	plugin.emitter.Emit("done")
}

// Finish finishes installation
func (plugin *Plugin) Finish() (err error) {
	err = plugin.Pkg.Install()
	if err != nil {
		plugin.Rollback()
		return
	}

	return plugin.PostInstall()
}

// Interrupt handles interruption signals (like CTRL+C)
func (plugin *Plugin) Interrupt() {
	channel := make(chan os.Signal, 1)

	signal.Notify(channel, os.Interrupt)

	go func() {
		<-channel
		plugin.Rollback()
		os.Exit(1)
	}()
}

func (plugin *Plugin) List() (vers []string, err error) {
	path := variables.Prefix(plugin.name)
	vers = io.ListVersions(path)

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

func (plugin *Plugin) Remove() (err error) {
	if plugin.Version == "" {
		return errors.New("Version was not defined")
	}

	var (
		home = filepath.Join(variables.Home(), plugin.name)
		base = filepath.Join(home, plugin.Version)
	)

	// Need to remove proxies if this is a current version.
	// So we wouldn't confuse user
	if plugin.Current() == plugin.Version {
		err = plugin.removeProxy()
		if err != nil {
			return err
		}
	}

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

	if resp == nil {
		return resp, errors.New("Something went wrong with HTTP request")
	}

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

func (plugin *Plugin) Events() *emission.Emitter {
	return plugin.Pkg.Events()
}

func (plugin *Plugin) Link() (err error) {
	var (
		base    = variables.Path(plugin.name, plugin.Version)
		current = variables.Path(plugin.name)
	)

	err = io.Symlink(current, base)
	if err != nil {
		return
	}

	return plugin.Pkg.Link()
}

// IsInstalled checks if this version was already installed
func (plugin *Plugin) IsInstalled() bool {
	return variables.IsInstalled(plugin.name, plugin.Version)
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
