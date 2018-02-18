// Package commands defines CLI API
package commands

import (
	"os"

	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
	"github.com/schollz/closestmatch"
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

// Execute the command
func Execute() {
	cm := closestmatch.New(plugins.Plugins, []int{2})

	// Until https://github.com/spf13/cobra/pull/369 is landed
	args := os.Args[1:]
	cmd, args, _ := Command.Find(args)
	name := cmd.Name()

	if name == "ec" && hasHelp(args) == false {
		augment()
	}

	// Searching for closest plugin name
	if isLanguageRelated(name, args) && info.HasLanguage(args) == false {
		possible := info.PossibleLanguage(args)
		print.ClosestLangWarning(possible, cm.Closest(possible))
		return
	}

	Command.Execute()
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

func isLanguageRelated(name string, args []string) bool {
	if hasHelp(args) {
		return false
	}

	names := []string{"ec", "ls", "rm"}
	for _, elem := range names {
		if elem == name {
			return len(args) != 0
		}
	}

	return false
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
