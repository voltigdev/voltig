/*
Package models provides data models for Voltig CLI.
*/
package models

import "voltig/config"

// ToModel converts a config.Package to a models.Package.
func ToModel(pkg config.Package) Package {
	return Package{
		Name:         []string(pkg.Name),
		Manager:      pkg.Manager,
		Version:      pkg.Version,
		Optional:     pkg.Optional,
		Dependencies: pkg.Dependencies,
	}
}
