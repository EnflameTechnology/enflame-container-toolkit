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

package efdevices

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"enflame/enflame-container-toolkit/internal/info/proc/devices"
	"enflame/enflame-container-toolkit/internal/logger"
)

var errInvalidDeviceNode = errors.New("invalid device node")

// Interface provides a set of utilities for interacting with ENFLAME devices on the system.
type Interface struct {
	devices.Devices

	logger logger.Interface

	dryRun bool
	// devRoot is the root directory where device nodes are expected to exist.
	devRoot string

	mknoder
}

// New constructs a new Interface struct with the specified options.
func New(opts ...Option) (*Interface, error) {
	i := &Interface{}
	for _, opt := range opts {
		opt(i)
	}

	if i.logger == nil {
		i.logger = logger.New()
	}
	if i.devRoot == "" {
		i.devRoot = "/"
	}
	if i.Devices == nil {
		devices, err := devices.GetENFLAMEDevices()
		if err != nil {
			return nil, fmt.Errorf("failed to create devices info: %v", err)
		}
		i.Devices = devices
	}

	if i.dryRun {
		i.mknoder = &mknodLogger{i.logger}
	} else {
		i.mknoder = &mknodUnix{}
	}
	return i, nil
}

// CreateENFLAMEControlDevices creates the ENFLAME control device nodes at the configured devRoot.
func (m *Interface) CreateENFLAMEControlDevices() error {
	controlNodes := []string{"enflamectl", "enflame-modeset", "enflame-uvm", "enflame-uvm-tools"}
	for _, node := range controlNodes {
		err := m.CreateENFLAMEDevice(node)
		if err != nil {
			return fmt.Errorf("failed to create device node %s: %w", node, err)
		}
	}
	return nil
}

// CreateENFLAMEDevice creates the specified ENFLAME device node at the configured devRoot.
func (m *Interface) CreateENFLAMEDevice(node string) error {
	node = filepath.Base(node)
	if !strings.HasPrefix(node, "enflame") {
		return fmt.Errorf("invalid device node %q: %w", node, errInvalidDeviceNode)
	}

	major, err := m.Major(node)
	if err != nil {
		return fmt.Errorf("failed to determine major: %w", err)
	}

	minor, err := m.Minor(node)
	if err != nil {
		return fmt.Errorf("failed to determine minor: %w", err)
	}

	return m.createDeviceNode(filepath.Join("dev", node), int(major), int(minor))
}

// createDeviceNode creates the specified device node with the require major and minor numbers.
// If a devRoot is configured, this is prepended to the path.
func (m *Interface) createDeviceNode(path string, major int, minor int) error {
	path = filepath.Join(m.devRoot, path)
	if _, err := os.Stat(path); err == nil {
		m.logger.Infof("Skipping: %s already exists", path)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat %s: %v", path, err)
	}

	return m.Mknode(path, major, minor)
}

// Major returns the major number for the specified ENFLAME device node.
// If the device node is not supported, an error is returned.
func (m *Interface) Major(node string) (int64, error) {
	var valid bool
	var major devices.Major
	switch node {
	case "enflame-uvm", "enflame-uvm-tools":
		major, valid = m.Get(devices.ENFLAMEUVM)
	case "enflame-modeset", "enflamectl":
		major, valid = m.Get(devices.ENFLAMEGCU)
	}

	if valid {
		return int64(major), nil
	}

	return 0, errInvalidDeviceNode
}

// Minor returns the minor number for the specified ENFLAME device node.
// If the device node is not supported, an error is returned.
func (m *Interface) Minor(node string) (int64, error) {
	switch node {
	case "enflame-modeset":
		return devices.ENFLAMEModesetMinor, nil
	case "enflame-uvm-tools":
		return devices.ENFLAMEUVMToolsMinor, nil
	case "enflame-uvm":
		return devices.ENFLAMEUVMMinor, nil
	case "enflamectl":
		return devices.ENFLAMECTLMinor, nil
	}

	return 0, errInvalidDeviceNode
}
