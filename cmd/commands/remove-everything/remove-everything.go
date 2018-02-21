// Package version defines "version" command i.e. outputs current version
package removeEverything

import (
	"os"
	"path/filepath"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/shell"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/list"
	"github.com/spf13/cobra"
)

// Command config
var Command = &cobra.Command{
	Use:   "remove-everything",
	Short: "Removes everything related to the eclectica",
	Run:   run,
}

// Assume yes to the prompt
var assumeYes bool

// Init
func init() {
	flags := Command.PersistentFlags()

	flags.BoolVarP(&assumeYes, "assume-yes", "y", false, "Assume yes to the prompt")
}

func run(c *cobra.Command, args []string) {

	if assumeYes == false {
		response := list.GetWith("Are you sure?", []string{"yes", "no"})

		if response == "no" {
			return
		}
	}

	// Get ec binary
	path, err := os.Executable()
	print.Error(err)

	// Remove main executable
	err = os.Remove(path)
	print.Error(err)

	// Get ec-proxy binary
	ecProxy := filepath.Join(filepath.Dir(path), "ec-proxy")

	// Remove proxy executable
	err = os.Remove(ecProxy)
	print.Error(err)

	// Remove all the languages
	err = os.RemoveAll(variables.Base())
	print.Error(err)

	// Remove everything in rc files and restart the shell
	err = shell.New().Remove()
	print.Error(err)
}
