package sequence

// ResultType represents the result type
type ResultType int

// Result types
const (
	ResultASCII ResultType = iota
	ResultC0
	ResultUTF8

	ResultCSI
	ResultCS

	ResultOSC
	ResultOSCTitle
)

// Result will be passed by channel
type Result struct {
	Type  ResultType
	Value interface{}
}

// CSIResult contains the CSI sequence result
type CSIResult struct {
	param        []string
	intermediate string
	final        byte
}

// OSCResult contains the OSC sequence result
type OSCResult struct {
	param []string
}
