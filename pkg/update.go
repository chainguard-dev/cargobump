package pkg

import (
	"fmt"
	"log"

	"github.com/hectorj2f/cargobump/pkg/run"
	"github.com/hectorj2f/cargobump/pkg/types"
	"golang.org/x/mod/semver"
)

func Update(patches map[string]*types.Package, pkgs []types.CargoPackage, cargoRoot string) error {
	for _, p := range pkgs {
		v, exists := patches[p.Name]
		if exists {
			log.Printf("Update package: %s\n", p.Name)
			if semver.Compare(p.Version, patches[p.Name].Version) > 0 {
				return fmt.Errorf("package %s with version '%s' is already at version %s", p.Name, v.Version, p.Version)
			}
			if output, err := run.CargoUpdatePackage(p.Name, v.Version, cargoRoot); err != nil {
				return fmt.Errorf("failed to run cargo update '%v' with error: '%v'", output, err)
			}
			log.Printf("Package updated successfully: %s to version %s\n", p.Name, v.Version)
		}
	}
	return nil
}
