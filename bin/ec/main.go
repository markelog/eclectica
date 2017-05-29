package main

import (
	"github.com/markelog/eclectica/cmd/commands/root"

	// Commands
	"github.com/markelog/eclectica/cmd/commands/ls"
	"github.com/markelog/eclectica/cmd/commands/path"
	"github.com/markelog/eclectica/cmd/commands/rm"
	"github.com/markelog/eclectica/cmd/commands/version"
)

func main() {
	root.Register(rm.Command)
	root.Register(ls.Command)
	root.Register(version.Command)
	root.Register(path.Command)

	root.Execute()
}
