/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package run

import (
	"fmt"
	"os/exec"
	"strings"
)

func CargoUpdatePackage(name, oldVersion, newVersion, cargoRoot string) (string, error) {
	cmd := exec.Command("cargo", "update", "--precise", newVersion, "--package", fmt.Sprintf("%s@%s", name, oldVersion)) //nolint:gosec
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
