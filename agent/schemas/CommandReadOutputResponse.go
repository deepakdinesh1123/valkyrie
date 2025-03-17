// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    commandReadOutputResponse, err := UnmarshalCommandReadOutputResponse(bytes)
//    bytes, err = commandReadOutputResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalCommandReadOutputResponse(data []byte) (CommandReadOutputResponse, error) {
	var r CommandReadOutputResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CommandReadOutputResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CommandReadOutputResponse struct {
	// Command ID        
	CommandID    string  `json:"commandId"`
	Msg          string  `json:"msg"`
	MsgType      *string `json:"msgType,omitempty"`
	Stdout       string  `json:"stdout"`
	Success      bool    `json:"success"`
}
