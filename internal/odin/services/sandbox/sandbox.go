package sandbox

type FlakeConfig struct {
	// NIXPKGS_URL represents the URL for nixpkgs
	NIXPKGS_URL string `json:"nixpkgs_url"`

	// Languages is a list of programming languages/tools to be included in the development shell
	Languages []string `json:"languages"`

	// SystemDependencies is a list of system packages to be included in the development shell
	SystemDependencies []string `json:"system_dependencies"`

	Services []string `json:"services"`
}
