// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    deleteDirectory, err := UnmarshalDeleteDirectory(bytes)
//    bytes, err = deleteDirectory.Marshal()

package schemas

import "encoding/json"

func UnmarshalDeleteDirectory(data []byte) (DeleteDirectory, error) {
	var r DeleteDirectory
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DeleteDirectory) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a request to delete a directory in a sandbox.
type DeleteDirectory struct {
	// Type of the message                        
	MsgType                               *string `json:"msgType,omitempty"`
	// Path of the directory to be deleted        
	Path                                  string  `json:"path"`
}
