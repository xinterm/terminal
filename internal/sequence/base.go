package sequence

type baseSequence struct {
	handler func()
}

func (seq *baseSequence) shouldHandOver() bool {
	return true
}

func (seq *baseSequence) reset() {
	// An empty base if nothing to reset
}

func (seq *baseSequence) process(byte) {
}

func (seq *baseSequence) handle() {
	if seq.handler != nil {
		seq.handler()
	}
}
