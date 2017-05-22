package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Command config
var Command = &cobra.Command{
	Use:   "version",
	Short: "Print version of Eclectica",
	Run:   run,
}

// Version number
const Version = "0.0.13"

// Runner
func run(c *cobra.Command, args []string) {
	fmt.Println(Version)
}
