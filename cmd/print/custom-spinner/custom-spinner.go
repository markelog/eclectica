// Package CustomSpinner provides spinner with additional custom logic
package CustomSpinner

import (
	"fmt"
	"sync"
	"time"

	"github.com/markelog/curse"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/cmd/print/spinner"
)

var (
	white   = print.White
	gray    = print.Gray
	reset   = print.Reset
	timeout = print.Timeout
)

// Spin essential struct
type Spin struct {
	SpinArgs
	Spinner *spinner.Spinner

	mutex *sync.Mutex
}

// SpinArgs is arguments struct for New() method
type SpinArgs struct {
	Header, Item, Note, Message string
}

// New returns custom spinner struct
func New(args *SpinArgs) *Spin {
	me := &Spin{
		mutex: &sync.Mutex{},
	}
	me.Set(args)

	return me
}

// Set new data for the spinner
func (me *Spin) Set(args *SpinArgs) {
	me.mutex.Lock()

	if len(args.Header) != 0 {
		me.Header = args.Header
	}

	if len(args.Item) != 0 {
		me.Item = args.Item
	}

	if len(args.Note) != 0 {
		me.Note = args.Note
	}

	if len(args.Message) != 0 {
		me.Message = args.Message
	}

	me.mutex.Unlock()

	if me.Spinner == nil {
		me.constructSpinner()
	}
}

// Start the spinner
func (me Spin) Start() {
	me.Spinner.Start()
}

// Stop the spinner
func (me Spin) Stop() {
	me.Spinner.Stop()
}

func (me *Spin) constructSpinner() {
	cursed, _ := curse.New()

	before := func() {}

	started := false
	prefix := func() {
		me.mutex.Lock()
		defer me.mutex.Unlock()

		cursed.MoveUp(1)

		if started {
			cursed.EraseCurrentLine()
		}
		started = true

		print.InStyle(me.Header, me.Item)
	}

	postfix := func() {
		me.mutex.Lock()
		defer me.mutex.Unlock()

		var (
			message string
		)

		if len(me.Message) > 0 {
			message = white + "(" + gray + me.Message + white + ")"
		}
		fmt.Println(gray, me.Note, message, reset)

		time.Sleep(timeout)
	}

	after := func() {
		me.mutex.Lock()
		defer me.mutex.Unlock()

		cursed.MoveUp(1)
		cursed.EraseCurrentLine()
		print.InStyleln(me.Header, me.Item)
	}

	me.Spinner = spinner.New(before, after, prefix, postfix)
}
