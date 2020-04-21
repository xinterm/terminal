package sequence

import (
	"github.com/xinterm/terminal/util"
)

type sequencer interface {
	shouldEnter(byte) bool
	shouldExit() bool
	reset()
	process(byte)
	startSubSequence() bool
	dispatch()
}

// SeqCtl controls the sequences
type SeqCtl struct {
	log util.Logger

	dp *dispatch

	curSeq sequencer

	childrenTree map[sequencer][]sequencer
	parentTree   map[sequencer]sequencer

	rootSeq *rootSequence

	ascSeq        *byteSequence
	c0Seq         *byteSequence
	escSeq        *escapeSequence
	utf8Seq       *utf8Sequence
	unresolvedSeq *byteSequence

	csiSeq           *csiSequence
	oscSeq           *oscSequence
	unresolvedEscSeq *byteSequence
	invalidEscSeq    *byteSequence

	csiFinalSeq *byteSequence
	csiParamSeq *byteSequence
	csiInterSeq *byteSequence

	oscParamSeq *byteSequence
	oscEscSeq   *escapeSequence
	oscSTSeq    *byteSequence
	oscBELSeq   *byteSequence
}

// NewSeqCtl creates a new SeqCtl instance
func NewSeqCtl(log util.Logger) *SeqCtl {
	ctl := &SeqCtl{
		log: log,
		dp:  newDispatch(log),
	}

	ctl.registerSequences()

	ctl.curSeq = ctl.rootSeq

	return ctl
}

// PollResult waits for the sequence result
func (ctl *SeqCtl) PollResult() *Result {
	return <-ctl.dp.result
}

func (ctl *SeqCtl) checkExitStatus() {
	if ctl.curSeq.shouldExit() {
		ctl.curSeq.dispatch()
		ctl.curSeq = ctl.parentTree[ctl.curSeq]
		ctl.checkExitStatus()
	}
}

// Parse the incoming byte
func (ctl *SeqCtl) Parse(c byte) {
	if !ctl.curSeq.startSubSequence() {
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
		ctl.Parse(c)
		return
	}
}

func (ctl *SeqCtl) addTree(parent sequencer, children ...sequencer) {
	for _, child := range children {
		ctl.childrenTree[parent] = append(ctl.childrenTree[parent], child)
		ctl.parentTree[child] = parent
	}
}

func (ctl *SeqCtl) registerSequences() {
	ctl.rootSeq = &rootSequence{}

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

	ctl.unresolvedSeq = &byteSequence{}
	ctl.unresolvedSeq.byteRange = [][2]byte{{0x00, 0xff}}
	ctl.unresolvedSeq.handler = func() {
		ctl.log.Errorf("Unresolved sequence: 0x%x", ctl.unresolvedSeq.char)
	}

	ctl.addTree(ctl.rootSeq, ctl.ascSeq, ctl.c0Seq, ctl.escSeq, ctl.utf8Seq, ctl.unresolvedSeq)
}

func (ctl *SeqCtl) registerEscapeSequences() {
	ctl.escSeq = &escapeSequence{}

	ctl.registerCSISequences()

	ctl.registerOSCSequences()

	ctl.unresolvedEscSeq = &byteSequence{}
	ctl.unresolvedEscSeq.addRange(0x40, 0x5f)
	ctl.unresolvedEscSeq.handler = func() {
		ctl.log.Errorf("Unresolved escape sequence: 0x%x", ctl.unresolvedEscSeq.char)
	}

	ctl.invalidEscSeq = &byteSequence{}
	ctl.invalidEscSeq.addRange(0x00, 0xff)
	ctl.invalidEscSeq.handler = func() {
		ctl.log.Errorf("Invalid escape sequence: 0x%x", ctl.invalidEscSeq.char)
	}

	ctl.addTree(ctl.escSeq, ctl.csiSeq, ctl.oscSeq, ctl.unresolvedEscSeq, ctl.invalidEscSeq)
}

func (ctl *SeqCtl) registerCSISequences() {
	ctl.csiSeq = &csiSequence{}
	ctl.csiSeq.handler = func() {
		ctl.escSeq.exit = true
		ctl.dp.dispatchCSI(ctl.csiSeq.param.String(), ctl.csiSeq.intermediate.String(), ctl.csiSeq.final)
	}

	ctl.csiFinalSeq = &byteSequence{}
	ctl.csiFinalSeq.byteRange = [][2]byte{{0x40, 0x7e}}
	ctl.csiFinalSeq.handler = func() {
		ctl.csiSeq.final = ctl.csiFinalSeq.char
	}

	ctl.csiParamSeq = &byteSequence{}
	ctl.csiParamSeq.byteRange = [][2]byte{{0x30, 0x3f}}
	ctl.csiParamSeq.handler = func() {
		ctl.csiSeq.param.WriteByte(ctl.csiParamSeq.char)
	}

	ctl.csiInterSeq = &byteSequence{}
	ctl.csiInterSeq.byteRange = [][2]byte{{0x20, 0x2f}}
	ctl.csiInterSeq.handler = func() {
		ctl.csiSeq.intermediate.WriteByte(ctl.csiInterSeq.char)
	}

	ctl.addTree(ctl.csiSeq, ctl.csiFinalSeq, ctl.csiParamSeq, ctl.csiInterSeq)
}

func (ctl *SeqCtl) registerOSCSequences() {
	ctl.oscSeq = &oscSequence{}
	ctl.oscSeq.handler = func() {
		ctl.oscSeq.exit = true
		ctl.dp.dispatchOSC(ctl.oscSeq.param.String())
	}

	ctl.oscEscSeq = &escapeSequence{}

	ctl.oscSTSeq = &byteSequence{}
	ctl.oscSTSeq.byteRange = [][2]byte{{0x5c, 0x5c}}
	ctl.oscSTSeq.handler = func() {
		ctl.oscEscSeq.exit = true
		ctl.oscSeq.exit = true
	}

	ctl.oscBELSeq = &byteSequence{}
	ctl.oscBELSeq.byteRange = [][2]byte{{0x07, 0x07}}
	ctl.oscBELSeq.handler = func() {
		ctl.oscSeq.exit = true
	}

	ctl.oscParamSeq = &byteSequence{}
	ctl.oscParamSeq.byteRange = [][2]byte{{0x08, 0x0d}, {0x20, 0xff}}
	ctl.oscParamSeq.handler = func() {
		ctl.oscSeq.param.WriteByte(ctl.oscParamSeq.char)
	}

	ctl.addTree(ctl.oscEscSeq, ctl.oscSTSeq)
	ctl.addTree(ctl.oscSeq, ctl.oscEscSeq, ctl.oscBELSeq, ctl.oscParamSeq)
}
