// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalWrite, err := UnmarshalTerminalWrite(bytes)
//    bytes, err = terminalWrite.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalWrite(data []byte) (TerminalWrite, error) {
	var r TerminalWrite
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TerminalWrite) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Write to terminal
type TerminalWrite struct {
	// input to write to terminal                        
	Input                                        string  `json:"input"`
	MsgType                                      *string `json:"msgType,omitempty"`
	// Unique identifier for the terminal session        
	TerminalID                                   string  `json:"terminalId"`
}
