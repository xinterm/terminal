package sequence

import (
	"strings"
)

type escapeSequence struct {
	baseSequence
	handOver bool
	exit     bool
}

func (seq *escapeSequence) shouldEnter(c byte) bool {
	return c == 0x1b
}

func (seq *escapeSequence) shouldExit() bool {
	return seq.exit
}

func (seq *escapeSequence) shouldHandOver() bool {
	return seq.handOver
}

func (seq *escapeSequence) reset() {
	seq.handOver = false
	seq.exit = false
}

func (seq *escapeSequence) process(c byte) {
	seq.handOver = true
}

type escapeSTSequence struct {
	baseSequence

	byteRange [][2]byte

	char     byte
	param    strings.Builder
	handOver bool
	exit     bool
}

func (seq *escapeSTSequence) shouldEnter(c byte) bool {
	for _, r := range seq.byteRange {
		if c >= r[0] && c <= r[1] {
			return true
		}
	}
	return false
}

func (seq *escapeSTSequence) shouldExit() bool {
	return seq.exit
}

func (seq *escapeSTSequence) shouldHandOver() bool {
	return seq.handOver
}

func (seq *escapeSTSequence) reset() {
	seq.char = 0
	seq.param.Reset()
	seq.handOver = false
	seq.exit = false
}

func (seq *escapeSTSequence) process(c byte) {
	seq.char = c
	seq.handOver = true
}
