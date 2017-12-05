package print

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-errors/errors"
	"github.com/markelog/curse"
	"github.com/mgutz/ansi"
	"gopkg.in/cavaliercoder/grab.v1"

	"github.com/markelog/eclectica/cmd/print/spinner"
	"github.com/markelog/eclectica/variables"
)

var (
	Gray    = ansi.ColorCode("240")
	White   = ansi.ColorCode("white+b")
	Reset   = ansi.ColorCode("reset")
	Timeout = 200 * time.Millisecond
)

func InStyle(name, entity string) {
	name = ansi.Color(name, "white+b")
	entity = ansi.Color(" "+entity+" ", "cyan+h")

	fmt.Print(name, entity)
}

func InStyleln(name, entity string) {
	InStyle(name, entity)
	fmt.Println()
}

func Version(args ...interface{}) {
	var color string
	version := args[0].(string)

	if len(args) == 1 {
		color = Gray
	} else {
		color = ansi.ColorCode(args[1].(string))
	}

	fmt.Println(color, "  ", version, Reset)
}

func CurrentVersion(version string) {
	fmt.Println(ansi.Color("  ♥ "+version, "cyan"))
}

func Error(err error) {
	if err == nil {
		return
	}

	stderr := log.New(os.Stderr, "", 0)

	stderr.Println()
	stderr.Print(ansi.Color("> ", "red"))

	stderr.Printf("%v", err)
	stderr.Println()

	if variables.IsDebug() {
		stderr.Println(errors.Wrap(err, 2).ErrorStack())
	}

	stderr.Println()
	stderr.Println()

	os.Exit(1)
}

func Download(response *grab.Response, version string) string {
	Error(response.Error)

	cursed, _ := curse.New()

	sizeAndTransfer := func() (size, transfer string) {
		size = humanize.Bytes(response.Size)

		transfer = humanize.Bytes(response.BytesTransferred())
		transfer = strings.Replace(transfer, " MB", "", 1)

		return
	}

	before := func() {}

	prefix := func() {
		Error(response.Error)

		cursed.MoveUp(1)
		cursed.EraseCurrentLine()

		size, transfer := sizeAndTransfer()
		text := fmt.Sprintf("(%s/%s ", transfer, size)

		InStyle("Version", version)
		fmt.Print(Gray, text, Reset)
	}

	postfix := func() {
		progress := int(100 * response.Progress())
		text := fmt.Sprintf("%d%%)", progress)

		fmt.Println(Gray, text, Reset)

		time.Sleep(Timeout)
	}

	after := func() {
		cursed.MoveUp(1)
		cursed.EraseCurrentLine()
		InStyleln("Version", version)
	}

	spin := &spinner.Spinner{
		Before:  before,
		After:   after,
		Prefix:  prefix,
		Postfix: postfix,
	}

	spin.Start()
	for response.IsComplete() == false {
		time.Sleep(time.Millisecond * 100)
	}
	spin.Stop()

	return response.Filename
}

func Warning(note, command string) {
	fmt.Println()
	fmt.Print(ansi.Color("> ", "red"))
	fmt.Print(note)

	if command != "" {
		fmt.Println()
		fmt.Println()

		fmt.Print(ansi.Color("> ", "green") + command)
	}

	fmt.Println()
	fmt.Println()
}
