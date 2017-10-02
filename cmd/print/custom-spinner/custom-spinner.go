package CustomSpinner

import (
	"fmt"
	"time"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/cmd/print/spinner"

	"github.com/markelog/curse"
)

var (
	white   = print.White
	gray    = print.Gray
	reset   = print.Reset
	timeout = print.Timeout
)

type SpinArgs struct {
	Header, Item, Note, Message string
}

type Spin struct {
	SpinArgs
	Spinner *spinner.Spinner
}

func New(args *SpinArgs) *Spin {
	me := &Spin{}
	me.Set(args)

	return me
}

func (me *Spin) Set(args *SpinArgs) {
	if me.Spinner == nil {
		me.constructSpinner()
	}

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
}

func (me Spin) Start() {
	me.Spinner.Start()
}

func (me Spin) Stop() {
	me.Spinner.Stop()
}

func (me *Spin) constructSpinner() {
	cursed, _ := curse.New()

	before := func() {}

	started := false
	prefix := func() {
		cursed.MoveUp(1)

		if started {
			cursed.EraseCurrentLine()
		}
		started = true

		print.InStyle(me.Header, me.Item)
	}

	postfix := func() {
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
		cursed.MoveUp(1)
		cursed.EraseCurrentLine()
		print.InStyleln(me.Header, me.Item)
	}

	me.Spinner = &spinner.Spinner{
		Before:  before,
		After:   after,
		Prefix:  prefix,
		Postfix: postfix,
	}
}
