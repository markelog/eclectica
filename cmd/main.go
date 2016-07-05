package main

import (
  "os"
  "fmt"
  "time"

  "github.com/urfave/cli"
  "github.com/cavaliercoder/grab"
  "github.com/markelog/eclectica/detect"
)

func download(url string) {

  // start file download
  fmt.Printf("Downloading %s...\n", url)
  respch, err := grab.GetAsync(".", url)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", url, err)
    os.Exit(1)
  }

  // block until HTTP/1.1 GET response is received
  fmt.Printf("Initializing download...\n")
  resp := <-respch

  // print progress until transfer is complete
  for !resp.IsComplete() {
    fmt.Printf("\033[1AProgress %d / %d bytes (%d%%)\033[K\n", resp.BytesTransferred(), resp.Size, int(100*resp.Progress()))
    time.Sleep(200 * time.Millisecond)
  }

  // clear progress line
  fmt.Printf("\033[1A\033[K")

  // check for errors
  if resp.Error != nil {
    fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", url, resp.Error)
    os.Exit(1)
  }

  fmt.Printf("Successfully downloaded to ./%s\n", resp.Filename)
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

    download(dists["url"])
  }
}
