package types

type CargoPackage struct {
	Name         string
	Version      string
	Source       string
	Dependencies []string
}

type Package struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}

// Used to marshal from yaml/json file to get the list of packages
type PackageList struct {
	Packages []Package `json:"packages" yaml:"packages"`
}
