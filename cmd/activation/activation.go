package activation

import (
  "fmt"
  "os"
  "time"

  "github.com/markelog/archive"
  "github.com/cavaliercoder/grab"

  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/variables"
  "github.com/markelog/eclectica/directory"
)

func Activate(language string) {
  info, err := plugins.Detect(language)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  path := fmt.Sprintf("%s/%s/%s", variables.Home, info["name"], info["version"])

  if _, err := os.Stat(path); err == nil {
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
