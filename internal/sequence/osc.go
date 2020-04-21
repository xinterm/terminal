package sequence

import (
	"strings"
)

type oscSequence struct {
	baseSequence
	param strings.Builder
	exit  bool
}

func (seq *oscSequence) reset() {
	seq.param.Reset()
	seq.exit = false
}

func (seq *oscSequence) shouldEnter(c byte) bool {
	return c == ']'
}

func (seq *oscSequence) shouldExit() bool {
	return seq.exit
}
