package cmd

import (
	"fmt"
  "os"
  "errors"

	"github.com/spf13/cobra"

  "github.com/markelog/eclectica/cmd/info"
  "github.com/markelog/eclectica/plugins"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove language version",
  RunE: func(cmd *cobra.Command, args []string) error {
    _, version := info.GetLanguage(args)

    if version == "" {
      return errors.New("Can't remove without specific version")
    }

    return nil
  },
	Run: func(cmd *cobra.Command, args []string) {
    if len(args) == 0 {
      remove(info.Ask())
    } else {
      remove(info.GetLanguage(args))
    }
	},
}

func remove(language, version string) {
  err := plugins.Remove(language, version)

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
