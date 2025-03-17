// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    uninstallNixPackage, err := UnmarshalUninstallNixPackage(bytes)
//    bytes, err = uninstallNixPackage.Marshal()

package schemas

import "encoding/json"

func UnmarshalUninstallNixPackage(data []byte) (UninstallNixPackage, error) {
	var r UninstallNixPackage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UninstallNixPackage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type UninstallNixPackage struct {
	MsgType                    *string `json:"msgType,omitempty"`
	// Nix package to uninstall        
	PkgName                    string  `json:"pkgName"`
}
