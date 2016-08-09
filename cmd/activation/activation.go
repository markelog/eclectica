package activation

import (
  "fmt"
  "os"
  "time"
  "strings"
  "errors"

  "github.com/markelog/archive"
  "github.com/cavaliercoder/grab"
  "github.com/fatih/color"
  "github.com/tj/go-spin"
  "github.com/dustin/go-humanize"
  "github.com/sethgrid/curse"

  "github.com/markelog/eclectica/cmd/helpers"
  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/variables"
  "github.com/markelog/eclectica/directory"
)

func checkErrors(err error) {
  if err == nil {
    return
  }

  color.Set(color.FgRed)
  fmt.Print("> ")
  color.Unset()
  fmt.Fprintf(os.Stderr, "%v", err)
  fmt.Println()
  os.Exit(1)
}

func check404(resp *grab.Response, version string) {
  if resp.HTTPResponse.StatusCode == 404 {
    checkErrors(errors.New("Incorrect version " + version))
  }
}

func ActivateAndPrint(language string) {
  info, err := plugins.Detect(language)
  checkErrors(err)

  helpers.PrintInStyle("Language", info["name"])
  fmt.Println()
  helpers.PrintInStyle("Version", info["version"])
  fmt.Println()

  Activate(language)
}

func Activate(language string) {
  info, err := plugins.Detect(language)
  checkErrors(err)

  path := fmt.Sprintf("%s/%s/%s", variables.Home, info["name"], info["version"])

  if _, err := os.Stat(path); err == nil {
    err := plugins.Activate(info)

    checkErrors(err)

    os.Exit(0)
  }

  downloadPlace := download(info)
  extractionPlace, err := directory.Create(info["name"])
  checkErrors(err)

  // Just in case archive was downloaded, but not extracted files
  os.Remove(extractionPlace)

  err = archive.Extract(downloadPlace, extractionPlace)
  checkErrors(err)

  downloadPath := fmt.Sprintf("%s/%s", extractionPlace, info["filename"])
  extractionPath := fmt.Sprintf("%s/%s", extractionPlace, info["version"])

  err = os.Rename(downloadPath, extractionPath)
  checkErrors(err)

  err = plugins.Activate(info)
  checkErrors(err)
}

func download(info map[string]string) string {
  url := info["url"]

  // Start file download
  respch, err := grab.GetAsync(os.TempDir(), url)
  checkErrors(err)

  resp := <-respch

  check404(resp, info["version"])
  checkErrors(resp.Error)

  // Print progress until transfer is complete
  s := spin.New()
  c, _ := curse.New()
  started := false
  for !resp.IsComplete() {
    size := humanize.Bytes(resp.Size)
    transfered := humanize.Bytes(resp.BytesTransferred())
    transfered = strings.Replace(transfered, " MB", "", 1)

    c.MoveUp(1)
    if started {
      c.EraseCurrentLine()
    }
    started = true

    helpers.PrintInStyle("Version", info["version"])

    color.Set(color.FgBlack)
    fmt.Print("(")
    fmt.Printf("%s/%s ", transfered, size)

    color.Set(color.FgCyan)
    fmt.Print(s.Next())
    color.Unset()

    color.Set(color.FgBlack)
    fmt.Printf(" %d%%", int(100*resp.Progress()))
    fmt.Print(")")
    fmt.Println()

    time.Sleep(200 * time.Millisecond)
  }

  c.MoveUp(1)
  c.EraseCurrentLine()

  helpers.PrintInStyle("Version", info["version"])
  fmt.Println()

  checkErrors(resp.Error) // Don't know how to reproduce
  check404(resp, info["version"])

  return resp.Filename
}
