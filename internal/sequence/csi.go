package sequence

import (
	"strings"
)

type csiSequence struct {
	baseSequence
	param          strings.Builder
	intermediate   strings.Builder
	final          byte
	processCounter int
}

func (seq *csiSequence) reset() {
	seq.param.Reset()
	seq.intermediate.Reset()
	seq.final = 0
	seq.processCounter = 0
}

func (seq *csiSequence) shouldEnter(c byte) bool {
	return c == '['
}

func (seq *csiSequence) shouldExit() bool {
	return seq.final != 0
}

func (seq *csiSequence) process(c byte) {
	seq.processCounter++
}

func (seq *csiSequence) shouldHandOver() bool {
	return seq.processCounter > 1
}
