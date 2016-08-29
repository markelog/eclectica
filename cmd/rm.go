package cmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove language version",
	Run: func(cmd *cobra.Command, args []string) {
		if info.HasVersion(args) == false {
			print.Error(errors.New("Can't remove without specific version"))
		}

		var (
			language string
			version  string
			err      error
		)

		if len(args) == 0 {
			language, version, err = info.Ask()
		} else {
			language, version = info.GetLanguage(args)
		}

		remove(language, version, err)
	},
}

func remove(language, version string, err error) {
	print.Error(err)

	err = plugins.New(language).Remove(version)

	print.Error(err)

}

func init() {
	RootCmd.AddCommand(rmCmd)
}
