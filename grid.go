package terminal

// Cell describes the single character
type Cell struct {
	Char  rune
	Width int
	FG    uint32
	BG    uint32
}

// Line includes the glyphs in a line
type Line struct {
	Cells []*Cell
}

// Grid for display
type Grid struct {
	Title   string
	Lines   []*Line
	CursorX int
	CursorY int
}

// AddCell adds a cell to the line
func (l *Line) AddCell(c *Cell) {
	l.Cells = append(l.Cells, c)
}
