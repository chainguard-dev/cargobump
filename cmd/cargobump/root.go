/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package cargobump

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"chainguard.dev/apko/pkg/log"
	charmlog "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"sigs.k8s.io/release-utils/version"

	"github.com/chainguard-dev/cargobump/pkg"
	"github.com/chainguard-dev/cargobump/pkg/parser"
	"github.com/chainguard-dev/cargobump/pkg/types"
)

type rootCLIFlags struct {
	packages  string
	bumpFile  string
	cargoRoot string
	update    bool
}

var rootFlags rootCLIFlags

func New() *cobra.Command {
	var logPolicy []string
	var level log.CharmLogLevel

	cmd := &cobra.Command{
		Use:   "cargobump <file-to-bump>",
		Short: "cargobump cli",
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			out, err := log.Writer(logPolicy)
			if err != nil {
				return fmt.Errorf("failed to create log writer: %w", err)
			}

			slog.SetDefault(slog.New(charmlog.NewWithOptions(out, charmlog.Options{ReportTimestamp: true, Level: charmlog.Level(level)})))

			return nil
		},

		// Uncomment the following line if your bare application
		// has an action associated with it:
		RunE: func(_ *cobra.Command, _ []string) error {
			if rootFlags.packages == "" && rootFlags.bumpFile == "" {
				return fmt.Errorf("no packages or bump file provides, use --packages/--bump-file")
			}

			if rootFlags.bumpFile != "" && rootFlags.packages != "" {
				return fmt.Errorf("use either --packages or --bump-file")
			}
			var file *os.File
			parse := parser.NewParser()
			var patches map[string]*types.Package

			if rootFlags.bumpFile != "" {
				var err error
				file, err = os.Open(rootFlags.bumpFile)
				if err != nil {
					return fmt.Errorf("failed reading file: %w", err)
				}
				defer file.Close()

				patches, err = parse.ParseBumpFile(file)

				if err != nil {
					return fmt.Errorf("failed to parse the bump file: %w", err)
				}
			} else {
				var err error
				patches, err = parsePackageList(rootFlags.packages)
				if err != nil {
					return fmt.Errorf("failed to parse package list: %w", err)
				}
			}

			cargoLockFile, err := os.Open(filepath.Join(rootFlags.cargoRoot, "Cargo.lock"))
			if err != nil {
				return fmt.Errorf("failed reading file: %w", err)
			}
			defer cargoLockFile.Close()

			pkgs, err := parse.ParseCargoLock(cargoLockFile)
			if err != nil {
				return fmt.Errorf("failed to parse Cargo.lock file: %w", err)
			}

			if err = pkg.Update(patches, pkgs, rootFlags.cargoRoot, rootFlags.update); err != nil {
				return fmt.Errorf("failed to update packages: %w", err)
			}

			return nil
		},
	}
	cmd.PersistentFlags().StringSliceVar(&logPolicy, "log-policy", []string{"builtin:stderr"}, "log policy (e.g. builtin:stderr, /tmp/log/foo)")
	cmd.PersistentFlags().Var(&level, "log-level", "log level (e.g. debug, info, warn, error)")

	cmd.AddCommand(version.WithFont("starwars"))

	cmd.DisableAutoGenTag = true

	flagSet := cmd.Flags()
	flagSet.StringVar(&rootFlags.cargoRoot, "cargoroot", "", "path to the Cargo.lock root")
	flagSet.StringVar(&rootFlags.packages, "packages", "", "A space-separated list of dependencies to update in form package@version")
	flagSet.StringVar(&rootFlags.bumpFile, "bump-file", "", "The input file to read dependencies to bump from")
	flagSet.BoolVar(&rootFlags.update, "run-update", false, "Run 'cargo update' prior upgrading any dependency")

	return cmd
}

// parsePackageList converts the package list string into a map of package names to their versions.
func parsePackageList(pkgs string) (map[string]*types.Package, error) {
	patches := make(map[string]*types.Package)
	for _, pkg := range strings.Fields(pkgs) {
		parts := strings.Split(pkg, "@")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid package %q, want <name@version>", pkg)
		}
		if p, ok := patches[parts[0]]; ok {
			return nil, fmt.Errorf("duplicate package %s@%s found, already defined as %s@%s", parts[0], parts[1], p.Name, p.Version)
		}
		patches[parts[0]] = &types.Package{
			Name:    parts[0],
			Version: parts[1],
		}
	}
	return patches, nil
}
