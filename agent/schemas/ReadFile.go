// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    readFile, err := UnmarshalReadFile(bytes)
//    bytes, err = readFile.Marshal()

package schemas

import "encoding/json"

func UnmarshalReadFile(data []byte) (ReadFile, error) {
	var r ReadFile
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ReadFile) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a request to read a file from a sandbox.
type ReadFile struct {
	// Type of the message                
	MsgType                       *string `json:"msgType,omitempty"`
	// Path of the file to be read        
	Path                          string  `json:"path"`
}
