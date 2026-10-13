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
// EF_CLOSE
/*
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	// "github.com/ENFLAME/go-eflib/pkg/efml"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"

	"enflame/enflame-container-toolkit/internal/discover"
	"enflame/enflame-container-toolkit/internal/edits"
	"enflame/enflame-container-toolkit/internal/info/drm"
	"enflame/enflame-container-toolkit/internal/logger"
)

// GetGCUDeviceSpecs returns the CDI device specs for the full GCU represented by 'device'.
func (l *efmllib) GetGCUDeviceSpecs(i int, d device.Device) ([]specs.Device, error) {
	edits, err := l.GetGCUDeviceEdits(d)
	if err != nil {
		return nil, fmt.Errorf("failed to get edits for device: %v", err)
	}

	var deviceSpecs []specs.Device
	names, err := l.deviceNamers.GetDeviceNames(i, convert{d})
	if err != nil {
		return nil, fmt.Errorf("failed to get device name: %v", err)
	}
	for _, name := range names {
		spec := specs.Device{
			Name:           name,
			ContainerEdits: *edits.ContainerEdits,
		}
		deviceSpecs = append(deviceSpecs, spec)
	}

	return deviceSpecs, nil
}

// GetGCUDeviceEdits returns the CDI edits for the full GCU represented by 'device'.
func (l *efmllib) GetGCUDeviceEdits(d device.Device) (*cdi.ContainerEdits, error) {
	device, err := newFullGCUDiscoverer(l.logger, l.devRoot, l.enflameCTKPath, d)
	if err != nil {
		return nil, fmt.Errorf("failed to create device discoverer: %v", err)
	}

	editsForDevice, err := edits.FromDiscoverer(device)
	if err != nil {
		return nil, fmt.Errorf("failed to create container edits for device: %v", err)
	}

	return editsForDevice, nil
}

// byPathHookDiscoverer discovers the entities required for injecting by-path DRM device links
type byPathHookDiscoverer struct {
	logger        logger.Interface
	devRoot       string
	enflameCTKPath string
	pciBusID      string
	deviceNodes   discover.Discover
}

var _ discover.Discover = (*byPathHookDiscoverer)(nil)

// newFullGCUDiscoverer creates a discoverer for the full GCU defined by the specified device.
func newFullGCUDiscoverer(logger logger.Interface, devRoot string, enflameCTKPath string, d device.Device) (discover.Discover, error) {
	// TODO: The functionality to get device paths should be integrated into the go-eflib/pkg/device.Device interface.
	// This will allow reuse here and in other code where the paths are queried such as the ENFLAME device plugin.
	minor, ret := d.GetMinorNumber()
	if ret != efml.SUCCESS {
		return nil, fmt.Errorf("error getting GCU device minor number: %v", ret)
	}
	path := fmt.Sprintf("/dev/enflame%d", minor)

	pciInfo, ret := d.GetPciInfo()
	if ret != efml.SUCCESS {
		return nil, fmt.Errorf("error getting PCI info for device: %v", ret)
	}
	pciBusID := getBusID(pciInfo)

	drmDeviceNodes, err := drm.GetDeviceNodesByBusID(pciBusID)
	if err != nil {
		return nil, fmt.Errorf("failed to determine DRM devices for %v: %v", pciBusID, err)
	}

	deviceNodePaths := append([]string{path}, drmDeviceNodes...)

	deviceNodes := discover.NewCharDeviceDiscoverer(
		logger,
		devRoot,
		deviceNodePaths,
	)

	byPathHooks := &byPathHookDiscoverer{
		logger:        logger,
		devRoot:       devRoot,
		enflameCTKPath: enflameCTKPath,
		pciBusID:      pciBusID,
		deviceNodes:   deviceNodes,
	}

	deviceFolderPermissionHooks := newDeviceFolderPermissionHookDiscoverer(
		logger,
		devRoot,
		enflameCTKPath,
		deviceNodes,
	)

	dd := discover.Merge(
		deviceNodes,
		byPathHooks,
		deviceFolderPermissionHooks,
	)

	return dd, nil
}

// Devices returns the empty list for the by-path hook discoverer
func (d *byPathHookDiscoverer) Devices() ([]discover.Device, error) {
	return nil, nil
}

// Hooks returns the hooks for the GCU device.
// The following hooks are detected:
//  1. A hook to create /dev/dri/by-path symlinks
func (d *byPathHookDiscoverer) Hooks() ([]discover.Hook, error) {
	links, err := d.deviceNodeLinks()
	if err != nil {
		return nil, fmt.Errorf("failed to discover DRA device links: %v", err)
	}
	if len(links) == 0 {
		return nil, nil
	}

	var args []string
	for _, l := range links {
		args = append(args, "--link", l)
	}

	hook := discover.CreateEnflameCTKHook(
		d.enflameCTKPath,
		"create-symlinks",
		args...,
	)

	return []discover.Hook{hook}, nil
}

// Mounts returns an empty slice for a full GCU
func (d *byPathHookDiscoverer) Mounts() ([]discover.Mount, error) {
	return nil, nil
}

func (d *byPathHookDiscoverer) deviceNodeLinks() ([]string, error) {
	devices, err := d.deviceNodes.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to discover device nodes: %v", err)
	}

	if len(devices) == 0 {
		return nil, nil
	}

	selectedDevices := make(map[string]bool)
	for _, d := range devices {
		selectedDevices[d.HostPath] = true
	}

	candidates := []string{
		fmt.Sprintf("/dev/dri/by-path/pci-%s-card", d.pciBusID),
		fmt.Sprintf("/dev/dri/by-path/pci-%s-render", d.pciBusID),
	}

	var links []string
	for _, c := range candidates {
		linkPath := filepath.Join(d.devRoot, c)
		device, err := os.Readlink(linkPath)
		if err != nil {
			d.logger.Warningf("Failed to evaluate symlink %v; ignoring", linkPath)
			continue
		}

		deviceNode := device
		if !filepath.IsAbs(device) {
			deviceNode = filepath.Join(filepath.Dir(linkPath), device)
		}
		if !selectedDevices[deviceNode] {
			d.logger.Debugf("ignoring device symlink %v -> %v since %v is not mounted", linkPath, device, deviceNode)
			continue
		}
		d.logger.Debugf("adding device symlink %v -> %v", linkPath, device)
		links = append(links, fmt.Sprintf("%v::%v", device, linkPath))
	}

	return links, nil
}

// getBusID provides a utility function that returns the string representation of the bus ID.
func getBusID(p efml.PciInfo) string {
	var bytes []byte
	for _, b := range p.BusId {
		if byte(b) == '\x00' {
			break
		}
		bytes = append(bytes, byte(b))
	}
	id := strings.ToLower(string(bytes))

	if id != "0000" {
		id = strings.TrimPrefix(id, "0000")
	}

	return id
}

*/