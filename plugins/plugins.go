package plugins

import (
  "os"
  "io/ioutil"
  "fmt"
  "errors"

  "github.com/cavaliercoder/grab"
  "github.com/markelog/archive"

  "github.com/markelog/eclectica/plugins/nodejs"
  "github.com/markelog/eclectica/plugins/rust"
  "github.com/markelog/eclectica/variables"
)

var (
  List = []string{
    "node",
    "rust",
  }
)

func Activate(info map[string]string) (err error) {
  err = Extract(info)
  if err != nil {
    return
  }

  err = Install(info)
  if err != nil {
    return
  }

  return
}

func Download(info map[string]string) (*grab.Response, error) {
  // If already downloaded
  if _, err := os.Stat(info["destination-folder"]); err == nil {
    if err != nil {
      return nil, err
    }

    return nil, nil
  }

  response, err := grab.GetAsync(info["archive-folder"], info["url"])
  if err != nil {
    return nil, err
  }

  resp := <-response

  if resp.HTTPResponse.StatusCode == 404 {
    return resp, errors.New("Incorrect version " + info["version"])
  }

  return resp, nil
}

func Extract(info map[string]string) error {
  extractionPlace, err := createDir(fmt.Sprintf("%s/%s", variables.Home(), info["name"]))
  if err != nil {
    return err
  }

  // Just in case archive was downloaded, but not extracted
  // i.e. is below steps have failed this issues comes up in the second run
  err = os.RemoveAll(fmt.Sprintf("%s/%s/%s", variables.Home(), info["name"], info["filename"]))

  err = archive.Extract(info["archive-path"], extractionPlace)
  if err != nil {
    return err
  }

  downloadPath := fmt.Sprintf("%s/%s", extractionPlace, info["filename"])
  extractionPath := info["destination-folder"]

  err = os.Rename(downloadPath, extractionPath)
  if err != nil {
    return err
  }

  return nil
}

func Versions(name string) (versions []string) {
  versions = []string{}
  path := variables.Home() + "/" + name

  if _, err := os.Stat(path); os.IsNotExist(err) {
    return
  }

  folders, _ := ioutil.ReadDir(path)
  for _, folder := range folders {
    versions = append(versions, folder.Name())
  }

  return
}

func Version(langauge, version string) (map[string]string, error) {
  var (
    info map[string]string
    err error
  )

  switch {
    case langauge == "node":
      info, err = nodejs.Version(version)
    case langauge == "rust":
      info, err = rust.Version(version)
  }

  return info, err
}

func Remove(langauge, version string) error {
  switch {
    case langauge == "node":
      return nodejs.Remove(version)
    case langauge == "rust":
      return rust.Remove(version)
  }

  return nil
}

func Install(data map[string]string) error {
  switch {
    case data["name"] == "node":
      return nodejs.Install(data)
    case data["name"] == "rust":
      return rust.Install(data)
  }

  return nil
}

func RemoteList(langauge string) (map[string][]string, error) {
  var (
    versions []string
    err error
  )

  switch {
    case langauge == "node":
      versions, err = nodejs.ListVersions()
    case langauge == "rust":
      versions, err = rust.ListVersions()
  }

  if err != nil {
    return nil, err
  }

  return ComposeVersions(versions), nil
}

func CurrentVersion(langauge string) string {
  switch {
    case langauge == "node":
      return nodejs.CurrentVersion()
    case langauge == "rust":
      return rust.CurrentVersion()
  }

  return ""
}

func createDir(path string) (string, error) {
  err := os.MkdirAll(path, 0700)

  if err != nil {
    return "", err
  }

  return path, nil
}
