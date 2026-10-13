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
	// "github.com/ENFLAME/go-eflib/pkg/efml"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"

	"enflame/enflame-container-toolkit/internal/discover"
	"enflame/enflame-container-toolkit/internal/edits"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/efcaps"
)

// GetMIGDeviceSpecs returns the CDI device specs for the full GCU represented by 'device'.
func (l *efmllib) GetMIGDeviceSpecs(i int, d device.Device, j int, mig device.MigDevice) ([]specs.Device, error) {
	edits, err := l.GetMIGDeviceEdits(d, mig)
	if err != nil {
		return nil, fmt.Errorf("failed to get edits for device: %v", err)
	}

	names, err := l.deviceNamers.GetMigDeviceNames(i, convert{d}, j, convert{mig})
	if err != nil {
		return nil, fmt.Errorf("failed to get device name: %v", err)
	}
	var deviceSpecs []specs.Device
	for _, name := range names {
		spec := specs.Device{
			Name:           name,
			ContainerEdits: *edits.ContainerEdits,
		}
		deviceSpecs = append(deviceSpecs, spec)
	}
	return deviceSpecs, nil
}

// GetMIGDeviceEdits returns the CDI edits for the MIG device represented by 'mig' on 'parent'.
func (l *efmllib) GetMIGDeviceEdits(parent device.Device, mig device.MigDevice) (*cdi.ContainerEdits, error) {
	gcu, ret := parent.GetMinorNumber()
	if ret != efml.SUCCESS {
		return nil, fmt.Errorf("error getting GCU minor: %v", ret)
	}

	gi, ret := mig.GetGcuInstanceId()
	if ret != efml.SUCCESS {
		return nil, fmt.Errorf("error getting GCU Instance ID: %v", ret)
	}

	ci, ret := mig.GetComputeInstanceId()
	if ret != efml.SUCCESS {
		return nil, fmt.Errorf("error getting Compute Instance ID: %v", ret)
	}

	editsForDevice, err := l.GetEditsForComputeInstance(gcu, gi, ci)
	if err != nil {
		return nil, fmt.Errorf("failed to create container edits for MIG device: %v", err)
	}

	return editsForDevice, nil
}

// GetEditsForComputeInstance returns the CDI edits for a particular compute instance defined by the (gcu, gi, ci) tuple
func (l *efmllib) GetEditsForComputeInstance(gcu int, gi int, ci int) (*cdi.ContainerEdits, error) {
	computeInstance, err := newComputeInstanceDiscoverer(l.logger, l.devRoot, gcu, gi, ci)
	if err != nil {
		return nil, fmt.Errorf("failed to create discoverer for Compute Instance: %v", err)
	}

	editsForDevice, err := edits.FromDiscoverer(computeInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to create container edits for Compute Instance: %v", err)
	}

	return editsForDevice, nil
}

// newComputeInstanceDiscoverer returns a discoverer for the specified compute instance
func newComputeInstanceDiscoverer(logger logger.Interface, devRoot string, gcu int, gi int, ci int) (discover.Discover, error) {
	parentPath := fmt.Sprintf("/dev/enflame%d", gcu)

	migCaps, err := efcaps.NewMigCaps()
	if err != nil {
		return nil, fmt.Errorf("error getting MIG capability device paths: %v", err)
	}

	giCap := efcaps.NewGCUInstanceCap(gcu, gi)
	giCapDevicePath, err := migCaps.GetCapDevicePath(giCap)
	if err != nil {
		return nil, fmt.Errorf("failed to get GI cap device path: %v", err)
	}

	ciCap := efcaps.NewComputeInstanceCap(gcu, gi, ci)
	ciCapDevicePath, err := migCaps.GetCapDevicePath(ciCap)
	if err != nil {
		return nil, fmt.Errorf("failed to get CI cap device path: %v", err)
	}

	deviceNodes := discover.NewCharDeviceDiscoverer(
		logger,
		devRoot,
		[]string{
			parentPath,
			giCapDevicePath,
			ciCapDevicePath,
		},
	)

	return deviceNodes, nil
}
*/