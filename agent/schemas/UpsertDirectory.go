// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    upsertDirectory, err := UnmarshalUpsertDirectory(bytes)
//    bytes, err = upsertDirectory.Marshal()

package schemas

import "encoding/json"

func UnmarshalUpsertDirectory(data []byte) (UpsertDirectory, error) {
	var r UpsertDirectory
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UpsertDirectory) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a request to add or update a directory in a sandbox.
type UpsertDirectory struct {
	// Type of the message                                          
	MsgType                                                 *string `json:"msgType,omitempty"`
	// Path where the directory should be created or updated        
	Path                                                    string  `json:"path"`
}
