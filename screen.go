package terminal

// Content for a single screen
type screen struct {
	lines []*Line

	sbBuf        *scrollBack
	sbLineNumber int

	cursorX int
	cursorY int

	rows int
	cols int

	wrapLine bool
}

func newScreen() *screen {
	return &screen{
		lines: []*Line{&Line{}},
		sbBuf: newScrollBack(),
	}
}

func (s *screen) setWrapLine(wrapLine bool) {
	s.wrapLine = wrapLine
}

func (s *screen) resize(rows, cols int) {
	s.rows = rows
	s.cols = cols
}

func (s *screen) printChar(r rune) {
	s.lines[s.cursorY].AddCell(&Cell{
		Char: r,
	})
	s.cursorX += charWidth(r)
}

func (s *screen) moveCursorX(x int) {

}

func (s *screen) moveCursorY(x int) {

}

func (s *screen) moveRelCursorX(x int) {

}

func (s *screen) moveRelCursorY(y int) {

}

func (s *screen) shiftTab() {

}
