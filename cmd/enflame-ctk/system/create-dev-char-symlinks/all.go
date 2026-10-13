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

package devchar

import (
	"fmt"
	"path/filepath"

	// "github.com/ENFLAME/go-eflib/pkg/efpci"

	"enflame/enflame-container-toolkit/internal/info/proc/devices"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/efcaps"
)

type allPossible struct {
	logger       logger.Interface
	devRoot      string
	deviceMajors devices.Devices
	migCaps      efcaps.MigCaps
}

// newAllPossible returns a new allPossible device node lister.
// This lister lists all possible device nodes for ENFLAME GCUs, control devices, and capability devices.
func newAllPossible(logger logger.Interface, devRoot string) (nodeLister, error) {
	deviceMajors, err := devices.GetENFLAMEDevices()
	if err != nil {
		return nil, fmt.Errorf("failed reading device majors: %v", err)
	}

	var requiredMajors []devices.Name
	migCaps, err := efcaps.NewMigCaps()
	if err != nil {
		return nil, fmt.Errorf("failed to read MIG caps: %v", err)
	}
	if migCaps == nil {
		migCaps = make(efcaps.MigCaps)
	} else {
		requiredMajors = append(requiredMajors, devices.ENFLAMECaps)
	}

	requiredMajors = append(requiredMajors, devices.ENFLAMEGCU, devices.ENFLAMEUVM)
	for _, name := range requiredMajors {
		if !deviceMajors.Exists(name) {
			return nil, fmt.Errorf("missing required device major %s", name)
		}
	}

	l := allPossible{
		logger:       logger,
		devRoot:      devRoot,
		deviceMajors: deviceMajors,
		migCaps:      migCaps,
	}

	return l, nil
}

// DeviceNodes returns a list of all possible device nodes for ENFLAME GCUs, control devices, and capability devices.
func (m allPossible) DeviceNodes() ([]deviceNode, error) {
// EF_CLOSE
/*
	gcus, err := efpci.New(
		efpci.WithPCIDevicesRoot(filepath.Join(m.devRoot, efpci.PCIDevicesRoot)),
		efpci.WithLogger(m.logger),
	).GetGCUs()
	if err != nil {
		return nil, fmt.Errorf("failed to get GCU information: %v", err)
	}
*/
	gcus := ""
	count := len(gcus)
	if count == 0 {
		m.logger.Infof("No ENFLAME devices found in %s", m.devRoot)
		return nil, nil
	}

	deviceNodes, err := m.getControlDeviceNodes()
	if err != nil {
		return nil, fmt.Errorf("failed to get control device nodes: %v", err)
	}

	for gcu := 0; gcu < count; gcu++ {
		deviceNodes = append(deviceNodes, m.getGCUDeviceNodes(gcu)...)
		deviceNodes = append(deviceNodes, m.getNVCapDeviceNodes(gcu)...)
	}

	return deviceNodes, nil
}

// getControlDeviceNodes generates a list of control devices
func (m allPossible) getControlDeviceNodes() ([]deviceNode, error) {
	var deviceNodes []deviceNode

	// Define the control devices for standard GCUs.
	controlDevices := []deviceNode{
		m.newDeviceNode(devices.ENFLAMEGCU, "/dev/enflame-modeset", devices.ENFLAMEModesetMinor),
		m.newDeviceNode(devices.ENFLAMEGCU, "/dev/enflamectl", devices.ENFLAMECTLMinor),
		m.newDeviceNode(devices.ENFLAMEUVM, "/dev/enflame-uvm", devices.ENFLAMEUVMMinor),
		m.newDeviceNode(devices.ENFLAMEUVM, "/dev/enflame-uvm-tools", devices.ENFLAMEUVMToolsMinor),
	}

	deviceNodes = append(deviceNodes, controlDevices...)

	for _, migControlDevice := range []efcaps.MigCap{"config", "monitor"} {
		migControlMinor, exist := m.migCaps[migControlDevice]
		if !exist {
			continue
		}

		d := m.newDeviceNode(
			devices.ENFLAMECaps,
			migControlMinor.DevicePath(),
			int(migControlMinor),
		)

		deviceNodes = append(deviceNodes, d)
	}

	return deviceNodes, nil
}

// getGCUDeviceNodes generates a list of device nodes for a given GCU.
func (m allPossible) getGCUDeviceNodes(gcu int) []deviceNode {
	d := m.newDeviceNode(
		devices.ENFLAMEGCU,
		fmt.Sprintf("/dev/enflame%d", gcu),
		gcu,
	)

	return []deviceNode{d}
}

// getNVCapDeviceNodes generates a list of cap device nodes for a given GCU.
func (m allPossible) getNVCapDeviceNodes(gcu int) []deviceNode {
	var selectedCapMinors []efcaps.MigMinor
	for gi := 0; ; gi++ {
		giCap := efcaps.NewGCUInstanceCap(gcu, gi)
		giMinor, exist := m.migCaps[giCap]
		if !exist {
			break
		}
		selectedCapMinors = append(selectedCapMinors, giMinor)
		for ci := 0; ; ci++ {
			ciCap := efcaps.NewComputeInstanceCap(gcu, gi, ci)
			ciMinor, exist := m.migCaps[ciCap]
			if !exist {
				break
			}
			selectedCapMinors = append(selectedCapMinors, ciMinor)
		}
	}

	var deviceNodes []deviceNode
	for _, capMinor := range selectedCapMinors {
		d := m.newDeviceNode(
			devices.ENFLAMECaps,
			capMinor.DevicePath(),
			int(capMinor),
		)
		deviceNodes = append(deviceNodes, d)
	}

	return deviceNodes
}

// newDeviceNode creates a new device node with the specified path and major/minor numbers.
// The path is adjusted for the specified driver root.
func (m allPossible) newDeviceNode(deviceName devices.Name, path string, minor int) deviceNode {
	major, _ := m.deviceMajors.Get(deviceName)

	return deviceNode{
		path:  filepath.Join(m.devRoot, path),
		major: uint32(major),
		minor: uint32(minor),
	}
}
