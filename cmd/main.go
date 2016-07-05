package main

import (
  "os"
  "fmt"

  "github.com/urfave/cli"

  "github.com/markelog/eclectica/detect"
)


func main() {
  cli.AppHelpTemplate = `
Usage: e <name>, <name>@<version>

`
  if len(os.Args) == 2 {
    cli.NewApp().Run(os.Args)
  } else {
    err := detect.Detect("node")

    if err != nil {
      fmt.Println(err)
    }
  }
}
