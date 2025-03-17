// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    deleteFile, err := UnmarshalDeleteFile(bytes)
//    bytes, err = deleteFile.Marshal()

package schemas

import "encoding/json"

func UnmarshalDeleteFile(data []byte) (DeleteFile, error) {
	var r DeleteFile
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DeleteFile) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a request to delete a file in a sandbox.
type DeleteFile struct {
	// Type of the message                   
	MsgType                          *string `json:"msgType,omitempty"`
	// Path of the file to be deleted        
	Path                             string  `json:"path"`
}
