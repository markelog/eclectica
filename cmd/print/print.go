package print

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-errors/errors"

	"github.com/cavaliercoder/grab"
	"github.com/dustin/go-humanize"
	"github.com/markelog/curse"
	"github.com/markelog/eclectica/variables"
	"github.com/mgutz/ansi"
)

var (
	gray  = ansi.ColorCode("240")
	reset = ansi.ColorCode("reset")
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
		color = gray
	} else {
		color = ansi.ColorCode(args[1].(string))
	}

	fmt.Println(color, "  ", version, reset)
}

func CurrentVersion(version string) {
	fmt.Println(ansi.Color("  ♥ "+version, "cyan"))
}

func Error(err error) {
	if err == nil {
		return
	}

	fmt.Println()
	fmt.Print(ansi.Color("> ", "red"))

	fmt.Fprintf(os.Stderr, "%v", err)

	if variables.IsDebug() {
		fmt.Println(errors.Wrap(err, 2).ErrorStack())
	}

	fmt.Println()
	fmt.Println()

	os.Exit(1)
}

func Download(response *grab.Response, version string) string {
	Error(response.Error)

	c, _ := curse.New()

	before := func() {
		time.Sleep(500 * time.Millisecond)
	}

	started := false
	prefix := func() {
		Error(response.Error)
		size := humanize.Bytes(response.Size)
		transfered := humanize.Bytes(response.BytesTransferred())
		transfered = strings.Replace(transfered, " MB", "", 1)

		c.MoveUp(1)

		if started {
			c.EraseCurrentLine()
		}
		started = true
		text := fmt.Sprintf("(%s/%s ", transfered, size)

		InStyle("Version", version)
		fmt.Print(gray, text, reset)
	}

	postfix := func() {
		progress := int(100 * response.Progress())
		text := fmt.Sprintf("%d%%)", progress)

		fmt.Println(gray, text, reset)

		time.Sleep(200 * time.Millisecond)
	}

	after := func() {
		c.EraseCurrentLine()
		InStyle("Version", version)
		fmt.Println()
	}

	s := &Spinner{
		Before:  before,
		After:   after,
		Prefix:  prefix,
		Postfix: postfix,
	}

	s.Start()
	for response.IsComplete() == false {
		time.Sleep(200 * time.Millisecond)
	}
	s.Stop()

	return response.Filename
}

func CustomSpin(header, item, message string) *Spinner {
	c, _ := curse.New()

	before := func() {}

	started := false
	prefix := func() {
		c.MoveUp(1)

		if started {
			c.EraseCurrentLine()
		}
		started = true

		InStyle(header, item)
	}

	postfix := func() {
		fmt.Println(gray, message, reset)

		time.Sleep(300 * time.Millisecond)
	}

	after := func() {
		c.EraseCurrentLine()
		InStyle(header, item)
		fmt.Println()
	}

	s := &Spinner{
		Before:  before,
		After:   after,
		Prefix:  prefix,
		Postfix: postfix,
	}

	return s
}

func Warning(message, command string) {
	fmt.Println()
	fmt.Print(ansi.Color("> ", "red"))
	fmt.Print(message)

	if command != "" {
		fmt.Println()
		fmt.Println()

		fmt.Print(ansi.Color("> ", "green") + command)
	}

	fmt.Println()
	fmt.Println()
}
