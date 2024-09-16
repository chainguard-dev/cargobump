package run

import (
	"os/exec"
	"strings"
)

func CargoUpdatePackage(name, version, cargoRoot string) (string, error) {
	cmd := exec.Command("cargo", "update", "--precise", version, "--package", name) //nolint:gosec
	cmd.Dir = cargoRoot
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(bytes)), err
	}
	return "", nil
}
