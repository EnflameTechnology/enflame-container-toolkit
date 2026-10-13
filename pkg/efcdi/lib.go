/**
# Copyright (c) NVIDIA CORPORATION.  All rights reserved.
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

package efcdi

import (
	"fmt"

	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	// "github.com/ENFLAME/go-eflib/pkg/eflib/info"
	// "github.com/ENFLAME/go-eflib/pkg/efml"
	"tags.cncf.io/container-device-interface/pkg/cdi"

	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/lookup/root"
	// "enflame/enflame-container-toolkit/internal/platform-support/tegra/csv"
	"enflame/enflame-container-toolkit/pkg/efcdi/spec"
	"enflame/enflame-container-toolkit/pkg/efcdi/transform"
)

type wrapper struct {
	Interface

	vendor string
	class  string

	mergedDeviceOptions []transform.MergedDeviceOption
}

type efcdilib struct {
	logger             logger.Interface
	// EF_CLOSE
	// efmllib            efml.Interface
	mode               string
	// EF_CLOSE
	// devicelib          device.Interface
	deviceNamers       DeviceNamers
	driverRoot         string
	devRoot            string
	enflameCTKPath      string
	ldconfigPath       string
	configSearchPaths  []string
	librarySearchPaths []string

	csvFiles          []string
	csvIgnorePatterns []string

	vendor string
	class  string

	driver  *root.Driver
	// EF_CLOSE
	// infolib info.Interface

	mergedDeviceOptions []transform.MergedDeviceOption
}

// New creates a new efcdi library
func New(opts ...Option) (Interface, error) {
	l := &efcdilib{}
	for _, opt := range opts {
		opt(l)
	}
	if l.mode == "" {
		l.mode = ModeAuto
	}
	if l.logger == nil {
		l.logger = logger.New()
	}
	if len(l.deviceNamers) == 0 {
		indexNamer, _ := NewDeviceNamer(DeviceNameStrategyIndex)
		l.deviceNamers = []DeviceNamer{indexNamer}
	}
	if l.driverRoot == "" {
		l.driverRoot = "/"
	}
	if l.devRoot == "" {
		l.devRoot = l.driverRoot
	}
	if l.enflameCTKPath == "" {
		l.enflameCTKPath = "/usr/bin/enflame-ctk"
	}

// EF_CLOSE
/*
	if l.infolib == nil {
		l.infolib = info.New()
	}
*/
	l.driver = root.New(
		root.WithLogger(l.logger),
		root.WithDriverRoot(l.driverRoot),
		root.WithLibrarySearchPaths(l.librarySearchPaths...),
	)

	var lib Interface
	switch l.resolveMode() {
// EF_CLOSE
/*
	case ModeCSV:
		if len(l.csvFiles) == 0 {
			l.csvFiles = csv.DefaultFileList()
		}
		lib = (*csvlib)(l)
	case ModeManagement:
		if l.vendor == "" {
			l.vendor = "management.enflame.com"
		}
		lib = (*managementlib)(l)

	case ModeEfml:
		if l.efmllib == nil {
			l.efmllib = efml.New()
		}
		if l.devicelib == nil {
			l.devicelib = device.New(device.WithEfml(l.efmllib))
		}

		lib = (*efmllib)(l)

	case ModeWsl:
		lib = (*wsllib)(l)
	case ModeGds:
		if l.class == "" {
			l.class = "gds"
		}
		lib = (*gdslib)(l)
*/
	case ModeMofed:
		if l.class == "" {
			l.class = "mofed"
		}
		lib = (*mofedlib)(l)
	default:
		return nil, fmt.Errorf("unknown mode %q", l.mode)
	}

	w := wrapper{
		Interface:           lib,
		vendor:              l.vendor,
		class:               l.class,
		mergedDeviceOptions: l.mergedDeviceOptions,
	}
	return &w, nil
}

// GetSpec combines the device specs and common edits from the wrapped Interface to a single spec.Interface.
func (l *wrapper) GetSpec() (spec.Interface, error) {
// EF_CLOSE
/*
	deviceSpecs, err := l.GetAllDeviceSpecs()
	if err != nil {
		return nil, err
	}
*/
	edits, err := l.GetCommonEdits()
	if err != nil {
		return nil, err
	}

	return spec.New(
		// EF_CLOSE
		// spec.WithDeviceSpecs(deviceSpecs),
		spec.WithEdits(*edits.ContainerEdits),
		spec.WithVendor(l.vendor),
		spec.WithClass(l.class),
		spec.WithMergedDeviceOptions(l.mergedDeviceOptions...),
	)
}

// GetCommonEdits returns the wrapped edits and adds additional edits on top.
func (m *wrapper) GetCommonEdits() (*cdi.ContainerEdits, error) {
	edits, err := m.Interface.GetCommonEdits()
	if err != nil {
		return nil, err
	}
	edits.Env = append(edits.Env, "ENFLAME_VISIBLE_DEVICES=void")

	return edits, nil
}

// resolveMode resolves the mode for CDI spec generation based on the current system.
func (l *efcdilib) resolveMode() (rmode string) {
	if l.mode != ModeAuto {
		return l.mode
	}
	defer func() {
		l.logger.Infof("Auto-detected mode as %q", rmode)
	}()

// EF_CLOSE
/*
	isWSL, reason := l.infolib.HasDXCore()
	l.logger.Debugf("Is WSL-based system? %v: %v", isWSL, reason)

	if isWSL {
		return ModeWsl
	}

	isEfml, reason := l.infolib.HasEfml()
	l.logger.Debugf("Is EFML-based system? %v: %v", isEfml, reason)

	isTegra, reason := l.infolib.IsTegraSystem()
	l.logger.Debugf("Is Tegra-based system? %v: %v", isTegra, reason)

	if isTegra && !isEfml {
		return ModeCSV
	}
*/
	return ModeEfml
}

// getTopsVersion returns the TOPS version of the current system.
func (l *efcdilib) getTopsVersion() (string, error) {
	// EF_CLOSE
	/*
	if hasEFML, reason := l.infolib.HasEfml(); !hasEFML {
		return "", fmt.Errorf("efml not detected: %v", reason)
	}
	if l.efmllib == nil {
		return "", fmt.Errorf("efml library not initialized")
	}
	r := l.efmllib.Init()
	if r != efml.SUCCESS {
		return "", fmt.Errorf("failed to initialize efml: %v", r)
	}
	defer func() {
		if r := l.efmllib.Shutdown(); r != efml.SUCCESS {
			l.logger.Warningf("failed to shutdown EFML: %v", r)
		}
	}()

	version, r := l.efmllib.SystemGetDriverVersion()
	if r != efml.SUCCESS {
		return "", fmt.Errorf("failed to get driver version: %v", r)
	}
	return version, nil
	*/
	return "N/A", nil
}
