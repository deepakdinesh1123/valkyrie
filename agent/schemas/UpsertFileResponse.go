// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    upsertFileResponse, err := UnmarshalUpsertFileResponse(bytes)
//    bytes, err = upsertFileResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalUpsertFileResponse(data []byte) (UpsertFileResponse, error) {
	var r UpsertFileResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UpsertFileResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a response to a file upsert request in a sandbox.
type UpsertFileResponse struct {
	// Name of the file that was upserted        
	FileName                             string  `json:"fileName"`
	Msg                                  string  `json:"msg"`
	// Type of the message                       
	MsgType                              *string `json:"msgType,omitempty"`
	// Path where the file was upserted          
	Path                                 string  `json:"path"`
	Success                              bool    `json:"success"`
}
