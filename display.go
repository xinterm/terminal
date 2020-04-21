package terminal

// Display deals with the content of screen
type display struct {
	wrapLine bool
}

func (d *display) SetWrapLine(wrapLine bool) {
	d.wrapLine = wrapLine
}

func (d *display) render(st *state) *Grid {
	scr := st.currentScreen()

	grid := &Grid{
		Lines:   make([]*Line, 0, scr.rows),
		CursorX: scr.cursorX,
		CursorY: scr.cursorY,
	}

	for _, line := range scr.lines {
		gridLine := &Line{}
		for _, cell := range line.Cells {
			gridLine.AddCell(cell)
		}
		grid.Lines = append(grid.Lines, gridLine)
	}

	return grid
}
