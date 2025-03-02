// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    commandWriteInput, err := UnmarshalCommandWriteInput(bytes)
//    bytes, err = commandWriteInput.Marshal()

package schemas

import "encoding/json"

func UnmarshalCommandWriteInput(data []byte) (CommandWriteInput, error) {
	var r CommandWriteInput
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CommandWriteInput) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CommandWriteInput struct {
	// Command ID        
	CommandID    string  `json:"commandId"`
	Input        *string `json:"input,omitempty"`
	MsgType      *string `json:"msgType,omitempty"`
}
