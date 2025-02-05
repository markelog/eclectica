// Package version defines "version" command i.e. outputs current version
package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Command config
var Command = &cobra.Command{
	Use:   "version",
	Short: "print version of eclectica",
	Run:   run,
}

// Version number
const Version = "0.9.0"

// Runner
func run(c *cobra.Command, args []string) {
	fmt.Println(Version)
}
