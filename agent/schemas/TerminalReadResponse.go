// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    terminalReadResponse, err := UnmarshalTerminalReadResponse(bytes)
//    bytes, err = terminalReadResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalTerminalReadResponse(data []byte) (TerminalReadResponse, error) {
	var r TerminalReadResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TerminalReadResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Response after reading from terminal
type TerminalReadResponse struct {
	// Indicates if the end of the stream has been reached       
	EOF                                                   *bool  `json:"eof,omitempty"`
	// optional message                                          
	Msg                                                   string `json:"msg"`
	// Content read from the terminal                            
	Output                                                string `json:"output"`
	// Success                                                   
	Success                                               bool   `json:"success"`
	// Unique identifier for the terminal session                
	TerminalID                                            string `json:"terminalId"`
}
