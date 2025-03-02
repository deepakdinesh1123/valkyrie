// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    uninstallNixPackageResponse, err := UnmarshalUninstallNixPackageResponse(bytes)
//    bytes, err = uninstallNixPackageResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalUninstallNixPackageResponse(data []byte) (UninstallNixPackageResponse, error) {
	var r UninstallNixPackageResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UninstallNixPackageResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type UninstallNixPackageResponse struct {
	Msg     string  `json:"msg"`
	MsgType *string `json:"msgType,omitempty"`
	Success bool    `json:"success"`
}
