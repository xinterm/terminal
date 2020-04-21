package terminal

import (
	"testing"
)

func TestTotalLines(t *testing.T) {
	sb := newScrollBack()
	sb.blockLines = 2
	sb.setMaxLines(9)

	if sb.totalLines != 0 {
		t.Errorf("New line buffer size is not 0")
	}

	sb.newLines(7)
	if sb.totalLines != 7 {
		t.Errorf("Total lines not correct after newLines")
	}

	sb.newLines(20)
	if sb.totalLines != 9 {
		t.Errorf("Total lines larger than max lines")
	}

	sb.removeHeadLines(4)
	if sb.totalLines != 5 {
		t.Errorf("Total lines error after remove head lines")
	}

	sb.setMaxLines(3)
	if sb.totalLines != 3 {
		t.Errorf("Total lines error after reset max lines")
	}
}

func TestCircle(t *testing.T) {
	sb := newScrollBack()
	sb.blockLines = 3
	sb.setMaxLines(10)

	sb.newLines(13)
	if sb.endBlockNo != 0 {
		t.Errorf("End block No does not come back to 0 in circle")
	}
	if sb.startBlockNo != 1 || sb.startSubLineNo != 0 {
		t.Errorf("Start position does not shift in circle")
	}

	sb.newLines(5)
}

func TestNormalize(t *testing.T) {
	sb := newScrollBack()
	sb.blockLines = 2
	sb.setMaxLines(9)

	sb.newLines(2345)
	sb.normalize()
	if sb.startBlockNo != 0 {
		t.Errorf("Start block No is not 0 after normalize")
	}
}

func TestBufSize(t *testing.T) {
	sb := newScrollBack()
	sb.blockLines = 2

	sb.setMaxLines(9)
	sb.newLines(2345)
	if len(sb.buf) != 5 {
		t.Errorf("buf size not correct")
	}

	sb.setMaxLines(8)
	sb.newLines(123)
	if len(sb.buf) != 5 {
		t.Errorf("buf size after new max lines not correct")
	}

	sb.setMaxLines(20)
	sb.newLines(789)
	if len(sb.buf) != 11 {
		t.Errorf("buf size after new max lines not correct")
	}
}

func TestLineBuffer(t *testing.T) {
	sb := newScrollBack()
	sb.blockLines = 7
	sb.setMaxLines(29)

	for i := 0; i < 100; i++ {
		line := sb.newLine()
		line.Cells = append(line.Cells, &Cell{
			Char: rune(i),
		})
	}
	if sb.firstLine().Cells[0].Char != 71 {
		t.Errorf("First line error")
	}
	if sb.lastLine().Cells[0].Char != 99 {
		t.Errorf("Last line error")
	}
	if sb.line(30) != nil {
		t.Errorf("Line No out of range is not nil")
	}
	if sb.line(-1) != nil {
		t.Errorf("Line No negative is not nil")
	}
	for i := 0; i < 29; i++ {
		if sb.line(i).Cells[0].Char != 71+rune(i) {
			t.Error("Get line error")
		}
	}
}
