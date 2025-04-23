package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("voltig.yml")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}
	if cfg == nil {
		t.Errorf("Expected non-nil config, got nil")
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	cfg, err := LoadConfig("invalid.yml")
	if err == nil {
		t.Errorf("Expected error loading invalid config")
	}
	if cfg != nil {
		t.Errorf("Expected nil config, got %v", cfg)
	}
}

func TestLoadConfigWithInsertedPackage(t *testing.T) {
	yamlContent := `
packages:
  - name: testpkg
    manager: brew
    version: 1.0.0
`
	tmpfile, err := os.CreateTemp("", "voltig-test-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Errorf("failed to remove temp file: %v", err)
		}
	}()

	if _, err := tmpfile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	cfg, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	if len(cfg.Packages) != 1 {
		t.Errorf("Expected 1 package, got %d", len(cfg.Packages))
	}
	if len(cfg.Packages[0].Name) != 1 || cfg.Packages[0].Name[0] != "testpkg" {
		t.Errorf("Expected package name 'testpkg', got '%v'", cfg.Packages[0].Name)
	}
}

func TestLoadConfigFindsInParent(t *testing.T) {
	// Create a temp subdir
	cwd, _ := os.Getwd()
	tempDir := filepath.Join(cwd, "testsubdir")
	if err := os.Mkdir(tempDir, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.Remove(tempDir); err != nil {
			t.Errorf("failed to remove temp dir: %v", err)
		}
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			t.Errorf("failed to chdir back: %v", err)
		}
	}()

	cfg, err := LoadConfig("voltig.yml")
	if err != nil {
		t.Errorf("Expected to find voltig.yml in parent dir, got error: %v", err)
	}
	if cfg == nil || len(cfg.Packages) == 0 {
		t.Errorf("Expected to load config from parent, got: %+v", cfg)
	}
}

func TestLoadConfigEmpty(t *testing.T) {
	cfg, err := LoadConfig("empty.yml")
	if err != nil {
		t.Errorf("Expected no error loading empty config, got: %v", err)
	}
	if cfg == nil {
		t.Errorf("Expected non-nil config for empty file")
	}
	if cfg != nil && len(cfg.Packages) != 0 {
		t.Errorf("Expected 0 packages, got %d", len(cfg.Packages))
	}
}
