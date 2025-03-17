// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    commandTerminate, err := UnmarshalCommandTerminate(bytes)
//    bytes, err = commandTerminate.Marshal()

package schemas

import "encoding/json"

func UnmarshalCommandTerminate(data []byte) (CommandTerminate, error) {
	var r CommandTerminate
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CommandTerminate) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CommandTerminate struct {
	// Command ID        
	CommandID    string  `json:"commandId"`
	MsgType      *string `json:"msgType,omitempty"`
}
