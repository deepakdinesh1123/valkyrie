// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    executeCommand, err := UnmarshalExecuteCommand(bytes)
//    bytes, err = executeCommand.Marshal()

package schemas

import "encoding/json"

func UnmarshalExecuteCommand(data []byte) (ExecuteCommand, error) {
	var r ExecuteCommand
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ExecuteCommand) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Execute a command in the sandbox environment
type ExecuteCommand struct {
	// Command                                                    
	Command                                   string              `json:"command"`
	// Environment variables                                      
	Env                                       []map[string]string `json:"env,omitempty"`
	MsgType                                   *string             `json:"msgType,omitempty"`
	// Enable stderr                                              
	Stderr                                    *bool               `json:"stderr,omitempty"`
	// Enable stdin                                               
	Stdin                                     *bool               `json:"stdin,omitempty"`
	// Enable stdout                                              
	Stdout                                    *bool               `json:"stdout,omitempty"`
	// Working directory for command execution                    
	WorkDir                                   *string             `json:"workDir,omitempty"`
}
