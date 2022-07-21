package terminal

import (
	"github.com/xinterm/terminal/internal/sequence"
	"github.com/xinterm/terminal/util"
)

type state struct {
	frm       *frame
	alterFrm  *frame
	alterMode bool

	title string

	scrollBackSize int
	rows           int
	cols           int

	log util.Logger
}

func newState(log util.Logger) *state {
	s := &state{
		frm:      newFrame(),
		alterFrm: newFrame(),
		log:      log,
	}
	return s
}

func (s *state) getFrame(flip bool) *frame {
	// XOR
	if s.alterMode != flip {
		return s.alterFrm
	}
	return s.frm
}

func (s *state) setScrollBackSize(size int) {
	s.scrollBackSize = size
	s.frm.scrl.setMaxLines(size + s.rows)
}

func (s *state) resize(rows, cols int) {
	if rows == s.rows && cols == s.cols {
		return
	}

	s.rows = rows
	s.cols = cols

	s.frm.resize(rows, cols)
	s.alterFrm.resize(rows, cols)

	s.frm.scrl.setMaxLines(s.scrollBackSize + rows)
	s.alterFrm.scrl.setMaxLines(rows)
}

func (s *state) sendResult(r *sequence.Result) {
	s.log.Debugf("Result: %#v", r)

	switch r.Type {
	case sequence.ResultASCII:
		s.getFrame(false).printChar(rune(r.Value.(byte)))
	case sequence.ResultUTF8:
		s.getFrame(false).printChar(r.Value.(rune))
	}
}

func (s *state) gridVisible(flip bool) *Grid {
	return s.getFrame(flip).visibleLines()
}

func (s *state) gridToEnd(endLine, displayLines int, flip bool) *Grid {
	return s.getFrame(flip).scrl.gridToEnd(endLine, displayLines)
}
