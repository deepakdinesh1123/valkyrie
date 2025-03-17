// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    upsertDirectoryResponse, err := UnmarshalUpsertDirectoryResponse(bytes)
//    bytes, err = upsertDirectoryResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalUpsertDirectoryResponse(data []byte) (UpsertDirectoryResponse, error) {
	var r UpsertDirectoryResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UpsertDirectoryResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Represents a response to a directory upsert request in a sandbox.
type UpsertDirectoryResponse struct {
	Msg                                     string  `json:"msg"`
	// Type of the message                          
	MsgType                                 *string `json:"msgType,omitempty"`
	// Path where the directory was upserted        
	Path                                    string  `json:"path"`
	Success                                 bool    `json:"success"`
}
