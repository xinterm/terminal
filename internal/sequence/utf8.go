package sequence

import (
	"unicode/utf8"
)

type utf8Sequence struct {
	baseSequence
	buf  []byte
	char rune
	exit bool
}

func (seq *utf8Sequence) shouldEnter(c byte) bool {
	return c >= 0x80
}

func (seq *utf8Sequence) shouldExit() bool {
	return seq.exit
}

func (seq *utf8Sequence) reset() {
	seq.buf = seq.buf[:0]
	seq.char = utf8.RuneError
	seq.exit = false
}

func (seq *utf8Sequence) process(c byte) {
	seq.buf = append(seq.buf, c)

	r, _ := utf8.DecodeRune(seq.buf)
	if r != utf8.RuneError {
		seq.char = r
		seq.exit = true
		return
	}

	if len(seq.buf) >= utf8.UTFMax {
		seq.exit = true
	}
}
