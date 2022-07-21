package terminal

import (
	"testing"
)

func TestTotalLines(t *testing.T) {
	ln := newLineBuffer()
	ln.blockLines = 2
	ln.setMaxLines(9)

	if ln.totalLines != 0 {
		t.Errorf("New line buffer size is not 0")
	}

	ln.newLines(7)
	if ln.totalLines != 7 {
		t.Errorf("Total lines not correct after newLines")
	}

	ln.newLines(20)
	if ln.totalLines != 9 {
		t.Errorf("Total lines larger than max lines")
	}

	ln.removeHeadLines(4)
	if ln.totalLines != 5 {
		t.Errorf("Total lines error after remove head lines")
	}

	ln.setMaxLines(3)
	if ln.totalLines != 3 {
		t.Errorf("Total lines error after reset max lines")
	}
}

func TestCircle(t *testing.T) {
	ln := newLineBuffer()
	ln.blockLines = 3
	ln.setMaxLines(10)

	ln.newLines(13)
	if ln.endBlockNo != 0 {
		t.Errorf("End block No does not come back to 0 in circle")
	}
	if ln.startBlockNo != 1 || ln.startSubLineNo != 0 {
		t.Errorf("Start position does not shift in circle")
	}

	ln.newLines(5)
}

func TestNormalize(t *testing.T) {
	ln := newLineBuffer()
	ln.blockLines = 2
	ln.setMaxLines(9)

	ln.newLines(2345)
	ln.normalize()
	if ln.startBlockNo != 0 {
		t.Errorf("Start block No is not 0 after normalize")
	}
}

func TestBufSize(t *testing.T) {
	ln := newLineBuffer()
	ln.blockLines = 2

	ln.setMaxLines(9)
	ln.newLines(2345)
	if len(ln.buf) != 5 {
		t.Errorf("buf size not correct")
	}

	ln.setMaxLines(8)
	ln.newLines(123)
	if len(ln.buf) != 5 {
		t.Errorf("buf size after new max lines not correct")
	}

	ln.setMaxLines(20)
	ln.newLines(789)
	if len(ln.buf) != 11 {
		t.Errorf("buf size after new max lines not correct")
	}
}

func TestLineBuffer(t *testing.T) {
	ln := newLineBuffer()
	ln.blockLines = 7
	ln.setMaxLines(29)

	for i := 0; i < 100; i++ {
		line := ln.newLine()
		line.cells = append(line.cells, &Cell{
			Char: rune(i),
		})
	}
	if ln.firstLine().cells[0].Char != 71 {
		t.Errorf("First line error")
	}
	if ln.lastLine().cells[0].Char != 99 {
		t.Errorf("Last line error")
	}
	if ln.line(30) != nil {
		t.Errorf("Line No out of range is not nil")
	}
	if ln.line(-1) != nil {
		t.Errorf("Line No negative is not nil")
	}
	for i := 0; i < 29; i++ {
		if ln.line(i).cells[0].Char != 71+rune(i) {
			t.Error("Get line error")
		}
	}
}
