// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    readDirectoryResponse, err := UnmarshalReadDirectoryResponse(bytes)
//    bytes, err = readDirectoryResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalReadDirectoryResponse(data []byte) (ReadDirectoryResponse, error) {
	var r ReadDirectoryResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ReadDirectoryResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a response to a directory read request in a sandbox.
type ReadDirectoryResponse struct {
	// Content of the directory that was read        
	Contents                                 string  `json:"contents"`
	Msg                                      string  `json:"msg"`
	// Type of the message                           
	MsgType                                  *string `json:"msgType,omitempty"`
	// Path of the directory that was read           
	Path                                     string  `json:"path"`
	Success                                  bool    `json:"success"`
}
