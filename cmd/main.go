package main

import (
  "os"
  "fmt"
  "time"
  "io/ioutil"
  // "strings"

  "github.com/markelog/archive"
  "github.com/urfave/cli"
  "github.com/cavaliercoder/grab"

  "github.com/markelog/eclectica/variables"
  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/directory"
  "github.com/markelog/eclectica/cmd/prompt"
)

func exists(path string) bool {
  _, err := os.Stat(path)
  return !os.IsNotExist(err)
}

func getVersions(name string) []string {
  path := variables.Home + "/" + name

  if !exists(path) {
    fmt.Println("There is no installed versions of " + name)
    os.Exit(0)
  }

  folders, err := ioutil.ReadDir(variables.Home + "/" + name)
  versions := []string{}

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  for _, folder := range folders {
    if folder.IsDir() {
      versions = append(versions, folder.Name())
    }
  }

  return versions
}

func listVersions(language string) {
  versions := getVersions(language)

  fmt.Println()
  for _, version := range versions {
    fmt.Println("  " + version)
  }
  fmt.Println()
}

func list() {
  language := prompt.List("Language", plugins.List).Language

  listVersions(language)
}

func askInfo() string {
  language := prompt.List("Language", plugins.List).Language
  version := prompt.List("Version", getVersions(language)).Version

  return language + "@" + version
}

func use() {
  activate(askInfo())
}

func main() {
  cli.AppHelpTemplate = `
Usage: e <name>, <name>@<version>

`
  if len(os.Args) == 1 {
    use()

  } else if os.Args[1] == "ls" {
    list()

  } else if os.Args[1] == "rm" {
    var nameAndVersion string

    if len(os.Args) == 2 {
      nameAndVersion = askInfo()
    } else {
      nameAndVersion = os.Args[2]
    }

    remove(nameAndVersion)

  } else {
    activate(os.Args[1])
  }
}

func remove(nameAndVersion string) {
  err := plugins.Remove(nameAndVersion)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func activate(language string) {
  info, err := plugins.Detect(language)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  path := fmt.Sprintf("%s/%s/%s", variables.Home, info["name"], info["version"])

  if exists(path) {
    err := plugins.Activate(info)

    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    os.Exit(0)
  }

  downloadPlace := download(info["url"])

  extractionPlace, err := directory.Create(info["name"])

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  err = archive.Extract(downloadPlace, extractionPlace)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  downloadPath := fmt.Sprintf("%s/%s", extractionPlace, info["filename"])
  extractionPath := fmt.Sprintf("%s/%s", extractionPlace, info["version"])
  err = os.Rename(downloadPath, extractionPath)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  err = plugins.Activate(info)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func download(url string) string {

  // Start file download
  fmt.Printf("Downloading %s...\n", url)

  respch, err := grab.GetAsync(os.TempDir(), url)

  if err != nil {
    fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", url, err)
    os.Exit(1)
  }

  // Block until HTTP/1.1 GET response is received
  fmt.Printf("Initializing download...\n")

  resp := <-respch

  // Print progress until transfer is complete
  for !resp.IsComplete() {
    fmt.Printf("\033[1AProgress %d / %d bytes (%d%%)\033[K\n", resp.BytesTransferred(), resp.Size, int(100*resp.Progress()))
    time.Sleep(200 * time.Millisecond)
  }

  // Clear progress line
  fmt.Printf("\033[1A\033[K")

  // Check for errors
  if resp.Error != nil {
    fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", url, resp.Error)
    os.Exit(1)
  }

  fmt.Printf("Downloaded\n")

  return resp.Filename
}
