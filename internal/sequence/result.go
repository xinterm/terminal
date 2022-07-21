package sequence

// ResultType represents the result type
type ResultType int

// Result types
const (
	ResultExit ResultType = iota

	ResultASCII
	ResultC0
	ResultUTF8
	ResultEscape

	ResultCSI
	ResultCS

	ResultOSC

	ResultTitle
)

// Result will be passed by channel
type Result struct {
	Type  ResultType
	Value interface{}
}

// EscapeResult contains the escape sequence result
type EscapeResult struct {
	Char  byte
	Param string
}

// CSIResult contains the CSI sequence result
type CSIResult struct {
	Param        []string
	Intermediate string
	Final        byte
}

// OSCResult contains the OSC sequence result
type OSCResult struct {
	Param []string
}
