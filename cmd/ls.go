package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/markelog/curse"
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

	versions, err := plugin.List()
	print.Error(err)

	current := plugin.Current()
	listVersions(versions, current)
}

func listLocal() {
	language := list.GetWith("Language", plugins.Plugins)

	listLocalVersions(language)
}

func listRemoteVersions(language string) {
	plugin := plugins.New(language)
	c, _ := curse.New()

	prefix := func() {
		c.MoveUp(1)
		c.EraseCurrentLine()
		print.InStyle("Language", language)
	}

	postfix := func() {
		fmt.Println()
		time.Sleep(200 * time.Millisecond)
	}

	s := &print.Spinner{
		Before:  func() { time.Sleep(500 * time.Millisecond) },
		After:   func() { fmt.Println() },
		Prefix:  prefix,
		Postfix: postfix,
	}

	s.Start()
	remoteList, err := plugin.ListRemote()
	s.Stop()

	print.Error(err)

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

	for _, plugin := range plugins.Plugins {
		if args[0] == plugin {
			print.InStyle("Language", plugin)
			fmt.Println()
			listRemoteVersions(plugin)
			return
		}
	}
}

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

func init() {
	RootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().BoolVarP(remoteInfo())
}
