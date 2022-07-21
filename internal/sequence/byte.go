package sequence

type byteSequence struct {
	baseSequence

	byteRange [][2]byte

	char byte
}

func (seq *byteSequence) addRange(low, hi byte) {
	seq.byteRange = append(seq.byteRange, [2]byte{low, hi})
}

func (seq *byteSequence) shouldEnter(c byte) bool {
	for _, r := range seq.byteRange {
		if c >= r[0] && c <= r[1] {
			return true
		}
	}
	return false
}

func (seq *byteSequence) shouldExit() bool {
	return true
}

func (seq *byteSequence) shouldHandOver() bool {
	return false
}

func (seq *byteSequence) reset() {
	seq.char = 0
}

func (seq *byteSequence) process(c byte) {
	seq.char = c
}
