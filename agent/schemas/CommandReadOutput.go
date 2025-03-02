// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    commandReadOutput, err := UnmarshalCommandReadOutput(bytes)
//    bytes, err = commandReadOutput.Marshal()

package schemas

import "encoding/json"

func UnmarshalCommandReadOutput(data []byte) (CommandReadOutput, error) {
	var r CommandReadOutput
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CommandReadOutput) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CommandReadOutput struct {
	// Command ID        
	CommandID    string  `json:"commandId"`
	MsgType      *string `json:"msgType,omitempty"`
}
