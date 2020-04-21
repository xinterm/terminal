package sequence

type escapeSequence struct {
	baseSequence
	char           byte
	processCounter int
	exit           bool
}

func (seq *escapeSequence) shouldEnter(c byte) bool {
	return c == 0x1b
}

func (seq *escapeSequence) shouldExit() bool {
	return seq.exit || seq.processCounter >= 3
}

func (seq *escapeSequence) reset() {
	seq.exit = false
	seq.char = 0
	seq.processCounter = 0
}

func (seq *escapeSequence) process(c byte) {
	switch seq.processCounter {
	case 0:
		break
	case 1:
		seq.char = c
	}
	seq.processCounter++
}

func (seq *escapeSequence) startSubSequence() bool {
	return seq.processCounter == 2
}
