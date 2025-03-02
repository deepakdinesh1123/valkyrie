// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    commandWriteInputResponse, err := UnmarshalCommandWriteInputResponse(bytes)
//    bytes, err = commandWriteInputResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalCommandWriteInputResponse(data []byte) (CommandWriteInputResponse, error) {
	var r CommandWriteInputResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CommandWriteInputResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CommandWriteInputResponse struct {
	// Command ID        
	CommandID    string  `json:"commandId"`
	Msg          string  `json:"msg"`
	MsgType      *string `json:"msgType,omitempty"`
	Success      bool    `json:"success"`
}
