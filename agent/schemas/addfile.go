// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    addfile, err := UnmarshalAddfile(bytes)
//    bytes, err = addfile.Marshal()

package schemas

import "encoding/json"

func UnmarshalAddfile(data []byte) (Addfile, error) {
	var r Addfile
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Addfile) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Add a file to the sandbox
type Addfile struct {
	// File content                                
	Content                                string  `json:"content"`
	// Name of the file                            
	FileName                               string  `json:"file_name"`
	MsgType                                *string `json:"msgType,omitempty"`
	// Path where to create the file               
	Path                                   string  `json:"path"`
	// ID of the sandbox to add the file to        
	SandboxID                              int64   `json:"sandboxId"`
}
