// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalclose, err := UnmarshalTerminalclose(bytes)
//    bytes, err = terminalclose.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalclose(data []byte) (Terminalclose, error) {
	var r Terminalclose
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Terminalclose) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Close terminal session
type Terminalclose struct {
	MsgType *string `json:"msgType,omitempty"`
}
