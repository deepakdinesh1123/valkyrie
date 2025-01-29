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
	// Command to execute                                               
	Command                                       string                `json:"command"`
	// Environment variables                                            
	Env                                           []EnvironmentVariable `json:"env,omitempty"`
	MsgType                                       *string               `json:"msgType,omitempty"`
	// ID of the sandbox to execute the command in                      
	SandboxID                                     int64                 `json:"sandboxId"`
	// Enable stderr                                                    
	Stderr                                        *bool                 `json:"stderr,omitempty"`
	// Enable stdin                                                     
	Stdin                                         *bool                 `json:"stdin,omitempty"`
	// Enable stdout                                                    
	Stdout                                        *bool                 `json:"stdout,omitempty"`
	// Working directory for command execution                          
	WorkDir                                       *string               `json:"workDir,omitempty"`
}

// Environment variable configuration
type EnvironmentVariable struct {
	// Environment variable name        
	Key                          string `json:"key"`
	// Environment variable value       
	Value                        string `json:"value"`
}
