// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    installNixPackageResponse, err := UnmarshalInstallNixPackageResponse(bytes)
//    bytes, err = installNixPackageResponse.Marshal()

package schemas

import "encoding/json"

func UnmarshalInstallNixPackageResponse(data []byte) (InstallNixPackageResponse, error) {
	var r InstallNixPackageResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *InstallNixPackageResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type InstallNixPackageResponse struct {
	Msg     string  `json:"msg"`
	MsgType *string `json:"msgType,omitempty"`
	Success bool    `json:"success"`
}
