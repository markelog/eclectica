package cmd

import (
  "os"
	"fmt"

	"github.com/spf13/cobra"

  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/cmd/prompt"
  "github.com/markelog/eclectica/cmd/info"
)

func exists(path string) bool {
  _, err := os.Stat(path)
  return !os.IsNotExist(err)
}

func listVersions(language string) {
  versions := info.Versions(language)

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

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List installed language versions",
	Run: func(cmd *cobra.Command, args []string) {
    list()
  },
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
