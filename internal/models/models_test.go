package models

import (
	"testing"
	"voltig/config"
)

func TestPackageStatusEquality(t *testing.T) {
	ps1 := PackageStatus{Name: "foo", Version: "1.0.0", Status: "installed"}
	ps2 := PackageStatus{Name: "foo", Version: "1.0.0", Status: "installed"}
	ps3 := PackageStatus{Name: "bar", Version: "2.0.0", Status: "missing"}

	if ps1 != ps2 {
		t.Errorf("Expected ps1 and ps2 to be equal")
	}
	if ps1 == ps3 {
		t.Errorf("Expected ps1 and ps3 to be different")
	}
}

func TestToModel(t *testing.T) {
	ps := ToModel(config.Package{
		Name:         []string{"testpkg"},
		Version:      "1.2.3",
		Manager:      "brew",
		Optional:     false,
		Dependencies: nil,
	})
	if ps.Name[0] != "testpkg" || ps.Version != "1.2.3" {
		t.Errorf("ToModel did not copy fields correctly")
	}
}
