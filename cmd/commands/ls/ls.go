// Package ls defines "list" command i.e. outputs installed or remote available languages
package ls

import (
	"fmt"

	"github.com/go-errors/errors"
	"github.com/schollz/closestmatch"
	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/list"
	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/versions"
)

// Is action remote?
var isRemote bool

// Command represents the ls command
var Command = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "list installed language versions",
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
  $ ec ls -r`

// Runner
func run(cmd *cobra.Command, args []string) {
	var (
		cm = closestmatch.New(plugins.Plugins, []int{2})
	)

	// Searching for closest plugin name
	if len(args) > 0 && info.HasLanguage(args) == false {
		possible := info.PossibleLanguage(args)
		print.ClosestLangWarning(possible, cm.Closest(possible))
		return
	}

	if isRemote {
		remote(args)
	} else {
		local(args)
	}
}

// List versions
func listVersions(vers []string, current string) {
	completeCurrent, _ := versions.Complete(current, vers)

	fmt.Println()
	for i, version := range vers {

		if completeCurrent == version {
			print.CurrentVersion(version)
			continue
		}

		print.Version(version)

		if i == 9 {
			print.Version("...", "white")
			break
		}
	}

	print.LastPrint()
}

// List local ones
func listLocalVersions(language string) {
	plugin := plugins.New(&plugins.Args{
		Language: language,
	})

	versions := plugin.List()

	if len(versions) == 0 {
		err := errors.New("There is no installed versions")
		print.Error(err)
	}

	current, _, err := io.GetVersion(plugin.Dots())

	print.Error(err)

	// In case we could find `.<language>-version` file i.e. there is no local version
	if current == "current" || current == "" {
		current = plugin.Current()
	}

	listVersions(versions, current)
}

// Ask for language and list local versions
func listLocal() {
	fmt.Println()

	language := list.List("langauge:", plugins.Plugins, 0)

	listLocalVersions(language)
}

// Ask for remote versions and list them
func listRemoteVersions(language string) {
	versions, err := info.AskRemoteVersions(language)
	print.Error(err)

	current := plugins.New(&plugins.Args{
		Language: language,
	}).Current()

	listVersions(versions, current)
}

// Ask for language and list remote versions
func listRemote() {
	fmt.Println()

	language := list.List("langauge:", plugins.Plugins, 0)

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
			print.FnInStyleln("langauge:", plugin)
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
			print.FnInStyleln("langauge:", plugin)
			listLocalVersions(plugin)
			return
		}
	}
}

// Init
func init() {
	flags := Command.PersistentFlags()

	flags.BoolVarP(&isRemote, "remote", "r", false, "Get remote versions")
}
