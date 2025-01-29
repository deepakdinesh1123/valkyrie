// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    addFile, err := UnmarshalAddFile(bytes)
//    bytes, err = addFile.Marshal()

package schemas

import "encoding/json"

func UnmarshalAddFile(data []byte) (AddFile, error) {
	var r AddFile
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddFile) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Add a file to the sandbox
type AddFile struct {
	// File content                                
	Content                                string  `json:"content"`
	// Name of the file                            
	FileName                               string  `json:"fileName"`
	MsgType                                *string `json:"msgType,omitempty"`
	// Path where to create the file               
	Path                                   string  `json:"path"`
	// ID of the sandbox to add the file to        
	SandboxID                              int64   `json:"sandboxId"`
}
