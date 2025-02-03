// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalClose, err := UnmarshalTerminalClose(bytes)
//    bytes, err = terminalClose.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalClose(data []byte) (TerminalClose, error) {
	var r TerminalClose
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TerminalClose) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Close terminal session
type TerminalClose struct {
	MsgType *string `json:"msgType,omitempty"`
	// Unique identifier for the terminal session
	TerminalID string `json:"terminalId"`
}
