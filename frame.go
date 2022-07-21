package terminal

// Content for a single frame
type frame struct {
	scrl *scroll

	displayLines   int
	startLine      LineAbsNo
	startLineSubNo int

	cursorX      int
	cursorY      int
	cursorLineNo LineAbsNo
	cursorPos    int

	rows int
	cols int

	tabs []int
}

func newFrame() *frame {
	frm := &frame{
		scrl: newScroll(),
	}
	frm.scrl.newLine()
	return frm
}

func (frm *frame) resize(rows, cols int) {
	if frm.displayLines >= frm.rows {
		// auto scroll
		frm.displayLines = min(rows, frm.scrl.totalDisplayLines)
	} else {
		frm.displayLines = min(rows, frm.displayLines)
	}

	frm.rows = rows
	frm.cols = cols
}

func (frm *frame) visibleLines() []*Line {
	return nil
}

func (frm *frame) printChar(r rune) {
	width := charWidth(r)

	c := &Cell{
		Char:  r,
		Width: width,
	}

	frm.scrl.putCell(c, frm.cursorLineNo, frm.cursorPos)

	//	frm.lines[frm.cursorY].addCell(&Cell{
	//		Char: r,
	//	})
	frm.cursorX += charWidth(r)
}

func (frm *frame) moveCursorX(x int) {

}

func (frm *frame) moveCursorY(x int) {

}

func (frm *frame) moveRelCursorX(x int) {

}

func (frm *frame) moveRelCursorY(y int) {

}

func (frm *frame) shiftTab() {

}
