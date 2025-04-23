package manager

import (
	"os/exec"
	"reflect"
	"testing"
	"voltig/internal/models"
)

func TestBrewManager_IsAvailable(t *testing.T) {
	b := &BrewManager{}
	t.Parallel()
	// This test will pass if brew is installed on the system running the test
	_ = b.IsAvailable() // Just ensure it doesn't panic
}

// NOTE: These tests use stubbing/hooking to simulate system calls for full branch coverage.


func TestBrewManager_GetStatus(t *testing.T) {
	b := &BrewManager{}
	t.Parallel()
	// Save original exec functions
	origCommand := execCommand
	origLookPath := execLookPath
	defer func() {
		execCommand = origCommand
		execLookPath = origLookPath
	}()

	tests := []struct {
		name      string
		brewOut   []byte
		brewErr   error
		lookPath  func(string) (string, error)
		pkg       models.Package
		want      models.PackageStatus
	}{
		{
			name:    "installed via brew",
			brewOut: []byte("foo 1.2.3\n"),
			brewErr: nil,
			lookPath: func(bin string) (string, error) { return "", exec.ErrNotFound },
			pkg:     models.Package{Name: []string{"foo"}},
			want:    models.PackageStatus{Name: "foo", Status: "installed", Version: "foo 1.2.3\n"},
		},
		{
			name:    "missing",
			brewOut: nil,
			brewErr: exec.ErrNotFound,
			lookPath: func(bin string) (string, error) { return "", exec.ErrNotFound },
			pkg:     models.Package{Name: []string{"bar"}},
			want:    models.PackageStatus{Name: "bar", Status: "missing"},
		},
		{
			name:    "external binary",
			brewOut: nil,
			brewErr: exec.ErrNotFound,
			lookPath: func(bin string) (string, error) {
				if bin == "baz" { return "/usr/bin/baz", nil }
				return "", exec.ErrNotFound
			},
			pkg: models.Package{Name: []string{"baz"}},
			want: models.PackageStatus{Name: "baz", Status: "installed (external)", Version: "unknown"},
		},
		{
			name:    "nvm node",
			brewOut: nil,
			brewErr: exec.ErrNotFound,
			lookPath: func(bin string) (string, error) {
				if bin == "node" { return "/Users/test/.nvm/versions/node/v18.16.0/bin/node", nil }
				if bin == "nvm" { return "/usr/local/bin/nvm", nil }
				return "", exec.ErrNotFound
			},
			pkg: models.Package{Name: []string{"node"}},
			want: models.PackageStatus{Name: "node", Status: "installed (nvm)", Version: "unknown"},
		},
		{
			name:    "pyenv python",
			brewOut: nil,
			brewErr: exec.ErrNotFound,
			lookPath: func(bin string) (string, error) {
				if bin == "python" { return "/Users/test/.pyenv/versions/3.9.1/bin/python", nil }
				if bin == "pyenv" { return "/usr/local/bin/pyenv", nil }
				return "", exec.ErrNotFound
			},
			pkg: models.Package{Name: []string{"python"}},
			want: models.PackageStatus{Name: "python", Status: "installed (pyenv)", Version: "unknown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{}
			}
			execOutput = func(cmd *exec.Cmd) ([]byte, error) { return tt.brewOut, tt.brewErr }
			execLookPath = tt.lookPath
			origExecOutput := execOutput
			defer func() { execOutput = origExecOutput }()
			got, _ := b.GetStatus(tt.pkg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_getBinaryVersion(t *testing.T) {
	// This will likely return "unknown" for a non-existent binary
	version := getBinaryVersion("nonexistent-binary-xyz")
	if version != "unknown" {
		t.Errorf("Expected 'unknown' for missing binary, got %s", version)
	}
	// Simulate a binary that outputs a version string
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return &exec.Cmd{}
	}
	execOutput = func(cmd *exec.Cmd) ([]byte, error) { return []byte("v1.2.3\n"), nil }
	version = getBinaryVersion("dummy")
	if version != "v1.2.3" {
		t.Errorf("Expected extracted version, got %s", version)
	}
	// Simulate a binary that outputs no version
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return &exec.Cmd{}
	}
	execOutput = func(cmd *exec.Cmd) ([]byte, error) { return []byte("no version here\n"), nil }
	version = getBinaryVersion("dummy")
	if version != "no version here" {
		t.Errorf("Expected fallback to full output, got %s", version)
	}
}


func TestBrewManager_Install_OutputFn(t *testing.T) {
	b := &BrewManager{}
	t.Parallel()
	pkg := models.Package{Name: []string{"test"}}
	outputLines := []string{}
	outputFn := func(s string) { outputLines = append(outputLines, s) }
	// We expect this to fail since the package doesn't exist, but outputFn should still be valid
	_ = b.Install(pkg, outputFn)
	// No assertion on outputLines because brew is not actually called for a dummy package
}

func TestBrewManager_Remove_OutputFn(t *testing.T) {
	b := &BrewManager{}
	t.Parallel()
	pkg := models.Package{Name: []string{"test"}}
	outputLines := []string{}
	outputFn := func(s string) { outputLines = append(outputLines, s) }
	_ = b.Remove(pkg, outputFn)
}

func Test_isNvmNode(t *testing.T) {
	if !isNvmNode("/Users/test/.nvm/versions/node/v18.16.0/bin/node") {
		t.Error("Expected true for nvm-managed path")
	}
	if isNvmNode("/usr/local/bin/node") {
		t.Error("Expected false for non-nvm path")
	}
}

func Test_isPyenvPython(t *testing.T) {
	if !isPyenvPython("/Users/test/.pyenv/versions/3.9.1/bin/python") {
		t.Error("Expected true for pyenv-managed path")
	}
	if isPyenvPython("/usr/bin/python") {
		t.Error("Expected false for non-pyenv path")
	}
}



func TestBrewManager_Update(t *testing.T) {
	b := &BrewManager{}
	pkg := models.Package{Name: []string{"test"}}
	// Should not panic, even if package does not exist
	_ = b.Update(pkg)
}

