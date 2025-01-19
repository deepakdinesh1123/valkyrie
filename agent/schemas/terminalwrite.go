// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalwrite, err := UnmarshalTerminalwrite(bytes)
//    bytes, err = terminalwrite.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalwrite(data []byte) (Terminalwrite, error) {
	var r Terminalwrite
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Terminalwrite) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Write to terminal
type Terminalwrite struct {
	// Content to write to terminal        
	Content                        string  `json:"content"`
	MsgType                        *string `json:"msgType,omitempty"`
}
