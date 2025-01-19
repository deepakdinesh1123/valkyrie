// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    newterminal, err := UnmarshalNewterminal(bytes)
//    bytes, err = newterminal.Marshal()

package schemas

import "encoding/json"

func UnmarshalNewterminal(data []byte) (Newterminal, error) {
	var r Newterminal
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Newterminal) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Create a new terminal session
type Newterminal struct {
	MsgType                   *string  `json:"msgType,omitempty"`
	// Nix flake configuration         
	NixFlake                  *string  `json:"nix_flake,omitempty"`
	// Nix shell configuration         
	NixShell                  *string  `json:"nix_shell,omitempty"`
	// Packages to install             
	Packages                  []string `json:"packages,omitempty"`
	// Shell type to use               
	Shell                     Shell    `json:"shell"`
}

// Shell type to use
type Shell string

const (
	Bash     Shell = "bash"
	Nix      Shell = "nix"
	NixShell Shell = "nix-shell"
	Sh       Shell = "sh"
)
