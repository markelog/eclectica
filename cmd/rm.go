package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
)

var rmExample = `
  Install specifc version
  $ ec rm rust@1.11.0

  Remove language version with interactive list
  $ ec rm go

  Remove with interactive list
  $ ec rm
`

var rmCmd = &cobra.Command{
	Use:     "rm [<language>@<version>]",
	Short:   "Remove language version",
	Example: rmExample,
	Run:     rmRunner,
}

func rmRunner(cmd *cobra.Command, args []string) {
	var (
		language string
		version  string
		err      error
	)

	language, version = info.GetLanguage(args)
	hasLanguage := info.HasLanguage(args)
	hasVersion := info.HasVersion(args)

	if hasVersion == false {
		if hasLanguage {
			version, err = info.AskVersion(language)
		} else {
			language, version, err = info.Ask()
		}
	}

	if isCurrent(language, version) {
		err = errors.New("Cannot remove active version")
	}

	print.Error(err)

	remove(language, version, err)
}

func isCurrent(language, version string) bool {
	current := plugins.New(language).Current()

	return current == version
}

func remove(language, version string, err error) {
	err = plugins.New(language).Remove(version)
	print.Error(err)
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
