// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    executeCommandResponse, err := UnmarshalExecuteCommandResponse(bytes)
//    bytes, err = executeCommandResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalExecuteCommandResponse(data []byte) (ExecuteCommandResponse, error) {
	var r ExecuteCommandResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ExecuteCommandResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ExecuteCommandResponse struct {
	// Command ID        
	CommandID    string  `json:"commandId"`
	Msg          string  `json:"msg"`
	MsgType      *string `json:"msgType,omitempty"`
	State        *State  `json:"state,omitempty"`
	// stdout            
	Stdout       string  `json:"stdout"`
	Success      bool    `json:"success"`
}

type State string

const (
	Exited   State = "exited"
	Running  State = "running"
	Starting State = "starting"
	Stopped  State = "stopped"
)
