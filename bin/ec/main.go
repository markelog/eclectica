package main

import (
	"github.com/markelog/eclectica/cmd/commands"

	// Commands
	"github.com/markelog/eclectica/cmd/commands/install"
	"github.com/markelog/eclectica/cmd/commands/ls"
	"github.com/markelog/eclectica/cmd/commands/path"
	removeEverything "github.com/markelog/eclectica/cmd/commands/remove-everything"
	"github.com/markelog/eclectica/cmd/commands/rm"
	"github.com/markelog/eclectica/cmd/commands/version"
)

func main() {
	commands.Register(install.Command)
	commands.Register(rm.Command)
	commands.Register(ls.Command)
	commands.Register(version.Command)
	commands.Register(path.Command)
	commands.Register(removeEverything.Command)

	commands.Execute()
}
