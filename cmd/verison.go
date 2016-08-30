package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "0.0.1"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of Eclectica",
	Run: func(c *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}
