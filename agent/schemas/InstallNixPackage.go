// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    installNixPackage, err := UnmarshalInstallNixPackage(bytes)
//    bytes, err = installNixPackage.Marshal()

package schemas

import "encoding/json"

func UnmarshalInstallNixPackage(data []byte) (InstallNixPackage, error) {
	var r InstallNixPackage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *InstallNixPackage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type InstallNixPackage struct {
	// Name of the channel           
	Channel                  *string `json:"channel,omitempty"`
	MsgType                  *string `json:"msgType,omitempty"`
	// Nix package to install        
	PkgName                  string  `json:"pkgName"`
}
