// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    deleteDirectoryResponse, err := UnmarshalDeleteDirectoryResponse(bytes)
//    bytes, err = deleteDirectoryResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalDeleteDirectoryResponse(data []byte) (DeleteDirectoryResponse, error) {
	var r DeleteDirectoryResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DeleteDirectoryResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a response to a directory delete request in a sandbox.
type DeleteDirectoryResponse struct {
	Msg                                      string  `json:"msg"`
	// Type of the message                           
	MsgType                                  *string `json:"msgType,omitempty"`
	// Path of the directory that was deleted        
	Path                                     string  `json:"path"`
	Success                                  bool    `json:"success"`
}
