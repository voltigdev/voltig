package models

// Package represents a software package in the system.
type Package struct {
	Name         []string
	Manager      string
	Version      string
	Optional     bool
	Dependencies []string
}

// PackageStatus represents the status of a package (installed, missing, etc.).
type PackageStatus struct {
	Name    string
	Status  string // e.g., installed, missing, outdated
	Version string
}
