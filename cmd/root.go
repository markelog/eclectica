package cmd

import (
	"fmt"
	"os"
  "strings"

	"github.com/spf13/cobra"

  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/cmd/activation"
  "github.com/markelog/eclectica/cmd/info"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:     "eclectica",
	Short:   "Version manager for any language",
	Long:    `Eclectica is version language manager for Node.js`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if len(os.Args) == 1 {
    use()
    return
  }

  language := strings.Split(os.Args[1], "@")[0]
  for _, plugin := range plugins.List {
    if strings.HasPrefix(language, plugin) {
      activation.Activate(os.Args[1])
      return
    }
  }

  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(-1)
  }
}

func init() {
	cobra.OnInitialize()
}

func use() {
  activation.Activate(info.Ask())
}
