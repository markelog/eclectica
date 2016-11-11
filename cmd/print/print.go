package print

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/markelog/curse"
)

func InStyle(name, entity string) {
	color.Set(color.Bold)
	fmt.Print(name)

	color.Set(color.FgCyan)
	fmt.Print(" ")
	fmt.Print(entity + " ")
	color.Unset()
}

func InStyleln(name, entity string) {
	InStyle(name, entity)
	fmt.Println()
}

func LaguageOrVersion(language, version string) {
	if language != "" {
		InStyle("Language", language)
		fmt.Println()
	}

	if version != "" {
		InStyle("Version", version)
		fmt.Println()
	}
}

func Error(err error) {
	if err == nil {
		return
	}

	fmt.Println()
	color.Set(color.FgRed)
	fmt.Print("> ")
	color.Unset()
	fmt.Fprintf(os.Stderr, "%v", err)
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

		InStyle("Version", version)

		color.Set(color.FgBlack)
		fmt.Print("(")
		fmt.Printf("%s/%s ", transfered, size)
		color.Unset()
	}

	postfix := func() {
		color.Set(color.FgBlack)
		fmt.Printf(" %d%%", int(100*response.Progress()))
		fmt.Print(")")
		fmt.Println()
		color.Unset()

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
		color.Set(color.FgBlack)
		fmt.Print(" ", message)
		color.Unset()
		fmt.Println()

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

func PostInstall(start, middle, end, command string) {
	fmt.Println()

	color.Set(color.FgRed)
	fmt.Print("> ")
	color.Unset()

	color.Set(color.Bold)
	fmt.Print(start)
	color.Set(color.FgRed)
	fmt.Print(middle)
	color.Unset()

	color.Set(color.Bold)
	fmt.Print(end)
	color.Unset()

	fmt.Println()
	fmt.Println()

	color.Set(color.FgGreen)
	fmt.Print("> ")
	color.Unset()

	fmt.Print(command)
	fmt.Println()
}
