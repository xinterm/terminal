package terminal

import (
	"github.com/xinterm/terminal/internal/sequence"
	"github.com/xinterm/terminal/util"
)

// Terminal describes a full featured virtual terminal
type Terminal struct {
	disp           display
	st             *state
	seqCtl         *sequence.SeqCtl
	RefreshDisplay func(*Grid)

	wrapLine bool

	minRefreshInterval int

	log *internalLog
}

// New creates a new terminal
func New(rows, cols int) *Terminal {
	t := &Terminal{
		log: &internalLog{},
	}

	t.st = newState(t.log)
	t.st.resize(rows, cols)

	t.seqCtl = sequence.NewSeqCtl(t.log)

	return t
}

// Run the terminal loop
func (t *Terminal) Run(updateDisplay func(*Grid)) {
	for {
		r := t.seqCtl.PollResult()
		t.st.sendResult(r)

		grid := t.disp.render(t.st)
		updateDisplay(grid)
	}
}

// Write implements io.Writer
func (t *Terminal) Write(p []byte) (n int, err error) {
	for _, c := range p {
		t.seqCtl.Parse(c)
	}
	return len(p), nil
}

// SetLogger accepts a Logger interface
func (t *Terminal) SetLogger(log util.Logger) {
	t.log.log = log
}

// SetHistoryLimit sets the scroll back size
func (t *Terminal) SetHistoryLimit(size int) {
	t.st.setScrollBackSize(size)
}

// SetMinRefreshInterval sets the minimal refresh interval in ms
func (t *Terminal) SetMinRefreshInterval(interval int) {
	t.minRefreshInterval = interval
}

// SetWrapLine sets the wrapline mode
func (t *Terminal) SetWrapLine(wrapLine bool) {
	t.wrapLine = wrapLine
	t.disp.SetWrapLine(t.wrapLine)
	//t.st.SetWrapLine(t.wrapLine)
}

// Resize the terminal
func (t *Terminal) Resize(rows, cols int) {
	t.st.resize(rows, cols)
}

// Search in terminal
func (t *Terminal) Search() int {
	return 0
}

// FlipScreen flips between the normal and alternative screen
func (t *Terminal) FlipScreen() {

}

// ScrollBackLineNumber gets the current scroll back line number
func (t *Terminal) ScrollBackLineNumber() int {
	return 0
}

// ScrollUp the screen
func (t *Terminal) ScrollUp(lineNumber int) int {
	current := 0
	return current
}

// ScrollDown the screen
func (t *Terminal) ScrollDown(lineNumber int) int {
	current := 0
	return current
}

// ScrollStart scroll to the start of the screen
func (t *Terminal) ScrollStart() {

}

// ScrollEnd scroll to the end of the screen
func (t *Terminal) ScrollEnd() {

}
