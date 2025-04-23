package manager

import (
	"errors"
	"testing"
	"voltig/internal/models"
)

// mockManager is a mock implementation of PackageManager for testing
// It tracks installed packages in a map.
type mockManager struct {
	installed map[string]bool
}

func (m *mockManager) Install(pkg models.Package, _ func(string)) error {
	m.installed[pkg.Name[0]] = true
	return nil
}
func (m *mockManager) Update(pkg models.Package, _ func(string)) error {
	m.installed[pkg.Name[0]] = true
	return nil
}
func (m *mockManager) Remove(pkg models.Package, _ func(string)) error {
	m.installed[pkg.Name[0]] = false
	return nil
}
func (m *mockManager) GetStatus(pkg models.Package) (models.PackageStatus, error) {
	if m.installed[pkg.Name[0]] {
		return models.PackageStatus{Name: pkg.Name[0], Version: pkg.Version, Status: "installed"}, nil
	}
	return models.PackageStatus{Name: pkg.Name[0], Version: pkg.Version, Status: "missing"}, errors.New("not installed")
}
func (m *mockManager) IsAvailable() bool { return true }


func TestInstall(t *testing.T) {
	mgr := &mockManager{installed: make(map[string]bool)}
	pkg := models.Package{Name: []string{"foo"}, Version: "1.0.0"}
	if err := mgr.Install(pkg, nil); err != nil {
		t.Fatalf("Install failed: %v", err)
	}
	if !mgr.installed["foo"] {
		t.Errorf("Expected foo to be installed")
	}
}

func TestUpdate(t *testing.T) {
	mgr := &mockManager{installed: make(map[string]bool)}
	pkg := models.Package{Name: []string{"foo"}, Version: "1.0.0"}
	if err := mgr.Update(pkg, nil); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if !mgr.installed["foo"] {
		t.Errorf("Expected foo to be installed after update")
	}

	if mgr.installed["bar"] {
		t.Errorf("Expected bar to not be installed after update")
	}
}

func TestRemove(t *testing.T) {
	mgr := &mockManager{installed: map[string]bool{"foo": true}}
	pkg := models.Package{Name: []string{"foo"}, Version: "1.0.0"}
	if err := mgr.Remove(pkg, nil); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}
	if mgr.installed["foo"] {
		t.Errorf("Expected foo to be removed")
	}
}

func TestGetStatus(t *testing.T) {
	mgr := &mockManager{installed: map[string]bool{"foo": true}}
	pkg := models.Package{Name: []string{"foo"}, Version: "1.0.0"}
	status, err := mgr.GetStatus(pkg)
	if err != nil || status.Status != "installed" {
		t.Errorf("Expected installed, got %v, err=%v", status.Status, err)
	}
	pkgMissing := models.Package{Name: []string{"bar"}, Version: "1.0.0"}
	status, err = mgr.GetStatus(pkgMissing)
	if err == nil || status.Status != "missing" {
		t.Errorf("Expected missing, got %v, err=%v", status.Status, err)
	}
}

func TestIsAvailable(t *testing.T) {
	mgr := &mockManager{}
	if !mgr.IsAvailable() {
		t.Error("Expected IsAvailable to return true")
	}
}

func TestPkgOperation(t *testing.T) {
	mgr := &mockManager{installed: make(map[string]bool)}
	pkgs := []models.Package{{Name: []string{"foo"}}, {Name: []string{"bar"}}}
	dummyOp := func(pkg models.Package, _ func(string)) error {
		mgr.installed[pkg.Name[0]] = true
		return nil
	}
	successes, failures := PkgOperation("Install", "installed", pkgs, dummyOp)
	if len(successes) != 2 || len(failures) != 0 {
		t.Errorf("Expected all installs to succeed, got successes=%v failures=%v", successes, failures)
	}
	// Simulate one failure
	failOp := func(pkg models.Package, _ func(string)) error {
		if pkg.Name[0] == "bar" {
			return errors.New("fail")
		}
		mgr.installed[pkg.Name[0]] = true
		return nil
	}
	successes, failures = PkgOperation("Install", "installed", pkgs, failOp)
	if len(successes) != 1 || len(failures) != 1 {
		t.Errorf("Expected one success and one failure, got successes=%v failures=%v", successes, failures)
	}
	// Explicitly cover the return statement (line 67 in manager.go)
	var emptyPkgs []models.Package
	ss, ff := PkgOperation("Install", "installed", emptyPkgs, dummyOp)
	if len(ss) != 0 || len(ff) != 0 {
		t.Errorf("Expected empty slices for empty input, got successes=%v failures=%v", ss, ff)
	}
}

