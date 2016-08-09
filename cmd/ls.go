package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
  "github.com/fatih/color"

  "github.com/markelog/eclectica/cmd/helpers"
  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/cmd/prompt"
  "github.com/markelog/eclectica/cmd/info"
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
  versions := info.Versions(language)
  current := plugins.CurrentVersion(language)

  listVersions(versions, current)
}

func listLocal() {
  language := prompt.List("Language", plugins.List).Language

  listLocalVersions(language)
}

func listRemoteVersions(language string) {
  list, _ := plugins.RemoteList(language)
  key := prompt.List("Mask", plugins.GetKeys(list)).Mask
  versions := plugins.GetElements(key, list)
  current := plugins.CurrentVersion(language)

  listVersions(versions, current)
}

func listRemote() {
  language := prompt.List("Language", plugins.List).Language

  listRemoteVersions(language)
}

func remote(args []string) {
  if len(args) == 0 {
    listRemote()
    return
  }

  for _, element := range plugins.List {
    if args[0] == element {
      helpers.PrintInStyle("Language", element)
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
      helpers.PrintInStyle("Language", element)
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
