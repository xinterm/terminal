package sequence

type rootSequence struct {
	baseSequence
}

func (seq *rootSequence) shouldEnter(c byte) bool {
	return true
}

func (seq *rootSequence) shouldExit() bool {
	return false
}
