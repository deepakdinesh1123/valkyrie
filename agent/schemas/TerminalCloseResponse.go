// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalCloseResponse, err := UnmarshalTerminalCloseResponse(bytes)
//    bytes, err = terminalCloseResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalCloseResponse(data []byte) (TerminalCloseResponse, error) {
	var r TerminalCloseResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TerminalCloseResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Response after closing a terminal session
type TerminalCloseResponse struct {
	// Message confirming terminal closure              
	Msg                                          string `json:"msg"`
	// Success                                          
	Success                                      bool   `json:"success"`
	// Unique identifier for the terminal session       
	TerminalID                                   string `json:"terminalId"`
}
