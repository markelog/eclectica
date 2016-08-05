package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
  "github.com/fatih/color"

  "github.com/markelog/eclectica/cmd/helpers"
  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/cmd/prompt"
  "github.com/markelog/eclectica/cmd/info"
)

func listVersions(language string) {
  versions := info.Versions(language)
  current := plugins.CurrentVersion(language)

  fmt.Println()
  for _, version := range versions {
    if current == version {

      color.Set(color.FgCyan)
      fmt.Println("♥ " + version)
      color.Unset()

    } else {
      color.Set(color.FgBlack)
      fmt.Println("  " + version)
      color.Unset()
    }
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

    if len(args) == 0 {
      list()
      return
    }

    for _, element := range plugins.List {
      if args[0] == element {
        helpers.PrintInStyle("Language", element)
        fmt.Println()
        listVersions(element)
        return
      }
    }
  },
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
