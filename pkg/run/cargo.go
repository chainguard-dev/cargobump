/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

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

func CargoUpdate(cargoRoot string) (string, error) {
	cmd := exec.Command("cargo", "update") //nolint:gosec
	cmd.Dir = cargoRoot
	if bytes, err := cmd.CombinedOutput(); err != nil {
		return strings.TrimSpace(string(bytes)), err
	}
	return "", nil
}
