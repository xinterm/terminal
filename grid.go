package terminal

// Cell describes the single character
type Cell struct {
	Char      rune
	Width     int
	FG        uint32
	BG        uint32
	Bold      bool
	Underline bool
	Italic    bool
	Blink     bool
}

// Line includes the glyphs in a line
type Line []*Cell

// Grid for display
type Grid struct {
	Title         string
	Lines         []*Line
	CursorX       int
	CursorY       int
	CursorVisible bool
	ScrollUp      int
}
