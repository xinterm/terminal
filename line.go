package terminal

type line struct {
	cells        []*Cell
	wrapPosition []int
}

func (l *line) setCell(c *Cell, pos int) {
	l.cells = append(l.cells, c)
}

func (l *line) rewrap(cols int) {
	if cols <= 0 {
		return
	}

	l.wrapPosition = l.wrapPosition[:0]

}

func (l *line) lineNumber() int {
	return len(l.wrapPosition) + 1
}
