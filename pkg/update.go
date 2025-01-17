/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package pkg

import (
	"fmt"
	"log"

	"golang.org/x/mod/semver"

	"github.com/chainguard-dev/cargobump/pkg/run"
	"github.com/chainguard-dev/cargobump/pkg/types"
)

func Update(patches map[string]*types.Package, pkgs []types.CargoPackage, cargoRoot string, update bool) error {
	// Run 'cargo update' prior upgrading any dependency
	if update {
		log.Printf("Running 'cargo update'...\n")
		if output, err := run.CargoUpdate(cargoRoot); err != nil {
			return fmt.Errorf("failed to run 'cargo update' '%v' with error: '%w'", output, err)
		}
	}
	for _, p := range pkgs {
		v, exists := patches[p.Name]
		if exists {
			if semver.Compare("v"+p.Version, "v"+patches[p.Name].Version) > 0 {
				log.Printf("package %s with version '%s' is already at version %s, skipping it!", p.Name, v.Version, p.Version)
				continue
			}
			log.Printf("Update package %s from version %s to %s\n", p.Name, p.Version, v.Version)

			if output, err := run.CargoUpdatePackage(p.Name, p.Version, v.Version, cargoRoot); err != nil {
				return fmt.Errorf("failed to run cargo update '%v' for package '%s' and version '%s' with error: '%w'", output, p.Name, p.Version, err)
			}

			log.Printf("Package updated successfully: %s to version %s\n", p.Name, v.Version)
		}

		// Try updating packages referring to a specific version
		packageVersion := p.Name + "@" + p.Version
		v, existsVersionRef := patches[packageVersion]
		if existsVersionRef {
			log.Printf("Update package with a specific version: %s\n", packageVersion)

			if semver.Compare("v"+p.Version, "v"+patches[packageVersion].Version) > 0 {
				return fmt.Errorf("warning: package %s with version '%s' is already at version %s", packageVersion, v.Version, p.Version)
			}

			if output, err := run.CargoUpdatePackage(p.Name, p.Version, v.Version, cargoRoot); err != nil {
				return fmt.Errorf("failed to run cargo update '%v' with error: '%w'", output, err)
			}

			log.Printf("Package updated successfully: %s to version %s\n", packageVersion, v.Version)
		}
	}

	return nil
}
