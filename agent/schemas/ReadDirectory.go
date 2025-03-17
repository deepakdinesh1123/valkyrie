// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    readDirectory, err := UnmarshalReadDirectory(bytes)
//    bytes, err = readDirectory.Marshal()

package schemas

import "encoding/json"

func UnmarshalReadDirectory(data []byte) (ReadDirectory, error) {
	var r ReadDirectory
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ReadDirectory) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a request to read a directory from a sandbox.
type ReadDirectory struct {
	// Type of the message                     
	MsgType                            *string `json:"msgType,omitempty"`
	// Path of the directory to be read        
	Path                               string  `json:"path"`
}
