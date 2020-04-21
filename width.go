package terminal

import (
	"unicode"

	"golang.org/x/text/width"
)

func charWidth(r rune) int {
	if unicode.IsControl(r) {
		return 0
	}

	p := width.LookupRune(r)
	switch p.Kind() {
	case width.Neutral:
		return 1
	case width.EastAsianAmbiguous:
		return 1
	case width.EastAsianNarrow, width.EastAsianHalfwidth:
		return 1
	case width.EastAsianWide, width.EastAsianFullwidth:
		return 2
	}

	return 0
}
