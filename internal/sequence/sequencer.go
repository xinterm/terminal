package sequence

type sequencer interface {
	shouldEnter(byte) bool
	shouldExit() bool
	shouldHandOver() bool

	reset()
	process(byte)
	handle()
}
