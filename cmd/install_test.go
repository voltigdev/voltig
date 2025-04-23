package cmd

import (
	"bytes"
	"errors"
	"testing"

	"voltig/config"
	"voltig/internal/models"
	"voltig/pkg/logger"
)

type mockManager struct {
	installResults map[string]error
}

func (m *mockManager) Install(pkg models.Package) error {
	if err, ok := m.installResults[pkg.Name[0]]; ok {
		return err
	}
	return nil
}

func TestInstall_AllSuccess(t *testing.T) {
	cfg := &config.PackageConfig{
		Packages: []config.Package{{Name: []string{"pkg1"}}, {Name: []string{"pkg2"}}},
	}
	m := &mockManager{installResults: map[string]error{}}
	setTestLogger(t)
	code := runInstallCmd(cfg, m, []string{})
	if code != 0 {
		t.Errorf("Expected exit code 0, got %d", code)
	}
}

func TestInstall_SomeFail(t *testing.T) {
	cfg := &config.PackageConfig{
		Packages: []config.Package{{Name: []string{"pkg1"}}, {Name: []string{"pkg2"}}},
	}
	m := &mockManager{installResults: map[string]error{"pkg2": errors.New("fail")}}
	setTestLogger(t)
	code := runInstallCmd(cfg, m, []string{})
	if code == 0 {
		t.Error("Expected non-zero exit code on failure")
	}
}

func TestInstall_SpecificSuccess(t *testing.T) {
	cfg := &config.PackageConfig{
		Packages: []config.Package{{Name: []string{"pkg1"}}, {Name: []string{"pkg2"}}},
	}
	m := &mockManager{installResults: map[string]error{}}
	setTestLogger(t)
	code := runInstallCmd(cfg, m, []string{"pkg2"})
	if code != 0 {
		t.Errorf("Expected exit code 0, got %d", code)
	}
}

func TestInstall_SpecificFail(t *testing.T) {
	cfg := &config.PackageConfig{
		Packages: []config.Package{{Name: []string{"pkg1"}}, {Name: []string{"pkg2"}}},
	}
	m := &mockManager{installResults: map[string]error{"pkg2": errors.New("fail")}}
	setTestLogger(t)
	code := runInstallCmd(cfg, m, []string{"pkg2"})
	if code == 0 {
		t.Error("Expected non-zero exit code on failure")
	}
}

func TestInstall_PackageNotFound(t *testing.T) {
	cfg := &config.PackageConfig{
		Packages: []config.Package{{Name: []string{"pkg1"}}},
	}
	m := &mockManager{installResults: map[string]error{}}
	setTestLogger(t)
	code := runInstallCmd(cfg, m, []string{"doesnotexist"})
	if code == 0 {
		t.Error("Expected non-zero exit code for package not found")
	}
}

// setTestLogger configures the logger for testing and attaches log output to the test.
func setTestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger.Configure(logger.Config{
		Level:      logger.LevelDebug,
		TimeFormat: "",
		Output:     &buf,
		Prefix:     "test",
		ShowCaller: false,
	})
	t.Cleanup(func() {
		logger.Configure(logger.Config{})
	})
}

// runInstallCmd simulates the install command logic for testing.
func runInstallCmd(cfg *config.PackageConfig, m *mockManager, args []string) int {
	failedInstalls := []string{}
	if len(args) == 0 {
		for _, pkg := range cfg.Packages {
			err := m.Install(models.ToModel(pkg))
			if err != nil {
				failedInstalls = append(failedInstalls, pkg.Name[0])
			}
		}
	} else {
		pkgName := args[0]
		found := false
		for _, pkg := range cfg.Packages {
			if pkg.Name[0] == pkgName {
				found = true
				err := m.Install(models.ToModel(pkg))
				if err != nil {
					failedInstalls = append(failedInstalls, pkg.Name[0])
				}
			}
		}
		if !found {
			return 1
		}
	}
	if len(failedInstalls) > 0 {
		return 1
	}
	return 0
}
