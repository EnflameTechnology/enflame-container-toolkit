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
	"path/filepath"
	"strings"

	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"

	"enflame/enflame-container-toolkit/internal/discover"
	"enflame/enflame-container-toolkit/internal/edits"
	"enflame/enflame-container-toolkit/internal/lookup/tops"
	"enflame/enflame-container-toolkit/pkg/efcdi/spec"
)

type managementlib efcdilib

var _ Interface = (*managementlib)(nil)

// GetAllDeviceSpecs returns all device specs for use in managemnt containers.
// A single device with the name `all` is returned.
func (m *managementlib) GetAllDeviceSpecs() ([]specs.Device, error) {
	devices, err := m.newManagementDeviceDiscoverer()
	if err != nil {
		return nil, fmt.Errorf("failed to create device discoverer: %v", err)
	}

	edits, err := edits.FromDiscoverer(devices)
	if err != nil {
		return nil, fmt.Errorf("failed to create edits from discoverer: %v", err)
	}

	if len(edits.DeviceNodes) == 0 {
		return nil, fmt.Errorf("no ENFLAME device nodes found")
	}

	device := specs.Device{
		Name:           "all",
		ContainerEdits: *edits.ContainerEdits,
	}
	return []specs.Device{device}, nil
}

// GetCommonEdits returns the common edits for use in managementlib containers.
func (m *managementlib) GetCommonEdits() (*cdi.ContainerEdits, error) {
	version, err := m.getTopsVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get TOPS version: %v", err)
	}

	driver, err := newDriverVersionDiscoverer(m.logger, m.driver, m.enflameCTKPath, m.ldconfigPath, version)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver library discoverer: %v", err)
	}

	edits, err := edits.FromDiscoverer(driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create edits from discoverer: %v", err)
	}

	return edits, nil
}

// getTopsVersion returns the TOPS version for use in managementlib containers.
func (m *managementlib) getTopsVersion() (string, error) {
	version, err := (*efcdilib)(m).getTopsVersion()
	if err == nil {
		return version, nil
	}

	libTopsPaths, err := tops.New(
		m.driver.Libraries(),
	).Locate(".*.*")
	if err != nil {
		return "", fmt.Errorf("failed to locate libtops.so: %v", err)
	}

	libTopsPath := libTopsPaths[0]

	version = strings.TrimPrefix(filepath.Base(libTopsPath), "libtops.so.")

	return version, nil
}

type managementDiscoverer struct {
	discover.Discover
}

// newManagementDeviceDiscoverer returns a discover.Discover that discovers device nodes for use in managementlib containers.
// EFML is not used to query devices and all device nodes are returned.
func (m *managementlib) newManagementDeviceDiscoverer() (discover.Discover, error) {
	deviceNodes := discover.NewCharDeviceDiscoverer(
		m.logger,
		m.devRoot,
		[]string{
			"/dev/enflame*",
			"/dev/enflame-caps/enflame-cap*",
			"/dev/enflame-modeset",
			"/dev/enflame-uvm-tools",
			"/dev/enflame-uvm",
			"/dev/enflamectl",
		},
	)

	deviceFolderPermissionHooks := newDeviceFolderPermissionHookDiscoverer(
		m.logger,
		m.devRoot,
		m.enflameCTKPath,
		deviceNodes,
	)

	d := discover.Merge(
		&managementDiscoverer{deviceNodes},
		deviceFolderPermissionHooks,
	)
	return d, nil
}

func (m *managementDiscoverer) Devices() ([]discover.Device, error) {
	devices, err := m.Discover.Devices()
	if err != nil {
		return devices, err
	}

	var filteredDevices []discover.Device
	for _, device := range devices {
		if m.nodeIsBlocked(device.HostPath) {
			continue
		}
		filteredDevices = append(filteredDevices, device)
	}

	return filteredDevices, nil
}

// nodeIsBlocked returns true if the specified device node should be ignored.
func (m managementDiscoverer) nodeIsBlocked(path string) bool {
	blockedPrefixes := []string{"enflame-fs", "enflame-efswitch", "enflame-eflink"}
	nodeName := filepath.Base(path)
	for _, prefix := range blockedPrefixes {
		if strings.HasPrefix(nodeName, prefix) {
			return true
		}
	}
	return false
}

// GetSpec is unsppported for the managementlib specs.
// managementlib is typically wrapped by a spec that implements GetSpec.
func (m *managementlib) GetSpec() (spec.Interface, error) {
	return nil, fmt.Errorf("GetSpec is not supported")
}

// GetGCUDeviceEdits is unsupported for the managementlib specs
func (m *managementlib) GetGCUDeviceEdits(device.Device) (*cdi.ContainerEdits, error) {
	return nil, fmt.Errorf("GetGCUDeviceEdits is not supported")
}

// GetGCUDeviceSpecs is unsupported for the managementlib specs
func (m *managementlib) GetGCUDeviceSpecs(int, device.Device) ([]specs.Device, error) {
	return nil, fmt.Errorf("GetGCUDeviceSpecs is not supported")
}

// GetMIGDeviceEdits is unsupported for the managementlib specs
func (m *managementlib) GetMIGDeviceEdits(device.Device, device.MigDevice) (*cdi.ContainerEdits, error) {
	return nil, fmt.Errorf("GetMIGDeviceEdits is not supported")
}

// GetMIGDeviceSpecs is unsupported for the managementlib specs
func (m *managementlib) GetMIGDeviceSpecs(int, device.Device, int, device.MigDevice) ([]specs.Device, error) {
	return nil, fmt.Errorf("GetMIGDeviceSpecs is not supported")
}

// GetDeviceSpecsByID returns the CDI device specs for the GCU(s) represented by
// the provided identifiers, where an identifier is an index or UUID of a valid
// GCU device.
func (l *managementlib) GetDeviceSpecsByID(...string) ([]specs.Device, error) {
	return nil, fmt.Errorf("GetDeviceSpecsByID is not supported")
}
*/