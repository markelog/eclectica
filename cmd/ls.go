package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
  "github.com/fatih/color"

  "github.com/markelog/list"

  "github.com/markelog/eclectica/cmd/print"
  "github.com/markelog/eclectica/plugins"
)

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
  versions := plugins.Versions(language)
  current := plugins.CurrentVersion(language)

  listVersions(versions, current)
}

func listLocal() {
  language := list.GetWith("Language", plugins.List)

  listLocalVersions(language)
}

func listRemoteVersions(language string) {
  remoteList, _ := plugins.RemoteList(language)
  mask := list.GetWith("Mask", plugins.GetKeys(remoteList))
  versions := plugins.GetElements(mask, remoteList)
  current := plugins.CurrentVersion(language)

  listVersions(versions, current)
}

func listRemote() {
  language := list.GetWith("Language", plugins.List)

  listRemoteVersions(language)
}

func remote(args []string) {
  if len(args) == 0 {
    listRemote()
    return
  }

  for _, element := range plugins.List {
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

  for _, element := range plugins.List {
    if args[0] == element {
      print.InStyle("Language", element)
      fmt.Println()
      listLocalVersions(element)
      return
    }
  }
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
  Use:     "ls",
  Short:   "List installed language versions",
  Run: func(cmd *cobra.Command, args []string) {
    if isRemote {
      remote(args)
    } else {
      local(args)
    }
  },
}

func init() {
  RootCmd.AddCommand(lsCmd)
}
