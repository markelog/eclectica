// Package strings provides simplified signatures for strings
package strings

import "golang.org/x/crypto/ssh/terminal"

// Elipsis returns trimmed string by provided maximum limit
func Elipsis(str string, max int) string {
	var (
		result = str
		length = len(str)
	)

	if max < 1 {
		return result
	}

	if length > max {
		if max > 3 {
			max -= 3
		}
		result = str[0:max] + "..."
	}
	return result
}

// ElipsisForTerminal returns trimmed string which should fit nicely to
// terminal pipe output
func ElipsisForTerminal(str string) (result string) {
	var (
		width, _, _ = terminal.GetSize(0)

		// Needed to compensate for the reset of the data already
		// printed to stdout/stderr
		subjectiveValue = 50
	)

	return Elipsis(str, width-subjectiveValue)
}
