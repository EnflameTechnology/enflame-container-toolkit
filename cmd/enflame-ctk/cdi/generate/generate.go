/**
# Copyright (c) 2022, NVIDIA CORPORATION.  All rights reserved.
# Copyright (c) 2024 Enflame Technologies, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
**/
//EF_SUPPORT

package generate

import (
	"fmt"
	"os"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	cdi "tags.cncf.io/container-device-interface/pkg/parser"

	"enflame/enflame-container-toolkit/internal/config"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/platform-support/tegra/csv"
	"enflame/enflame-container-toolkit/pkg/efcdi"
	"enflame/enflame-container-toolkit/pkg/efcdi/spec"
	"enflame/enflame-container-toolkit/pkg/efcdi/transform"
	"tags.cncf.io/container-device-interface/specs-go"
)

const (
	allDeviceName = "all"
)

type command struct {
	logger logger.Interface
}

type options struct {
	output               string
	format               string
	deviceNameStrategies cli.StringSlice
	driverRoot           string
	devRoot              string
	enflameCTKPath        string
	ldconfigPath         string
	mode                 string
	vendor               string
	class                string

	configSearchPaths  cli.StringSlice
	librarySearchPaths cli.StringSlice

	csv struct {
		files          cli.StringSlice
		ignorePatterns cli.StringSlice
	}
}

// NewCommand constructs a generate-cdi command with the specified logger
func NewCommand(logger logger.Interface) *cli.Command {
	c := command{
		logger: logger,
	}
	return c.build()
}

// build creates the CLI command
func (m command) build() *cli.Command {
	opts := options{}

	// Create the 'generate-cdi' command
	c := cli.Command{
		Name:  "generate",
		Usage: "Generate CDI specifications for use with CDI-enabled runtimes",
		Before: func(c *cli.Context) error {
			return m.validateFlags(c, &opts)
		},
		Action: func(c *cli.Context) error {
			return m.run(c, &opts)
		},
	}

	c.Flags = []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "config-search-path",
			Usage:       "Specify the path to search for config files when discovering the entities that should be included in the CDI specification.",
			Destination: &opts.configSearchPaths,
		},
		&cli.StringFlag{
			Name:        "output",
			Usage:       "Specify the file to output the generated CDI specification to. If this is '' the specification is output to STDOUT",
			Destination: &opts.output,
		},
		&cli.StringFlag{
			Name:        "format",
			Usage:       "The output format for the generated spec [json | yaml]. This overrides the format defined by the output file extension (if specified).",
			Value:       spec.FormatYAML,
			Destination: &opts.format,
		},
		&cli.StringFlag{
			Name:        "mode",
			Aliases:     []string{"discovery-mode"},
			Usage:       "The mode to use when discovering the available entities. One of [auto | efml | wsl]. If mode is set to 'auto' the mode will be determined based on the system configuration.",
			Value:       efcdi.ModeAuto,
			Destination: &opts.mode,
		},
		&cli.StringFlag{
			Name:        "dev-root",
			Usage:       "Specify the root where `/dev` is located. If this is not specified, the driver-root is assumed.",
			Destination: &opts.devRoot,
		},
		&cli.StringSliceFlag{
			Name:        "device-name-strategy",
			Usage:       "Specify the strategy for generating device names. If this is specified multiple times, the devices will be duplicated for each strategy. One of [index | uuid | type-index]",
			Value:       cli.NewStringSlice(efcdi.DeviceNameStrategyIndex, efcdi.DeviceNameStrategyUUID),
			Destination: &opts.deviceNameStrategies,
		},
		&cli.StringFlag{
			Name:        "driver-root",
			Usage:       "Specify the ENFLAME GCU driver root to use when discovering the entities that should be included in the CDI specification.",
			Destination: &opts.driverRoot,
		},
		&cli.StringSliceFlag{
			Name:        "library-search-path",
			Usage:       "Specify the path to search for libraries when discovering the entities that should be included in the CDI specification.\n\tNote: This option only applies to CSV mode.",
			Destination: &opts.librarySearchPaths,
		},
		&cli.StringFlag{
			Name:        "enflame-ctk-path",
			Usage:       "Specify the path to use for the enflame-ctk in the generated CDI specification. If this is left empty, the path will be searched.",
			Destination: &opts.enflameCTKPath,
		},
		&cli.StringFlag{
			Name:        "ldconfig-path",
			Usage:       "Specify the path to use for ldconfig in the generated CDI specification",
			Destination: &opts.ldconfigPath,
		},
		&cli.StringFlag{
			Name:        "vendor",
			Aliases:     []string{"cdi-vendor"},
			Usage:       "the vendor string to use for the generated CDI specification.",
			Value:       "enflame.com",
			Destination: &opts.vendor,
		},
		&cli.StringFlag{
			Name:        "class",
			Aliases:     []string{"cdi-class"},
			Usage:       "the class string to use for the generated CDI specification.",
			Value:       "gcu",
			Destination: &opts.class,
		},
		&cli.StringSliceFlag{
			Name:        "csv.file",
			Usage:       "The path to the list of CSV files to use when generating the CDI specification in CSV mode.",
			Value:       cli.NewStringSlice(csv.DefaultFileList()...),
			Destination: &opts.csv.files,
		},
		&cli.StringSliceFlag{
			Name:        "csv.ignore-pattern",
			Usage:       "Specify a pattern the CSV mount specifications.",
			Destination: &opts.csv.ignorePatterns,
		},
	}

	return &c
}

func (m command) validateFlags(c *cli.Context, opts *options) error {
	opts.format = strings.ToLower(opts.format)
	switch opts.format {
	case spec.FormatJSON:
	case spec.FormatYAML:
	default:
		return fmt.Errorf("invalid output format: %v", opts.format)
	}

	opts.mode = strings.ToLower(opts.mode)
	switch opts.mode {
	case efcdi.ModeAuto:
	case efcdi.ModeCSV:
	case efcdi.ModeEfml:
	case efcdi.ModeWsl:
	case efcdi.ModeManagement:
	default:
		return fmt.Errorf("invalid discovery mode: %v", opts.mode)
	}

	for _, strategy := range opts.deviceNameStrategies.Value() {
		_, err := efcdi.NewDeviceNamer(strategy)
		if err != nil {
			return err
		}
	}

	opts.enflameCTKPath = config.ResolveENFLAMECTKPath(m.logger, opts.enflameCTKPath)

	if outputFileFormat := formatFromFilename(opts.output); outputFileFormat != "" {
		m.logger.Debugf("Inferred output format as %q from output file name", outputFileFormat)
		if !c.IsSet("format") {
			opts.format = outputFileFormat
		} else if outputFileFormat != opts.format {
			m.logger.Warningf("Requested output format %q does not match format implied by output file name: %q", opts.format, outputFileFormat)
		}
	}

	if err := cdi.ValidateVendorName(opts.vendor); err != nil {
		return fmt.Errorf("invalid CDI vendor name: %v", err)
	}
	if err := cdi.ValidateClassName(opts.class); err != nil {
		return fmt.Errorf("invalid CDI class name: %v", err)
	}
	return nil
}

func (m command) run(c *cli.Context, opts *options) error {
	spec, err := m.generateSpec(opts)
	if err != nil {
		return fmt.Errorf("failed to generate CDI spec: %v", err)
	}
	m.logger.Infof("Generated CDI spec with version %v", spec.Raw().Version)

	if opts.output == "" {
		_, err := spec.WriteTo(os.Stdout)
		if err != nil {
			return fmt.Errorf("failed to write CDI spec to STDOUT: %v", err)
		}
		return nil
	}

	return spec.Save(opts.output)
}

func formatFromFilename(filename string) string {
	ext := filepath.Ext(filename)
	switch strings.ToLower(ext) {
	case ".json":
		return spec.FormatJSON
	case ".yaml", ".yml":
		return spec.FormatYAML
	}

	return ""
}

func (m command) generateSpec(opts *options) (spec.Interface, error) {
	var deviceNamers []efcdi.DeviceNamer
	for _, strategy := range opts.deviceNameStrategies.Value() {
		deviceNamer, err := efcdi.NewDeviceNamer(strategy)
		if err != nil {
			return nil, fmt.Errorf("failed to create device namer: %v", err)
		}
		deviceNamers = append(deviceNamers, deviceNamer)
	}

// EF_CLOSE
/*
	cdilib, err := efcdi.New(
		efcdi.WithLogger(m.logger),
		efcdi.WithDriverRoot(opts.driverRoot),
		efcdi.WithDevRoot(opts.devRoot),
		efcdi.WithENFLAMECTKPath(opts.enflameCTKPath),
		efcdi.WithLdconfigPath(opts.ldconfigPath),
		efcdi.WithDeviceNamers(deviceNamers...),
		efcdi.WithMode(opts.mode),
		efcdi.WithConfigSearchPaths(opts.configSearchPaths.Value()),
		efcdi.WithLibrarySearchPaths(opts.librarySearchPaths.Value()),
		efcdi.WithCSVFiles(opts.csv.files.Value()),
		efcdi.WithCSVIgnorePatterns(opts.csv.ignorePatterns.Value()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create CDI library: %v", err)
	}

	deviceSpecs, err := cdilib.GetAllDeviceSpecs()
	if err != nil {
		return nil, fmt.Errorf("failed to create device CDI specs: %v", err)
	}

	commonEdits, err := cdilib.GetCommonEdits()
	if err != nil {
		return nil, fmt.Errorf("failed to create edits common for entities: %v", err)
	}
*/
	searchLibPaths := []string{"/usr/lib/", "/usr/lib64/", "/usr/lib/x86_64-linux-gnu/"}
	searchBinPaths := []string{"/usr/bin/", "/usr/local/bin/", "/bin", "/usr/sbin", "/sbin"}

	var mountFiles []string

	for _, path := range searchLibPaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		err = filepath.WalkDir(absPath, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
			   return err
			}
			if strings.Contains(d.Name(), "libefml") {
				mountFiles = append(mountFiles, filePath)
			}
			return nil
		})
	}

	for _, path := range searchBinPaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		err = filepath.WalkDir(absPath, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
			   return err
			}
			if strings.Contains(d.Name(), "efsmi"){
				mountFiles = append(mountFiles, filePath)
			}
			return nil
		})
	}

	mountFiles = append(mountFiles, "/lib/firmware/enflame_ucode_bin")


	devicesPath := "/dev"
	if opts.class == "" {
            opts.class = "gcu"
	}

	invalid_class := opts.class + "_"

	var devices []specs.Device

	filepath.WalkDir(devicesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasPrefix(d.Name(), opts.class) || strings.HasPrefix(d.Name(), invalid_class){
			return nil
		}

		if strings.HasPrefix(d.Name(), "gcuctl") {
			return nil
		}
		deviceName := strings.TrimPrefix(d.Name(), opts.class)

		device := specs.Device{
			Name: deviceName,
			ContainerEdits: specs.ContainerEdits{
				DeviceNodes: []*specs.DeviceNode{
				       {
					      Path: devicesPath + "/" + d.Name(),
				       },
				},
			},
		}
		devices = append(devices, device)

		return nil
	})
        containerEdits := specs.ContainerEdits{
		Mounts: []*specs.Mount{},
	}

	for _, file := range mountFiles {
		mount := &specs.Mount{
			HostPath:      file,
			ContainerPath: file,
			Options:       []string{"ro", "nosuid", "nodev", "bind"},
			Type:          "bind",
		}
		containerEdits.Mounts = append(containerEdits.Mounts, mount)
	}

	deviceNode := &specs.DeviceNode{
		Path: "/dev/gcuctl",
	}
	containerEdits.DeviceNodes = append(containerEdits.DeviceNodes, deviceNode)


	return spec.New(
		spec.WithVendor(opts.vendor),
		spec.WithClass(opts.class),
		// EF_CLOSE
		// spec.WithDeviceSpecs(deviceSpecs),
		// spec.WithEdits(*commonEdits.ContainerEdits),
		spec.WithDeviceSpecs(devices),
		spec.WithEdits(containerEdits),
		spec.WithFormat(opts.format),
		spec.WithMergedDeviceOptions(
			transform.WithName(allDeviceName),
			transform.WithSkipIfExists(true),
		),
		spec.WithPermissions(0644),
	)
}
