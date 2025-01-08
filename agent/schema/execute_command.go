// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    executeCommand, err := UnmarshalExecuteCommand(bytes)
//    bytes, err = executeCommand.Marshal()

package schema

import "encoding/json"

func UnmarshalExecuteCommand(data []byte) (ExecuteCommand, error) {
	var r ExecuteCommand
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ExecuteCommand) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ExecuteCommand struct {
	Command   string   `json:"command"`
	Env       []string `json:"env,omitempty"`
	SandboxID int64    `json:"sandboxId"`
	Stderr    *bool    `json:"stderr,omitempty"`
	Stdin     *bool    `json:"stdin,omitempty"`
	Stdout    *bool    `json:"stdout,omitempty"`
	WorkDir   *string  `json:"workDir,omitempty"`
}
