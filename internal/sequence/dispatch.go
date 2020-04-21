package sequence

import (
	"strconv"
	"strings"

	"github.com/xinterm/terminal/util"
)

type dispatch struct {
	log    util.Logger
	result chan *Result
}

func newDispatch(log util.Logger) *dispatch {
	return &dispatch{
		log:    log,
		result: make(chan *Result, 512),
	}
}

func (dp *dispatch) dispatchASCII(c byte) {
	dp.result <- &Result{
		Type:  ResultASCII,
		Value: c,
	}
}

func (dp *dispatch) dispatchC0(c byte) {
	dp.result <- &Result{
		Type:  ResultC0,
		Value: c,
	}
}

func (dp *dispatch) dispatchUTF8(r rune) {
	dp.result <- &Result{
		Type:  ResultUTF8,
		Value: r,
	}
}

func (dp *dispatch) dispatchCSI(param, intermediate string, final byte) {
	dp.result <- &Result{
		Type: ResultCSI,
		Value: CSIResult{
			param:        strings.Split(param, ";"),
			intermediate: intermediate,
			final:        final,
		},
	}
}

func (dp *dispatch) dispatchOSC(param string) {
	dp.result <- &Result{
		Type: ResultOSC,
		Value: OSCResult{
			param: strings.Split(param, ";"),
		},
	}
}

func (dp *dispatch) getIntParams(param []string, defaultNumber int) []int {
	params := make([]int, 0)
	for _, frag := range param {
		var i int
		if frag == "" {
			i = defaultNumber
		} else {
			var err error
			i, err = strconv.Atoi(frag)
			if err != nil {
				dp.log.Debugf("Invalid CSI parameter format")
				break
			}
		}
		params = append(params, i)
	}
	return params
}

/*
func (dp *dispatch) dispatchC0Old(c byte) {
	switch c {
	case '\b':
		s.currentScreen().moveRelCursorX(-1)
	case '\t':
		s.currentScreen().shiftTab()
	case '\n':
		s.currentScreen().moveRelCursorY(1)
	case '\r':
		s.currentScreen().moveCursorX(0)
	default:
		s.log.Errorf("Unresolved C0 char: 0x%x", c)
	}
}

func (dp *dispatch) moveCursor() {
	param := s.csi.getIntParams(1)
	if len(param) == 0 {
		param = append(param, 1)
	}

	switch s.csi.final {
	case 'A':
		s.currentScreen().moveRelCursorY(-param[0])
	case 'B':
		s.currentScreen().moveRelCursorY(param[0])
	case 'C':
		s.currentScreen().moveRelCursorX(param[0])
	case 'D':
		s.currentScreen().moveRelCursorX(-param[0])
	case 'E':
		s.currentScreen().moveCursorX(0)
		s.currentScreen().moveRelCursorY(param[0])
	case 'F':
		s.currentScreen().moveCursorX(0)
		s.currentScreen().moveRelCursorY(-param[0])
	case 'G':
		s.currentScreen().moveCursorX(param[0])
	case 'H':
		if len(param) == 1 {
			param = append(param, 1)
		}
		s.currentScreen().moveCursorY(param[0])
		s.currentScreen().moveCursorX(param[1])
	}
}

func (dp *dispatch) handleCSI() {
	switch s.csi.final {
	case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H':
		s.moveCursor()
	default:
	}
}
*/
