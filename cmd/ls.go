package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/markelog/list"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List installed language versions",
	Run: func(cmd *cobra.Command, args []string) {
		if isRemote {
			remote(args)
		} else {
			local(args)
		}
	},
}

func listVersions(versions []string, current string) {
	fmt.Println()
	for _, version := range versions {
		if current == version {

			color.Set(color.FgCyan)
			fmt.Println("♥ " + version)
			color.Unset()

		} else {
			color.Set(color.FgBlack)
			fmt.Println("  " + version)
			color.Unset()
		}
	}
	fmt.Println()
}

func listLocalVersions(language string) {
	plugin := plugins.New(language)
	versions := plugin.List()
	current := plugin.Current()

	listVersions(versions, current)
}

func listLocal() {
	language := list.GetWith("Language", plugins.Plugins)

	listLocalVersions(language)
}

func listRemoteVersions(language string) {
	plugin := plugins.New(language)
	remoteList, _ := plugin.ListRemote()
	mask := list.GetWith("Mask", plugins.GetKeys(remoteList))
	versions := plugins.GetElements(mask, remoteList)
	current := plugin.Current()

	listVersions(versions, current)
}

func listRemote() {
	language := list.GetWith("Language", plugins.Plugins)

	listRemoteVersions(language)
}

func remote(args []string) {
	if len(args) == 0 {
		listRemote()
		return
	}

	for _, element := range plugins.Plugins {
		if args[0] == element {
			print.InStyle("Language", element)
			fmt.Println()
			listRemoteVersions(element)
			return
		}
	}
}

func local(args []string) {
	if len(args) == 0 {
		listLocal()
		return
	}

	for _, element := range plugins.Plugins {
		if args[0] == element {
			print.InStyle("Language", element)
			fmt.Println()
			listLocalVersions(element)
			return
		}
	}
}

func init() {
	RootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().BoolVarP(remoteInfo())
}
