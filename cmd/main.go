package main

import (
  "os"
  "fmt"
  "time"

  "github.com/markelog/eclectica/detect"
  "github.com/markelog/archive"

  "github.com/urfave/cli"
  "github.com/cavaliercoder/grab"
)

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


func main() {
  cli.AppHelpTemplate = `
Usage: e <name>, <name>@<version>

`
  if len(os.Args) == 1 {
    cli.NewApp().Run(os.Args)
  } else {
    dists, err := detect.Detect("node")

    if err != nil {
      fmt.Println(err)
    }

    path := download(dists["url"])

    archive.Extract(path, os.TempDir())
  }
}
