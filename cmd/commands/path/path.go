// Package path defines "path" command i.e. outputs proper $PATH
package path

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/shell"
)

// Command config
var Command = &cobra.Command{
	Use:    "path",
	Short:  "Echo path environment variable with eclectica specific keys",
	Run:    run,
	Hidden: true,
}

// Updates the path environment variable
func run(c *cobra.Command, args []string) {
	path := os.Getenv("PATH")
	addition := shell.Compose(plugins.Plugins)

	if strings.Contains(path, addition) {
		fmt.Print(path)
	} else {
		fmt.Print(addition + ":" + path)
	}

	os.Exit(0)
}
