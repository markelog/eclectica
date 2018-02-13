// Package commands defines CLI API
package commands

import (
	"os"
	"strings"

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
	flags.BoolVarP(&isRemote, "remote", "r", false, "ask for remote versions")
	flags.BoolVarP(&isLocal, "local", "l", false, "install to the current folder only")
	flags.BoolVarP(&withModules, "with-modules", "w", false, "reinstall global modules from the previous version (currently works only for node.js)")
}

func augment() {
	// Insert command "install" in args
	os.Args = append(os.Args[:1], append([]string{"install"}, os.Args[1:]...)...)
}

// Execute the command
func Execute() {

	// Until https://github.com/spf13/cobra/pull/369 is landed
	args := os.Args[1:]
	_, _, err := Command.Find(args)
	if err != nil && strings.HasPrefix(err.Error(), "unknown command") {
		augment()
	}

	Command.Execute()
}
