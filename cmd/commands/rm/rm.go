// Package rm defines "remove" command
package rm

import (
	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
)

// Command config
var Command = &cobra.Command{
	Use:     "rm [<language>@<version>]",
	Aliases: []string{"remove"},
	Short:   "remove language version",
	Example: example,
	Run:     run,
}

// Command example
var example = `
  Install specifc version
  $ ec rm rust@1.11.0

  Remove language version with interactive list
  $ ec rm go

  Remove with interactive list
  $ ec rm
`

// Runner
func run(cmd *cobra.Command, args []string) {
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
			print.InStyleln("Language", language)
			version, err = info.AskVersion(language)
		} else {
			language, version, err = info.Ask()
		}
	}

	print.Error(err)

	remove(language, version)
}

// Try to remove
func remove(language, version string) {
	err := plugins.New(&plugins.Args{
		Language: language,
		Version:  version,
	}).Remove()
	print.Error(err)
}
