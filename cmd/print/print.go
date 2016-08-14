package print

import(
  "fmt"
  "os"
  "strings"
  "time"

  "github.com/fatih/color"
  "github.com/tj/go-spin"
  "github.com/dustin/go-humanize"
  "github.com/sethgrid/curse"
  "github.com/cavaliercoder/grab"
)

func InStyle(name, entity string) {
  color.Set(color.Bold)
  fmt.Print(name)

  color.Set(color.FgCyan)
  fmt.Print(" ")
  fmt.Print(entity + " ")
  color.Unset()
}

func LaguageOrVersion(language, version string) {
  if language != "" {
    InStyle("Language", language)
    fmt.Println()
  }

  if version != "" {
    InStyle("Version", version)
    fmt.Println()
  }
}

func Error(err error) {
  if err == nil {
    return
  }

  fmt.Println(err)
  color.Set(color.FgRed)
  fmt.Print("> ")
  color.Unset()
  fmt.Fprintf(os.Stderr, "%v", err)
  fmt.Println()
  os.Exit(1)
}

func Download(response *grab.Response, version string) string {
  Error(response.Error)

  s := spin.New()
  c, _ := curse.New()
  started := false

  // Print progress until transfer is complete
  for response.IsComplete() == false {
    Error(response.Error)
    size := humanize.Bytes(response.Size)
    transfered := humanize.Bytes(response.BytesTransferred())
    transfered = strings.Replace(transfered, " MB", "", 1)

    c.MoveUp(1)
    if started {
      c.EraseCurrentLine()
    }
    started = true

    InStyle("Version", version)

    color.Set(color.FgBlack)
    fmt.Print("(")
    fmt.Printf("%s/%s ", transfered, size)
    color.Unset()

    color.Set(color.FgCyan)
    fmt.Print(s.Next())
    color.Unset()

    color.Set(color.FgBlack)
    fmt.Printf(" %d%%", int(100*response.Progress()))
    fmt.Print(")")
    fmt.Println()
    color.Unset()

    time.Sleep(200 * time.Millisecond)
  }

  c.MoveUp(1)
  c.EraseCurrentLine()

  InStyle("Version", version)
  fmt.Println()

  return response.Filename
}

