package pkg

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/chainguard-dev/cargobump/pkg/parser"
	"github.com/chainguard-dev/cargobump/pkg/types"
)

func TestUpdate(t *testing.T) {
	var file *os.File
	parse := parser.NewParser()
	var patches map[string]*types.Package

	bumpFile := "testdata/cargobump-deps.yaml"
	var err error
	file, err = os.Open(bumpFile)
	if err != nil {
		t.Fatalf("failed reading file: %v", err)
	}
	defer file.Close()

	patches, err = parse.ParseBumpFile(file)
	if err != nil {
		t.Fatalf("failed to parse the bump file: %v", err)
	}
	cargoRoot := "testdata"
	cargoLockFile, err := os.Open(filepath.Join(cargoRoot, "Cargo.lock.orig"))
	if err != nil {
		t.Fatalf("failed reading file: %v", err)
	}
	defer cargoLockFile.Close()
	cargoTomlFile, err := os.Open(filepath.Join(cargoRoot, "Cargo.toml.orig"))
	if err != nil {
		t.Fatalf("failed reading file: %v", err)
	}
	defer cargoTomlFile.Close()

	// copy the source file to the destination
	copyFile(t, cargoTomlFile.Name(), "testdata/Cargo.toml")
	defer os.Remove("testdata/Cargo.toml")

	// copy the source file to the destination
	copyFile(t, cargoLockFile.Name(), "testdata/Cargo.lock")
	defer os.Remove("testdata/Cargo.lock")

	pkgs, err := parse.ParseCargoLock(cargoLockFile)
	if err != nil {
		t.Fatalf("failed to parse Cargo.lock file: %v", err)
	}

	if err = Update(patches, pkgs, cargoRoot, false); err != nil {
		t.Errorf("failed to update packages: %v", err)
	}

	cargoLockFile, err = os.Open(filepath.Join(cargoRoot, "Cargo.lock"))
	if err != nil {
		t.Fatalf("failed reading file: %v", err)
	}
	defer cargoLockFile.Close()

	updatedPkgs, err := parse.ParseCargoLock(cargoLockFile)
	if err != nil {
		t.Fatalf("failed to parse Cargo.lock file: %v", err)
	}

	expectedPkgName := "mio"
	expectedPkgVersion := "0.8.11"
	found := false
	for _, pkg := range updatedPkgs {
		if pkg.Name == expectedPkgName && pkg.Version == expectedPkgVersion {
			found = true
		}
	}
	if !found {
		t.Errorf("failed to find '%s' updated package with version '%s'", expectedPkgName, expectedPkgVersion)
	}
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	_, err := exec.Command("cp", "-r", src, dst).Output()
	if err != nil {
		t.Fatal(err)
	}
}
