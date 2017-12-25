// Package commands defines CLI API
package commands

import (
	"os"

	"github.com/spf13/cobra"
)

// Is action remote?
var isRemote bool

// Is action local?
var isLocal bool

// Reinstall global modules from previous version?
var withModules bool

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

func init() {
	Command.SetHelpTemplate(help)
	Command.SetUsageTemplate(usage)

	cobra.OnInitialize()

	flags := Command.PersistentFlags()
	flags.BoolVarP(&isRemote, "remote", "r", false, "Get remote versions")
	flags.BoolVarP(&isLocal, "local", "l", false, "Install as local version i.e. language will be installed only to local folder")
	flags.BoolVarP(&withModules, "with-modules", "w", false, "Reinstall global modules from the previous version (currently works only for node.js)")
}

func augment() {
	// Insert command "install" in args
	os.Args = append(os.Args[:1], append([]string{"install"}, os.Args[1:]...)...)
}

// Execute the command
func Execute() {

	args := os.Args[1:]

	// Until https://github.com/spf13/cobra/pull/369 is landed
	// Workaround to "forward" to a know command when no know command found
	cmd, _, _ := Command.Find(args)

	if cmd.Use == use {
		augment()
	}

	Command.Execute()
}
