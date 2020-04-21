package terminal

import (
	"github.com/xinterm/terminal/internal/sequence"
	"github.com/xinterm/terminal/util"
)

type state struct {
	scr       *screen
	alterScr  *screen
	alterMode bool

	title string

	log util.Logger
}

func newState(log util.Logger) *state {
	s := &state{
		scr:      newScreen(),
		alterScr: newScreen(),
		log:      log,
	}
	return s
}

func (s *state) currentScreen() *screen {
	if s.alterMode {
		return s.alterScr
	}
	return s.scr
}

func (s *state) setScrollBackSize(size int) {
	s.scr.sbBuf.setMaxLines(size)
}

func (s *state) resize(rows, cols int) {
	s.scr.resize(rows, cols)
	s.scr.resize(rows, cols)
}

func (s *state) sendResult(r *sequence.Result) {
}

func (s *state) printChar(r rune) {
	s.currentScreen().printChar(r)
}
