// Package commands defines CLI API
package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var use = "ec [<language>@<version>]"

// Command config
var Command = &cobra.Command{
	Use:     use,
	Example: example,
	Hidden:  true,
}

// Register command to root command
func Register(cmd *cobra.Command) {
	Command.AddCommand(cmd)
}

// Execute the command
func Execute() {
	// Until https://github.com/spf13/cobra/pull/369 is landed
	args := os.Args[1:]
	cmd, args, _ := Command.Find(args)
	name := cmd.Name()

	if name == "ec" && hasHelp(args) == false {
		augment()
	}

	Command.Execute()
}

func init() {
	Command.SetHelpTemplate(help)
	Command.SetUsageTemplate(usage)

	cobra.OnInitialize()
}

func augment() {
	// Insert command "install" in args
	os.Args = append(os.Args[:1], append([]string{"install"}, os.Args[1:]...)...)
}

func hasHelp(args []string) bool {
	for _, elem := range args {
		if elem == `--help` || elem == `-h` {
			return true
		}
	}

	return false
}
