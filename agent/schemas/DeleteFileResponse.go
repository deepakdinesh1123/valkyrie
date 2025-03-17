// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    deleteFileResponse, err := UnmarshalDeleteFileResponse(bytes)
//    bytes, err = deleteFileResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalDeleteFileResponse(data []byte) (DeleteFileResponse, error) {
	var r DeleteFileResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DeleteFileResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a response to a file delete request in a sandbox.
type DeleteFileResponse struct {
	Msg                                 string  `json:"msg"`
	// Type of the message                      
	MsgType                             *string `json:"msgType,omitempty"`
	// Path of the file that was deleted        
	Path                                string  `json:"path"`
	Success                             bool    `json:"success"`
}
