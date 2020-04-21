package terminal

const (
	defaultBlockLines = 512
)

type scrollBack struct {
	buf [][]*Line

	// This could not be changed after any new line created
	blockLines int

	maxLines  int
	maxBlocks int

	startBlockNo   int
	startSubLineNo int
	endBlockNo     int
	endSubLineNo   int

	totalLines int
}

func newScrollBack() *scrollBack {
	return &scrollBack{
		blockLines: defaultBlockLines,
	}
}

func (sb *scrollBack) setMaxLines(maxLines int) {
	if maxLines < 0 {
		// This is the maximum int value
		maxLines = int(^uint(0) >> 1)
	}

	sb.maxLines = maxLines

	// This maxBlocks assures that startBlockNo and endBlockNo do not overlap in a circle
	sb.maxBlocks = (maxLines + 2*sb.blockLines - 2) / sb.blockLines

	if sb.totalLines > sb.maxLines {
		sb.removeHeadLines(sb.totalLines - sb.maxLines)
	}
}

// Rotate the sb.buf slice to startBlockNo == 0
// and delete extra elements
func (sb *scrollBack) normalize() {
	step := sb.startBlockNo

	if step != 0 {
		bufLen := len(sb.buf)
		temp := sb.buf[0]
		index := 0
		startIndex := 0
		for i := range sb.buf {
			swapIndex := (index + step) % bufLen
			if swapIndex == startIndex {
				sb.buf[index] = temp
				if i == bufLen-1 {
					break
				}
				startIndex++
				index = startIndex
				temp = sb.buf[index]
			} else {
				sb.buf[index] = sb.buf[swapIndex]
				index = swapIndex
			}
		}

		sb.startBlockNo = 0
		sb.endBlockNo = (sb.endBlockNo + bufLen - step) % bufLen
	}

	sb.buf = sb.buf[:sb.endBlockNo+1]
}

func (sb *scrollBack) _moveToNextLine(blockNo, subLineNo *int) {
	*subLineNo++
	if *subLineNo >= sb.blockLines {
		*subLineNo -= sb.blockLines
		*blockNo++
		if *blockNo >= sb.maxBlocks {
			*blockNo -= sb.maxBlocks
		}
	}
}

func (sb *scrollBack) _moveToPrevLine(blockNo, subLineNo *int) {
	*subLineNo--
	if *subLineNo < 0 {
		*subLineNo += sb.blockLines
		*blockNo--
		if *blockNo < 0 {
			*blockNo += sb.maxBlocks
		}
	}
}

func (sb *scrollBack) newLine() *Line {
	if sb.maxLines <= 0 {
		return nil
	}

	if sb.totalLines < sb.maxLines {
		sb.totalLines++
	} else {
		sb.buf[sb.startBlockNo][sb.startSubLineNo] = nil
		sb._moveToNextLine(&sb.startBlockNo, &sb.startSubLineNo)
	}

	if sb.endBlockNo >= len(sb.buf) {
		sb.buf = append(sb.buf, make([]*Line, sb.blockLines))
	}

	line := &Line{}
	sb.buf[sb.endBlockNo][sb.endSubLineNo] = line

	sb._moveToNextLine(&sb.endBlockNo, &sb.endSubLineNo)

	return line
}

func (sb *scrollBack) newLines(number int) {
	for i := 0; i < number; i++ {
		line := sb.newLine()
		if line == nil {
			break
		}
	}
}

func (sb *scrollBack) firstLine() *Line {
	if sb.totalLines <= 0 {
		return nil
	}
	return sb.buf[sb.startBlockNo][sb.startSubLineNo]
}

func (sb *scrollBack) lastLine() *Line {
	if sb.totalLines <= 0 {
		return nil
	}
	return sb.line(sb.totalLines - 1)
}

func (sb *scrollBack) line(lineNo int) *Line {
	if lineNo < 0 || lineNo >= sb.totalLines {
		return nil
	}

	tempLineNo := sb.startSubLineNo + lineNo
	tempBlockNo := sb.startBlockNo + tempLineNo/sb.blockLines

	blockNo := tempBlockNo % len(sb.buf)
	subLineNo := tempLineNo % sb.blockLines

	return sb.buf[blockNo][subLineNo]
}

func (sb *scrollBack) removeHeadLines(number int) {
	if number > sb.totalLines {
		number = sb.totalLines
	}

	for i := 0; i < number; i++ {
		sb.buf[sb.startBlockNo][sb.startSubLineNo] = nil
		sb._moveToNextLine(&sb.startBlockNo, &sb.startSubLineNo)
		sb.totalLines--
	}

	sb.normalize()
}

func (sb *scrollBack) removeTailLines(number int) {
	if number > sb.totalLines {
		number = sb.totalLines
	}

	for i := 0; i < number; i++ {
		sb._moveToPrevLine(&sb.endBlockNo, &sb.endSubLineNo)
		sb.buf[sb.endBlockNo][sb.endSubLineNo] = nil
		sb.totalLines--
	}
}
