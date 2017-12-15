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
var WithModules bool

// Command config
var Command = &cobra.Command{
	Use:     "ec [<language>@<version>]",
	Example: example,
	Hidden:  true,
}

// Add command to root command
func Register(cmd *cobra.Command) {
	Command.AddCommand(cmd)
}

// Init
func init() {
	Command.SetHelpTemplate(help)
	Command.SetUsageTemplate(usage)

	cobra.OnInitialize()

	flags := Command.PersistentFlags()
	flags.BoolVarP(&isRemote, "remote", "r", false, "Get remote versions")
	flags.BoolVarP(&isLocal, "local", "l", false, "Install as local version i.e. language will be installed only to local folder")
	flags.BoolVarP(&WithModules, "with-modules", "w", false, "Reinstall global modules from the previous version (currently works only for node.js)")
}

func Execute() {

	// Until https://github.com/spf13/cobra/pull/369 is landed
	// Workaround to "forward" to a know command when no know command found
	_, _, err := Command.Find(os.Args[1:])

	if err != nil && strings.HasPrefix(err.Error(), "unknown command") {

		// Insert command "install" in args
		os.Args = append(os.Args[:1], append([]string{"install"}, os.Args[1:]...)...)
	}

	Command.Execute()
}
