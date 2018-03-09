package list

import (
	"regexp"
	"strings"

	"github.com/markelog/list"
)

var (
	escapeReg = regexp.MustCompile(`[|\\{}()[\]^$+*?.]`)
)

func escape(str string) string {
	return escapeReg.ReplaceAllString(str, "")
}

func List(header string, options []string, indent int) string {
	strIndent := ""

	if indent > 0 {
		strIndent = strings.Repeat(" ", indent)
	}

	l := list.New(strIndent+header, options)

	if indent > 0 {
		l.SetIndent(indent + len(header)/2)
	}

	// Show the list
	l.Show()

	// Waiting for the user input
	return l.Get()
}
