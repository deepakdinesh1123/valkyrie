// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalRead, err := UnmarshalTerminalRead(bytes)
//    bytes, err = terminalRead.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalRead(data []byte) (TerminalRead, error) {
	var r TerminalRead
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TerminalRead) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Read from terminal
type TerminalRead struct {
	MsgType                                      *string `json:"msgType,omitempty"`
	// Unique identifier for the terminal session        
	TerminalID                                   string  `json:"terminalId"`
}
