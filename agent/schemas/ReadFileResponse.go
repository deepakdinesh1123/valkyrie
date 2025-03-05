// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    readFileResponse, err := UnmarshalReadFileResponse(bytes)
//    bytes, err = readFileResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalReadFileResponse(data []byte) (ReadFileResponse, error) {
	var r ReadFileResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ReadFileResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a response to a file read request in a sandbox.
type ReadFileResponse struct {
	// Content of the file that was read        
	Content                             string  `json:"content"`
	Msg                                 string  `json:"msg"`
	// Type of the message                      
	MsgType                             *string `json:"msgType,omitempty"`
	// Path of the file that was read           
	Path                                string  `json:"path"`
	Success                             bool    `json:"success"`
}
