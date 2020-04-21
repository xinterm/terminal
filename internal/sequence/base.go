package sequence

type baseSequence struct {
	handler func()
}

func (seq *baseSequence) reset() {
	// An empty base if nothing to reset
}

func (seq *baseSequence) process(byte) {
}

func (seq *baseSequence) startSubSequence() bool {
	return true
}

func (seq *baseSequence) dispatch() {
	if seq.handler != nil {
		seq.handler()
	}
}
