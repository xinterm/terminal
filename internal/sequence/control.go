package sequence

import (
	"github.com/xinterm/terminal/util"
)

// Control the sequences
type Control struct {
	log util.Logger

	dp *dispatch

	curSeq sequencer

	childrenTree map[sequencer][]sequencer
	parentTree   map[sequencer]sequencer

	rootSeq *rootSequence

	ascSeq  *byteSequence
	c0Seq   *byteSequence
	escSeq  *escapeSequence
	utf8Seq *utf8Sequence

	csiSeq    *csiSequence
	oscSeq    *escapeSTSequence
	escSubSeq *escapeSTSequence
}

// NewControl creates a new Control instance
func NewControl(log util.Logger) *Control {
	ctl := &Control{
		log:          log,
		dp:           newDispatch(log),
		childrenTree: make(map[sequencer][]sequencer),
		parentTree:   make(map[sequencer]sequencer),
	}

	ctl.registerSequences()

	ctl.curSeq = ctl.rootSeq

	return ctl
}

// ResultEvent return the sequence result channel
func (ctl *Control) ResultEvent() chan *Result {
	return ctl.dp.result
}

func (ctl *Control) checkExitStatus() {
	if ctl.curSeq == nil {
		return
	}

	if ctl.curSeq.shouldExit() {
		ctl.curSeq.handle()
		ctl.curSeq = ctl.parentTree[ctl.curSeq]
		ctl.checkExitStatus()
	}
}

func (ctl *Control) parse(c byte) {
	if ctl.curSeq == nil {
		return
	}

	if !ctl.curSeq.shouldHandOver() {
		ctl.curSeq.process(c)
		ctl.checkExitStatus()
		return
	}

	for _, child := range ctl.childrenTree[ctl.curSeq] {
		if !child.shouldEnter(c) {
			continue
		}
		ctl.curSeq = child
		ctl.curSeq.reset()
		ctl.parse(c)
		return
	}
}

// Parse the incoming byte
func (ctl *Control) Parse(c byte) {
	ctl.parse(c)
}

func (ctl *Control) addTree(parent sequencer, children ...sequencer) {
	for _, child := range children {
		ctl.childrenTree[parent] = append(ctl.childrenTree[parent], child)
		ctl.parentTree[child] = parent
	}
}

func (ctl *Control) registerSequences() {
	ctl.rootSeq = &rootSequence{}

	ignoredSeq := &byteSequence{}
	ignoredSeq.byteRange = [][2]byte{{0x7f, 0x7f}}

	ctl.ascSeq = &byteSequence{}
	ctl.ascSeq.byteRange = [][2]byte{{0x20, 0x7e}}
	ctl.ascSeq.handler = func() {
		ctl.dp.dispatchASCII(ctl.ascSeq.char)
	}

	ctl.c0Seq = &byteSequence{}
	ctl.c0Seq.byteRange = [][2]byte{{0x0, 0x17}, {0x19, 0x19}, {0x1c, 0x1f}}
	ctl.c0Seq.handler = func() {
		ctl.dp.dispatchC0(ctl.c0Seq.char)
	}

	ctl.registerEscapeSequences()

	ctl.utf8Seq = &utf8Sequence{}
	ctl.utf8Seq.handler = func() {
		ctl.dp.dispatchUTF8(ctl.utf8Seq.char)
	}

	ctl.addTree(ctl.rootSeq, ignoredSeq, ctl.ascSeq, ctl.c0Seq, ctl.escSeq, ctl.utf8Seq)
}

func (ctl *Control) registerEscapeSequences() {
	ctl.escSeq = &escapeSequence{}

	ctl.registerCSISequences()
	ctl.registerOSCSequences()
	ctl.registerEscapeSubSequences()

	invalidEscSeq := &byteSequence{}
	invalidEscSeq.addRange(0x00, 0xff)
	invalidEscSeq.handler = func() {
		ctl.escSeq.exit = true
		ctl.log.Errorf("Invalid escape sequence: ESC 0x%x", invalidEscSeq.char)
	}

	ctl.addTree(ctl.escSeq, ctl.csiSeq, ctl.oscSeq, ctl.escSubSeq, invalidEscSeq)
}

func (ctl *Control) registerCSISequences() {
	ctl.csiSeq = &csiSequence{}
	ctl.csiSeq.handler = func() {
		ctl.escSeq.exit = true
		ctl.dp.dispatchCSI(ctl.csiSeq.param.String(), ctl.csiSeq.intermediate.String(), ctl.csiSeq.final)
	}

	csiFinalSeq := &byteSequence{}
	csiFinalSeq.byteRange = [][2]byte{{0x40, 0x7e}}
	csiFinalSeq.handler = func() {
		ctl.csiSeq.final = csiFinalSeq.char
	}

	csiParamSeq := &byteSequence{}
	csiParamSeq.byteRange = [][2]byte{{0x30, 0x3f}}
	csiParamSeq.handler = func() {
		ctl.csiSeq.param.WriteByte(csiParamSeq.char)
	}

	csiInterSeq := &byteSequence{}
	csiInterSeq.byteRange = [][2]byte{{0x20, 0x2f}}
	csiInterSeq.handler = func() {
		ctl.csiSeq.intermediate.WriteByte(csiInterSeq.char)
	}

	ctl.addTree(ctl.csiSeq, csiFinalSeq, csiParamSeq, csiInterSeq)
}

func (ctl *Control) registerOSCSequences() {
	ctl.oscSeq = &escapeSTSequence{}
	ctl.oscSeq.byteRange = [][2]byte{{0x5d, 0x5d}}
	ctl.oscSeq.handler = func() {
		ctl.escSeq.exit = true
		ctl.dp.dispatchOSC(ctl.oscSeq.param.String())
	}

	oscEscSeq := &escapeSequence{}
	oscEscSeq.handler = func() {
		ctl.oscSeq.exit = true
	}

	oscSTSeq := &byteSequence{}
	oscSTSeq.byteRange = [][2]byte{{0x5c, 0x5c}}
	oscSTSeq.handler = func() {
		oscEscSeq.exit = true
	}

	oscBELSeq := &byteSequence{}
	oscBELSeq.byteRange = [][2]byte{{0x07, 0x07}}
	oscBELSeq.handler = func() {
		ctl.oscSeq.exit = true
	}

	oscParamSeq := &byteSequence{}
	oscParamSeq.byteRange = [][2]byte{{0x08, 0x0d}, {0x20, 0xff}}
	oscParamSeq.handler = func() {
		ctl.oscSeq.param.WriteByte(oscParamSeq.char)
	}

	ctl.addTree(oscEscSeq, oscSTSeq)
	ctl.addTree(ctl.oscSeq, oscEscSeq, oscBELSeq, oscParamSeq)
}

func (ctl *Control) registerEscapeSubSequences() {
	ctl.escSubSeq = &escapeSTSequence{}
	ctl.escSubSeq.byteRange = [][2]byte{{0x30, 0x7e}}
	ctl.escSubSeq.handler = func() {
		ctl.escSeq.exit = true
		ctl.dp.dispatchEscape(ctl.escSubSeq.char, ctl.escSubSeq.param.String())
	}

	escSubEscSeq := &escapeSequence{}
	escSubEscSeq.handler = func() {
		ctl.escSubSeq.exit = true
	}

	escSubSTSeq := &byteSequence{}
	escSubSTSeq.byteRange = [][2]byte{{0x5c, 0x5c}}
	escSubSTSeq.handler = func() {
		escSubEscSeq.exit = true
	}

	escSubParamSeq := &byteSequence{}
	escSubParamSeq.byteRange = [][2]byte{{0x08, 0x0d}, {0x20, 0xff}}
	escSubParamSeq.handler = func() {
		ctl.escSubSeq.param.WriteByte(escSubParamSeq.char)
	}

	ctl.addTree(escSubEscSeq, escSubSTSeq)
	ctl.addTree(ctl.escSubSeq, escSubEscSeq, escSubParamSeq)
}
