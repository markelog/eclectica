package plugins

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/markelog/eclectica/variables"
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

func printShellMessage(name string) error {
	name = makeFirstUppercase(name)

	messageStart := name + ` has been installed, but it requires to restart your shell,
  please execute following command (you would need to do it only`
	messageMiddle := " once"
	messageEnd := "):"

	fmt.Println()

	color.Set(color.FgRed)
	fmt.Print("> ")
	color.Unset()

	color.Set(color.Bold)
	fmt.Print(messageStart)
	color.Set(color.FgRed)
	fmt.Print(messageMiddle)
	color.Unset()

	color.Set(color.Bold)
	fmt.Print(messageEnd)
	color.Unset()

	fmt.Println()
	fmt.Println()
	fmt.Println("                          $ " + variables.GetShellName())
	fmt.Println()

	return nil
}
