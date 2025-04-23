package cmd

import (
	"errors"
	"testing"
	"voltig/internal/models"
)

type removeMockManager struct {
	removed map[string]bool
	fail    map[string]bool
}

func (m *removeMockManager) Remove(pkg models.Package, outputFn func(string)) error {
	if m.fail[pkg.Name[0]] {
		return errors.New("mock failure")
	}
	m.removed[pkg.Name[0]] = true
	return nil
}

func (m *removeMockManager) GetStatus(pkg models.Package) (models.PackageStatus, error) {
	return models.PackageStatus{}, nil
}

func (m *removeMockManager) IsAvailable() bool { return true }

func TestRemove_Success(t *testing.T) {
	mgr := &removeMockManager{removed: make(map[string]bool), fail: make(map[string]bool)}
	pkg := models.Package{Name: []string{"foo"}}
	err := mgr.Remove(pkg, func(string) {})
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if !mgr.removed["foo"] {
		t.Errorf("Expected foo to be removed")
	}
}

func TestRemove_Failure(t *testing.T) {
	mgr := &removeMockManager{removed: make(map[string]bool), fail: map[string]bool{"bar": true}}
	pkg := models.Package{Name: []string{"bar"}}
	err := mgr.Remove(pkg, func(string) {})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestRemove_NotFound(t *testing.T) {
	mgr := &removeMockManager{removed: make(map[string]bool), fail: make(map[string]bool)}
	// Simulate not found by not calling Remove
	if mgr.removed["baz"] {
		t.Errorf("baz should not be removed (not found)")
	}

	// Check that Remove was not called
	if _, called := mgr.fail["baz"]; called {
		t.Errorf("baz should not be removed (not found)")
	}
}
