package models

type Package struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	PkgType  string `json:"pkgType"`
	Language string `json:"language"`
}

type SearchResult struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Language string `json:"language"`
}
