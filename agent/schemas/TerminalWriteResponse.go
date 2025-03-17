// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalWriteResponse, err := UnmarshalTerminalWriteResponse(bytes)
//    bytes, err = terminalWriteResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalWriteResponse(data []byte) (TerminalWriteResponse, error) {
	var r TerminalWriteResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TerminalWriteResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Response after writing to terminal
type TerminalWriteResponse struct {
	Msg                                          string  `json:"msg"`
	MsgType                                      *string `json:"msgType,omitempty"`
	Success                                      bool    `json:"success"`
	// Unique identifier for the terminal session        
	TerminalID                                   string  `json:"terminalId"`
}
