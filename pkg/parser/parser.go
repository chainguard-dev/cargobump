package parser

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/chainguard-dev/cargobump/pkg/types"
	"github.com/samber/lo"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

type cargoLockPackage struct {
	Name         string   `toml:"name"`
	Version      string   `toml:"version"`
	Dependencies []string `toml:"dependencies,omitempty"`
}

// Cargo Lock file based on https://doc.rust-lang.org/cargo/reference/manifest.html
type Lockfile struct {
	Packages []cargoLockPackage `toml:"package"`
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Used to marshal from yaml/json file to get the list of packages
type PackageList struct {
	Packages []types.Package `json:"packages" yaml:"packages"`
}

func (pa *Parser) ParseBumpFile(r io.Reader) (map[string]*types.Package, error) {
	bytes, _ := io.ReadAll(r)
	patches := map[string]*types.Package{}
	var packageList PackageList
	if err := yaml.Unmarshal(bytes, &packageList); err != nil {
		return patches, fmt.Errorf("unmarshaling file: %w", err)
	}
	for i, p := range packageList.Packages {
		if p.Name == "" {
			return patches, fmt.Errorf("invalid package spec at [%d], missing name", i)
		}
		if p.Version == "" {
			return patches, fmt.Errorf("invalid package spec at [%d], missing version", i)
		}
		if patches == nil {
			patches = make(map[string]*types.Package, 1)
		}
		patches[p.Name] = &packageList.Packages[i]
	}
	return patches, nil
}

func (pa *Parser) ParseCargoLock(r io.Reader) ([]types.CargoPackage, error) {
	var lockfile Lockfile
	decoder := toml.NewDecoder(r)
	if _, err := decoder.Decode(&lockfile); err != nil {
		return nil, xerrors.Errorf("decode error: %w", err)
	}

	//if _, err := r.Seek(0, io.SeekStart); err != nil {
	//	return nil, nil, xerrors.Errorf("seek error: %w", err)
	//}

	// We need to get version for unique dependencies for lockfile v3 from lockfile.Packages
	pkgs := lo.SliceToMap(lockfile.Packages, func(pkg cargoLockPackage) (string, cargoLockPackage) {
		return pkg.Name, pkg
	})

	var ps []types.CargoPackage
	for _, pkg := range lockfile.Packages {
		pkgID := packageID(pkg.Name, pkg.Version)
		p := types.CargoPackage{
			Name:    pkg.Name,
			Version: pkg.Version,
		}

		deps := parseDependencies(pkgID, pkg, pkgs)

		if len(deps) > 0 {
			p.Dependencies = append(p.Dependencies, deps...)
		}
		ps = append(ps, p)
	}
	sortCargoPkgs(ps)
	return ps, nil
}

func sortCargoPkgs(pkgs []types.CargoPackage) {
	sort.Slice(pkgs, func(i, j int) bool {
		return strings.Compare(packageID(pkgs[i].Name, pkgs[i].Version), packageID(pkgs[j].Name, pkgs[j].Version)) < 0
	})
}

func parseDependencies(pkgId string, pkg cargoLockPackage, pkgs map[string]cargoLockPackage) []string {
	var dependOn []string

	for _, pkgDep := range pkg.Dependencies {
		/*
			Dependency entries look like:
					"any-package" - if lock file contains only 1 version of dependency
					"any-package 0.1.2" if lock file contains more than 1 version of dependency
		*/
		fields := strings.Fields(pkgDep)
		switch len(fields) {
		case 1:
			name := fields[0]
			version, ok := pkgs[name]
			if !ok {
				log.Printf("can't find version for %s", name)
				continue
			}
			dependOn = append(dependOn, packageID(name, version.Version))
		// 2: non-unique dependency in new lock file
		// 3: old lock file
		case 2, 3:
			dependOn = append(dependOn, packageID(fields[0], fields[1]))
		default:
			log.Printf("wrong dependency format for %s", pkgDep)
			continue
		}
	}
	return dependOn
}

func packageID(name, version string) string {
	return fmt.Sprintf("%s@%s", name, version)
}
