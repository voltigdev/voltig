/*
Package shared provides shared data structures for Voltig CLI.
Will maybe used this in the future as shared package between cli and tui
*/
package shared

import (
	"voltig/config"
	"voltig/internal/models"
)

// LoadedConfigPkgs holds the list of package statuses loaded at TUI startup.
var LoadedConfigPkgs []models.PackageStatus

// LoadedConfigList holds the list of config packages loaded at TUI startup.
var LoadedConfigList []config.Package
