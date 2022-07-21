package terminal

// LineNoAbs is the absolute number of line
type LineNoAbs int64

type scroll struct {
	lineMap    map[LineNoAbs]*line
	displayMap map[int64]LineNoAbs

	maxLines int

	totalLines        int
	totalDisplayLines int

	startLineNo    LineNoAbs
	endLineNo      LineNoAbs
	startDisplayNo int64
	endDisplayNo   int64
}

func newScroll() *scroll {
	return &scroll{
		lineMap:    make(map[LineNoAbs]*line),
		displayMap: make(map[int64]LineNoAbs),
	}
}

func (scrl *scroll) recalculateDisplayLineNo(cols int) {
}

func (scrl *scroll) setMaxLines(maxLines int) {
	if maxLines < 0 {
		// This is the maximum int value
		maxLines = int(^uint(0) >> 1)
	}

	scrl.maxLines = maxLines

	if scrl.totalLines > scrl.maxLines {
		scrl.removeHeadLines(scrl.totalLines - scrl.maxLines)
	}
}

func (scrl *scroll) putCell(c *Cell, lineNo LineNoAbs, pos int) {
	scrl.lineMap
}

func (scrl *scroll) newLines(number int) LineNoAbs {
	if scrl.maxLines <= 0 {
		return -1
	}

	for i := 0; i < number; i++ {
		if scrl.totalLines < scrl.maxLines {
			scrl.totalLines++
		} else {
			scrl.removeHeadLines(1)
		}

		scrl.lineMap[scrl.endLineNo] = &line{}
		scrl.endLineNo++
	}

	return scrl.endLineNo
}

func (scrl *scroll) firstLine() LineNoAbs {
	return scrl.startLineNo
}

func (scrl *scroll) lastLine() LineNoAbs {
	return scrl.endLineNo
}

func (scrl *scroll) line(displayLineNo int) LineNoAbs {
	if lineNo < 0 || lineNo >= scrl.totalLines {
		return 0
	}

	tempLineNo := scrl.startSubLineNo + lineNo
	tempBlockNo := scrl.startBlockNo + tempLineNo/scrl.blockLines

	blockNo := tempBlockNo % len(scrl.buf)
	subLineNo := tempLineNo % scrl.blockLines

	return scrl.buf[blockNo][subLineNo]
}

func (scrl *scroll) removeHeadLines(number int) {
	if number > scrl.totalLines {
		number = scrl.totalLines
	}

	for n := 0; n < number; n++ {
		for i := 0; i < scrl.lineMap[scrl.startLineNo].lineNumber(); i++ {
			delete(scrl.displayMap, scrl.startDisplayNo)
			scrl.startDisplayNo++
		}

		delete(scrl.lineMap, scrl.startLineNo)
		scrl.startLineNo++
	}
}

func (scrl *scroll) removeTailLines(number int) {
	if number > scrl.totalLines {
		number = scrl.totalLines
	}

	for i := 0; i < number; i++ {
		scrl._moveToPrevLine(&scrl.endBlockNo, &scrl.endSubLineNo)
		scrl.buf[scrl.endBlockNo][scrl.endSubLineNo] = 0 ////////
		scrl.totalLines--
	}
}
