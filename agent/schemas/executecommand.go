// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    executecommand, err := UnmarshalExecutecommand(bytes)
//    bytes, err = executecommand.Marshal()

package schemas

import "encoding/json"

func UnmarshalExecutecommand(data []byte) (Executecommand, error) {
	var r Executecommand
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Executecommand) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Execute a command in the sandbox environment
type Executecommand struct {
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
