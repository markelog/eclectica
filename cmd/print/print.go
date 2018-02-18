// Package print have methods to print various eclectica info in style
package print

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
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

	// Gray color
	Gray = ansi.ColorCode("240")
	// White color
	White = ansi.ColorCode("white+b")
	// Reset ASCII-sequence
	Reset = ansi.ColorCode("reset")
	// Timeout for anything general
	Timeout = 200 * time.Millisecond
)

// InStyle prints header and text with style
func InStyle(name, entity string) {
	name = ansi.Color(name, "white+b")
	entity = ansi.Color(" "+entity+" ", "cyan+h")

	fmt.Print(name, entity, Reset)
}

// InStyleln prints header and text with style and newline char
func InStyleln(name, entity string) {
	InStyle(name, entity)
	fmt.Println()
}

// Version prints version in style
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

// CurrentVersion prints version with the heart symbol in style
func CurrentVersion(version string) {
	fmt.Println(ansi.Color("  ♥ "+version, "cyan"))
}

// Error prints error in style to stderr
func Error(err error) {
	if err == nil {
		return
	}

	red := ansi.Color("> ", "red")
	stderr := log.New(os.Stderr, "", 0)

	stderr.Println()

	stderr.Printf(red+"%v", err)

	if variables.IsDebug() {
		stderr.Println(errors.Wrap(err, 2).ErrorStack())
	}

	os.Exit(1)
}

// Download continuously prints download info
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
		Error(response.Error)

		cursed.MoveUp(1)
		cursed.EraseCurrentLine()
		InStyleln("Version", version)
	}

	spin := spinner.New(before, after, prefix, postfix)
	mutex := &sync.Mutex{}

	spin.Start()
	mutex.Lock()
	for response.IsComplete() == false {
		time.Sleep(time.Millisecond * 100)
	}
	mutex.Unlock()
	spin.Stop()

	return response.Filename
}

// Warning prints warning and how to fix it in style
func Warning(note, command string) {
	stderr := log.New(os.Stderr, "", 0)

	stderr.Println()
	stderr.Print(ansi.Color("> ", "red"), note)

	if command != "" {
		stderr.Println()
		stderr.Println()

		stderr.Print(ansi.Color("> ", "green") + command)
	}
}

func ClosestLangWarning(language, closest string) {
	incorrectOne := ansi.Color(language, "red")

	if closest != "" {
		getArgs := strings.Join(os.Args, " ")
		corrected := strings.Replace(getArgs, language, closest, 1)

		withColor := ansi.Color(corrected, "green")

		Warning(
			`Eclectica does not support "`+
				incorrectOne+`", perhaps you meant "`+withColor+`"`,
			"",
		)

		return
	}

	Warning(
		`Eclectica does not support "`+incorrectOne+`", whatever that is`,
		"",
	)
}
