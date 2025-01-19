// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalread, err := UnmarshalTerminalread(bytes)
//    bytes, err = terminalread.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalread(data []byte) (Terminalread, error) {
	var r Terminalread
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Terminalread) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Read from terminal
type Terminalread struct {
	MsgType                   *string  `json:"msgType,omitempty"`
	// Read timeout in seconds         
	Timeout                   *float64 `json:"timeout"`
}
