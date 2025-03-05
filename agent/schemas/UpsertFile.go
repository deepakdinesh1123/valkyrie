// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    upsertFile, err := UnmarshalUpsertFile(bytes)
//    bytes, err = upsertFile.Marshal()

package schemas

import "encoding/json"

func UnmarshalUpsertFile(data []byte) (UpsertFile, error) {
	var r UpsertFile
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UpsertFile) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a request to add or update a file in a sandbox.
type UpsertFile struct {
	// Content of the file                                     
	Content                                            string  `json:"content"`
	// Name of the file to be added or updated                 
	FileName                                           string  `json:"fileName"`
	// Type of the message                                     
	MsgType                                            *string `json:"msgType,omitempty"`
	// Diff patch to apply to the file                         
	Patch                                              string  `json:"patch"`
	// Path where the file should be created or updated        
	Path                                               string  `json:"path"`
}
