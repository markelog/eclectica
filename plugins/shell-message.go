package plugins

import (
	"bytes"
	"os"
	"strings"

	"github.com/markelog/eclectica/cmd/print"
)

func GetShellName() string {
	name := os.Getenv("SHELL")
	path := strings.Split(name, "/")

	return path[len(path)-1]
}

func makeFirstUppercase(name string) string {
	if len(name) < 2 {
		return strings.ToUpper(name)
	}

	bts := []byte(name)

	lc := bytes.ToUpper([]byte{bts[0]})
	rest := bts[1:]

	return string(bytes.Join([][]byte{lc, rest}, nil))
}

func printShellMessage(name string) {
	name = makeFirstUppercase(name)

	start := name + ` has been installed, but it requires to restart your shell,
  for this to take affect you need to execute following command (you would need to do it only`
	middle := " once"
	end := "):"
	command := "exec " + os.Getenv("SHELL")

	print.PostInstall(start, middle, end, command)
}
