// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
  "os"
	"fmt"

	"github.com/spf13/cobra"

  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/cmd/prompt"
  "github.com/markelog/eclectica/cmd/info"
)

func exists(path string) bool {
  _, err := os.Stat(path)
  return !os.IsNotExist(err)
}

func listVersions(language string) {
  versions := info.Versions(language)

  fmt.Println()
  for _, version := range versions {
    fmt.Println("  " + version)
  }
  fmt.Println()
}

func list() {
  language := prompt.List("Language", plugins.List).Language

  listVersions(language)
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   "List installed language versions",
	Run: func(cmd *cobra.Command, args []string) {
    list()
  },
}

func init() {
	RootCmd.AddCommand(lsCmd)
}
