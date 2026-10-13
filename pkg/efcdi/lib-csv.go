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
//EF_CLOSE
/*
import (
	"fmt"

	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"

	"enflame/enflame-container-toolkit/internal/discover"
	"enflame/enflame-container-toolkit/internal/edits"
	"enflame/enflame-container-toolkit/internal/platform-support/tegra"
	"enflame/enflame-container-toolkit/pkg/efcdi/spec"
)

type csvlib efcdilib

var _ Interface = (*csvlib)(nil)

// GetSpec should not be called for wsllib
func (l *csvlib) GetSpec() (spec.Interface, error) {
	return nil, fmt.Errorf("Unexpected call to csvlib.GetSpec()")
}

// GetAllDeviceSpecs returns the device specs for all available devices.
func (l *csvlib) GetAllDeviceSpecs() ([]specs.Device, error) {
	d, err := tegra.New(
		tegra.WithLogger(l.logger),
		tegra.WithDriverRoot(l.driverRoot),
		tegra.WithDevRoot(l.devRoot),
		tegra.WithENFLAMECTKPath(l.enflameCTKPath),
		tegra.WithLdconfigPath(l.ldconfigPath),
		tegra.WithCSVFiles(l.csvFiles),
		tegra.WithLibrarySearchPaths(l.librarySearchPaths...),
		tegra.WithIngorePatterns(l.csvIgnorePatterns...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create discoverer for CSV files: %v", err)
	}
	e, err := edits.FromDiscoverer(d)
	if err != nil {
		return nil, fmt.Errorf("failed to create container edits for CSV files: %v", err)
	}

	names, err := l.deviceNamers.GetDeviceNames(0, uuidIgnored{})
	if err != nil {
		return nil, fmt.Errorf("failed to get device name: %v", err)
	}
	var deviceSpecs []specs.Device
	for _, name := range names {
		deviceSpec := specs.Device{
			Name:           name,
			ContainerEdits: *e.ContainerEdits,
		}
		deviceSpecs = append(deviceSpecs, deviceSpec)
	}

	return deviceSpecs, nil
}

// GetCommonEdits generates a CDI specification that can be used for ANY devices
func (l *csvlib) GetCommonEdits() (*cdi.ContainerEdits, error) {
	d := discover.None{}
	return edits.FromDiscoverer(d)
}

// GetGCUDeviceEdits generates a CDI specification that can be used for GCU devices
func (l *csvlib) GetGCUDeviceEdits(device.Device) (*cdi.ContainerEdits, error) {
	return nil, fmt.Errorf("GetGCUDeviceEdits is not supported for CSV files")
}

// GetGCUDeviceSpecs returns the CDI device specs for the full GCU represented by 'device'.
func (l *csvlib) GetGCUDeviceSpecs(i int, d device.Device) ([]specs.Device, error) {
	return nil, fmt.Errorf("GetGCUDeviceSpecs is not supported for CSV files")
}

// GetMIGDeviceEdits generates a CDI specification that can be used for MIG devices
func (l *csvlib) GetMIGDeviceEdits(device.Device, device.MigDevice) (*cdi.ContainerEdits, error) {
	return nil, fmt.Errorf("GetMIGDeviceEdits is not supported for CSV files")
}

// GetMIGDeviceSpecs returns the CDI device specs for the full MIG represented by 'device'.
func (l *csvlib) GetMIGDeviceSpecs(int, device.Device, int, device.MigDevice) ([]specs.Device, error) {
	return nil, fmt.Errorf("GetMIGDeviceSpecs is not supported for CSV files")
}

// GetDeviceSpecsByID returns the CDI device specs for the GCU(s) represented by
// the provided identifiers, where an identifier is an index or UUID of a valid
// GCU device.
func (l *csvlib) GetDeviceSpecsByID(...string) ([]specs.Device, error) {
	return nil, fmt.Errorf("GetDeviceSpecsByID is not supported for CSV files")
}
*/