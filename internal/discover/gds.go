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

package discover

import (
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/lookup"
)

type gdsDeviceDiscoverer struct {
	None
	logger  logger.Interface
	devices Discover
	mounts  Discover
}

// NewGDSDiscoverer creates a discoverer for GCUDirect Storage devices and mounts.
func NewGDSDiscoverer(logger logger.Interface, driverRoot string, devRoot string) (Discover, error) {
	devices := NewCharDeviceDiscoverer(
		logger,
		devRoot,
		[]string{"/dev/enflame-fs*"},
	)

	udev := NewMounts(
		logger,
		lookup.NewDirectoryLocator(lookup.WithLogger(logger), lookup.WithRoot(driverRoot)),
		driverRoot,
		[]string{"/run/udev"},
	)

	cufile := NewMounts(
		logger,
		lookup.NewFileLocator(
			lookup.WithLogger(logger),
			lookup.WithRoot(driverRoot),
		),
		driverRoot,
		[]string{"/etc/cufile.json"},
	)

	d := gdsDeviceDiscoverer{
		logger:  logger,
		devices: devices,
		mounts:  Merge(udev, cufile),
	}

	return &d, nil
}

// Devices discovers the enflame-fs device nodes for use with GCUDirect Storage
func (d *gdsDeviceDiscoverer) Devices() ([]Device, error) {
	return d.devices.Devices()
}

// Mounts discovers the required mounts for GCUDirect Storage.
// If no devices are discovered the discovered mounts are empty
func (d *gdsDeviceDiscoverer) Mounts() ([]Mount, error) {
	devices, err := d.Devices()
	if err != nil || len(devices) == 0 {
		d.logger.Debugf("No enflame-fs devices detected; skipping detection of mounts")
		return nil, nil
	}

	return d.mounts.Mounts()
}
