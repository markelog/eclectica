package cmd

import (
	"os"

	"github.com/spf13/cobra"
  "github.com/spf13/pflag"

  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/cmd/info"
  "github.com/markelog/eclectica/cmd/print"
)

var isRemote bool

var RootCmd = &cobra.Command{
	Use:     "eclectica",
	Short:   "Version manager for any language",
	Long: 	 "Cool and eclectic version manager for any language",
	Example: example,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
  args := os.Args[1:]

  if info.HasCommand(args) {
    // Initialize cobra for other commands
    if err := RootCmd.Execute(); err != nil {
      os.Exit(1)
    }

    return
  }

  // We don't use cobra here, since we support `ec <language>@version` syntax

  // If nothing was passed - just show list of the local versions
  if len(args) == 0 {
    install(info.Ask())
    return
  }

  pflag.BoolVarP(remoteInfo())
  pflag.Parse()

  language, version := info.GetLanguage(args)
  hasLanguage := info.HasLanguage(args)
  hasVersion := info.HasVersion(args)

  print.LaguageOrVersion(language, version)

  // In case of `ec <language>@<version>`
  if hasLanguage && hasVersion {
    install(language, version)
    return
  }

  // If `--remote` or `-r` flag was passed
  if isRemote {

    // In case of `ec -r`
    if hasVersion {
      install(info.AskRemote())
      return

    // In case of `ec -r <language>` or `ec <language> -r`
    } else {
      install(language, info.AskRemoteVersion(language))
      return
    }

    return
  }

  // In case of `ec <language>`
  if hasLanguage && hasVersion == false {
    install(language, info.AskVersion(language))
    return
  }

  // We already know it will show an error
  RootCmd.Execute()
  os.Exit(1)
}

func install(language, version string) {
  plugin := plugins.New(language, version)

  response, err := plugin.Download()
  print.Error(err)

  if response == nil {
    plugin.Install()
    return
  }

  print.Download(response, version)

  err = plugin.Activate()
  print.Error(err)
}

func remoteInfo() (*bool, string, string, bool, string) {
  return &isRemote, "remote", "r", false, "Get remote versions"
}

func init() {
	RootCmd.SetUsageTemplate(usage)
	cobra.OnInitialize()
}
