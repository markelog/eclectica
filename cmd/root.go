package cmd

import (
	"fmt"
	"os"
  "strings"

	"github.com/spf13/cobra"
  "github.com/spf13/pflag"

  "github.com/markelog/eclectica/cmd/activation"
  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/cmd/info"
  "github.com/markelog/eclectica/cmd/helpers"
)

var isRemote bool

var RootCmd = &cobra.Command{
	Use:     "eclectica",
	Short:   "Version manager for any language",
	Long:    `Eclectica is version language manager for Node.js`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
  // We don't use cobra here, since we support `ec <language>@version` syntax
  pflag.BoolVarP(&isRemote, "remote", "r", false, "Get remote versions")
  pflag.Parse()

  // If `--remote` flag were passed - just show list for the remote versions
  if isRemote {
    activation.Activate(info.AskRemote())
    return
  }

  // If nothing were passed - just show list for the local versions
  if len(os.Args) == 1 {
    activation.Activate(info.Ask())
    return
  }

  data := strings.Split(os.Args[1], "@")

  // In case of `ec <language>`
  if len(data) == 1 {
    activateWithoutVersion(data[0])
    return
  }

  // In case of `ec <language>@<version>`
  if len(data) == 2 {
    activateWithVersion(data[0])
    return
  }

  // Initialize cobra for other commands
  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func activateWithoutVersion(language string) {
  var version string

  for _, plugin := range plugins.List {
    if strings.HasPrefix(language, plugin) {
      helpers.PrintInStyle("Language", language)
      fmt.Println()

      version = info.AskVersion(language)
      activation.Activate(language + "@" + version)
      return
    }
  }
}

func activateWithVersion(language string) {
  for _, plugin := range plugins.List {
    if strings.HasPrefix(language, plugin) {
      activation.ActivateAndPrint(os.Args[1])
      return
    }
  }
}

func init() {
	cobra.OnInitialize()
  RootCmd.PersistentFlags().BoolVarP(&isRemote, "remote", "r", false, "Get remote versions")
}
