package path

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/initiate"
	"github.com/markelog/eclectica/plugins"
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
	addition := initiate.Compose(plugins.Plugins)

	if strings.Contains(path, addition) {
		fmt.Print(path)
	} else {
		fmt.Print(path + ":" + addition)
	}

	os.Exit(0)
}
