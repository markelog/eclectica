package ls

import (
	"fmt"

	"github.com/markelog/list"
	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/cmd/flags"
	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins"
)

// Command represents the ls command
var Command = &cobra.Command{
	Use:     "ls",
	Short:   "List installed language versions",
	Example: example,
	Run:     run,
}

// Command example
var example = `
  List local language specific versions
  $ ec ls rust

  List remote language specific versions
  $ ec ls -r node

  List local versions
  $ ec ls

  List remote versions
  $ ec ls -r
`

// Runner
func run(cmd *cobra.Command, args []string) {
	if flags.IsRemote {
		remote(args)
	} else {
		local(args)
	}
}

// List versions
func listVersions(versions []string, current string) {
	fmt.Println()
	for _, version := range versions {
		if current == version {
			print.CurrentVersion(version)
		} else {
			print.Version(version)
		}
	}
	fmt.Println()
}

// List local ones
func listLocalVersions(language string) {
	plugin := plugins.New(language)

	versions, err := plugin.List()
	print.Error(err)

	current, err := io.GetVersion(plugin.Dots())
	print.Error(err)

	// In case we could find `.<language>-version` file i.e. there is no local version
	if current == "current" || current == "" {
		current = plugin.Current()
	}

	listVersions(versions, current)
}

// Ask for language and list local versions
func listLocal() {
	language := list.GetWith("Language", plugins.Plugins)

	listLocalVersions(language)
}

// Ask for remote versions and list them
func listRemoteVersions(language string) {
	versions, err := info.AskRemoteVersions(language)
	print.Error(err)

	current := plugins.New(language).Current()
	listVersions(versions, current)
}

// Ask for language and list remote versions
func listRemote() {
	language := list.GetWith("Language", plugins.Plugins)

	listRemoteVersions(language)
}

// Main entry point for remote output
func remote(args []string) {
	if len(args) == 0 {
		listRemote()
		return
	}

	for _, plugin := range plugins.Plugins {
		if args[0] == plugin {
			print.InStyle("Language", plugin)
			fmt.Println()
			listRemoteVersions(plugin)
			return
		}
	}
}

// Main entry point for local output
func local(args []string) {
	if len(args) == 0 {
		listLocal()
		return
	}

	for _, plugin := range plugins.Plugins {
		if args[0] == plugin {
			print.InStyle("Language", plugin)
			fmt.Println()
			listLocalVersions(plugin)
			return
		}
	}
}

// Init
func init() {
	Command.PersistentFlags().BoolVarP(flags.RemoteFlag())
}
