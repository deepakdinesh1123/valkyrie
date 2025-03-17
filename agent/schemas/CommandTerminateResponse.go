// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    commandTerminateResponse, err := UnmarshalCommandTerminateResponse(bytes)
//    bytes, err = commandTerminateResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalCommandTerminateResponse(data []byte) (CommandTerminateResponse, error) {
	var r CommandTerminateResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CommandTerminateResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CommandTerminateResponse struct {
	// Command ID        
	CommandID    string  `json:"commandId"`
	Msg          string  `json:"msg"`
	MsgType      *string `json:"msgType,omitempty"`
	Success      bool    `json:"success"`
}
