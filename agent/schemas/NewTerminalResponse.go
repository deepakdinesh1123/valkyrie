// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    newTerminalResponse, err := UnmarshalNewTerminalResponse(bytes)
//    bytes, err = newTerminalResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalNewTerminalResponse(data []byte) (NewTerminalResponse, error) {
	var r NewTerminalResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NewTerminalResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type NewTerminalResponse struct {
	Msg           string  `json:"msg"`
	MsgType       *string `json:"msgType,omitempty"`
	Success       bool    `json:"success"`
	// Terminal ID        
	TerminalID    string  `json:"terminalID"`
}
