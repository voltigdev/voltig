/*
Package config contains the configuration loading and parsing logic.
*/
package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

/*
CustomCommand represents a custom command to be executed.
*/
type CustomCommand struct {
	Summary     string   `yaml:"summary"`
	Command     string   `yaml:"command,omitempty"`
	Script      string   `yaml:"script,omitempty"`
	Args        []string `yaml:"args,omitempty"`
	WorkDir     string   `yaml:"workDir,omitempty"`
	Environment []string `yaml:"environment,omitempty"`
	Shell       string   `yaml:"shell,omitempty"`
}

/*
PackageConfig represents the configuration for packages.
*/
type PackageConfig struct {
	Packages []Package                `yaml:"packages"`
	Commands map[string]CustomCommand `yaml:"commands"`
}

/*
PackageName represents a package name that can be either a string or an array of strings.
*/
type PackageName []string

// UnmarshalYAML implements the yaml.Unmarshaler interface for PackageName.
func (p *PackageName) UnmarshalYAML(value *yaml.Node) error {
	// Try to unmarshal as a string
	var s string
	if err := value.Decode(&s); err == nil {
		*p = []string{s}
		return nil
	}

	// Try to unmarshal as an array of strings
	var a []string
	if err := value.Decode(&a); err != nil {
		return fmt.Errorf("package name must be a string or an array of strings: %w", err)
	}
	*p = a
	return nil
}

/*
Package represents a package to be installed.
*/
type Package struct {
	Name         PackageName `yaml:"name"`
	Manager      string      `yaml:"manager"`
	Version      string      `yaml:"version"`
	Optional     bool        `yaml:"optional"`
	Dependencies []string    `yaml:"dependencies"`
}

/*
LoadConfig loads the configuration from a YAML file.
*/
func LoadConfig(path string) (*PackageConfig, error) {
	// Try direct open first
	f, err := os.Open(path)
	if err != nil {
		// If not found, search upward for voltig.yml
		if os.IsNotExist(err) {
			if found, ok := findConfigInParents("voltig.yml"); ok {
				f, err = os.Open(found)
			}
		}
		if f == nil || err != nil {
			return nil, err
		}
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close file: %v\n", err)
		}
	}()
	var cfg PackageConfig
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		if err == io.EOF {
			// Return empty config if file is empty
			return &PackageConfig{}, nil
		}
		return nil, err
	}
	return &cfg, nil
}

/*
findConfigInParents searches upwards for the given filename, starting from the current working directory.
*/
func findConfigInParents(filename string) (string, bool) {
	dir, _ := os.Getwd()
	for {
		candidate := filepath.Join(dir, filename)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root
		}
		dir = parent
	}
	return "", false
}
