package cmd

import (
	"fmt"
  "os"

	"github.com/spf13/cobra"

  "github.com/markelog/eclectica/cmd/info"
  "github.com/markelog/eclectica/plugins"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove language version",
	Run: func(cmd *cobra.Command, args []string) {
		var nameAndVersion string

    if len(args) == 0 {
      nameAndVersion = info.Ask()
    } else {
      nameAndVersion = args[0]
    }

    remove(nameAndVersion)
	},
}

func remove(nameAndVersion string) {
  err := plugins.Remove(nameAndVersion)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
