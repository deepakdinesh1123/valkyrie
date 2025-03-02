// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    newTerminal, err := UnmarshalNewTerminal(bytes)
//    bytes, err = newTerminal.Marshal()

package schemas

import "encoding/json"

func UnmarshalNewTerminal(data []byte) (NewTerminal, error) {
	var r NewTerminal
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *NewTerminal) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Create a new terminal session
type NewTerminal struct {
	// Environment variables to be added                  
	Env                                 map[string]string `json:"env"`
	MsgType                             *string           `json:"msgType,omitempty"`
	// Nix flake configuration                            
	NixFlake                            *string           `json:"nixFlake"`
	// Nix shell configuration                            
	NixShell                            *string           `json:"nixShell"`
	// Packages to install                                
	Packages                            []string          `json:"packages"`
	// Shell type to use                                  
	Shell                               *Shell            `json:"shell"`
}

type Shell string

const (
	Bash     Shell = "bash"
	Nix      Shell = "nix"
	NixShell Shell = "nix-shell"
	Sh       Shell = "sh"
)
