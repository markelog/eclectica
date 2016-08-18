package plugins

import (
  "os"
  "io/ioutil"
  "fmt"
  "errors"

  "github.com/cavaliercoder/grab"
  "github.com/markelog/archive"

  "github.com/markelog/eclectica/variables"
  "github.com/markelog/eclectica/plugins/nodejs"
  "github.com/markelog/eclectica/plugins/rust"
)

var (
  Plugins = []string{
    "node",
    "rust",
  }
)

type Pkg interface {
  Install(string) error
  ListRemote() ([]string, error)
  Info(string) (map[string]string, error)
  Current() string
}

type Plugin struct {
  name string
  version string
  pkg Pkg
  info map[string]string
}

func New(args... string) *Plugin {
  var pkg Pkg
  var version string
  name := args[0]

  if len(args) == 2 {
    version = args[1]
  } else {
    version = ""
  }

  plugin := &Plugin{
    name: name,
    version: version,
  }

  switch {
    case name == "node":
      pkg = &nodejs.Node{}
    case name == "rust":
      pkg = &rust.Rust{}
  }

  plugin.pkg = pkg

  if len(args) == 2 {
    info, _ := plugin.Info()
    plugin.info = info
  }

  return plugin
}

func (plugin *Plugin) Install() error {
  if plugin.version == "" {
    return errors.New("Version was not defined")
  }

  return plugin.pkg.Install(plugin.version)
}

func (plugin *Plugin) Info() (map[string]string, error) {
  if plugin.version == "" {
    return nil, errors.New("Version was not defined")
  }

  info, err := plugin.pkg.Info(plugin.version)
  if err != nil {
    return nil, err
  }

  info["archive-folder"] = os.TempDir()
  info["archive-path"] = fmt.Sprintf("%s%s.tar.gz", info["archive-folder"], info["filename"])

  info["destination-folder"] = fmt.Sprintf("%s/%s/%s", variables.Home(), plugin.name, plugin.version)

  return info, nil
}

func (plugin *Plugin) Current() string {
  return plugin.pkg.Current()
}

func (plugin *Plugin) ListRemote() (map[string][]string, error) {
  versions, err := plugin.pkg.ListRemote()
  if err != nil {
    return nil, err
  }

  return ComposeVersions(versions), nil
}

func (plugin *Plugin) Remove(version string) error {
  var err error
  home := fmt.Sprintf("%s/%s", variables.Home(), plugin.name)
  base := fmt.Sprintf("%s/%s", home, version)

  err = os.RemoveAll(base)

  if err != nil {
    return err
  }

  return nil
}

func (plugin *Plugin) Activate() (err error) {
  err = plugin.Extract()
  if err != nil {
    return
  }

  err = plugin.Install()
  if err != nil {
    return
  }

  return
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

  extractionPlace, err := createDir(fmt.Sprintf("%s/%s", variables.Home(), plugin.name))
  if err != nil {
    return err
  }

  // Just in case archive was downloaded, but not extracted
  // i.e. is below steps have failed this issues comes up in the second run
  os.RemoveAll(fmt.Sprintf("%s/%s", extractionPlace, plugin.info["filename"]))

  err = archive.Extract(plugin.info["archive-path"], extractionPlace)
  if err != nil {
    return err
  }

  downloadPath := fmt.Sprintf("%s/%s", extractionPlace, plugin.info["filename"])
  extractionPath := plugin.info["destination-folder"]

  err = os.Rename(downloadPath, extractionPath)
  if err != nil {
    return err
  }

  return nil
}

func (plugin *Plugin) List() (versions []string) {
  versions = []string{}
  path := variables.Home() + "/" + plugin.name

  if _, err := os.Stat(path); os.IsNotExist(err) {
    return
  }

  folders, _ := ioutil.ReadDir(path)
  for _, folder := range folders {
    versions = append(versions, folder.Name())
  }

  return
}

func createDir(path string) (string, error) {
  err := os.MkdirAll(path, 0700)

  if err != nil {
    return "", err
  }

  return path, nil
}
