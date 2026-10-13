/*
# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
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
*/

package devices

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Device major numbers and device names for ENFLAME devices
const (
	ENFLAMEUVMMinor      = 0
	ENFLAMEUVMToolsMinor = 1
	ENFLAMECTLMinor      = 255
	ENFLAMEModesetMinor  = 254

	ENFLAMEFrontend = Name("enflame-frontend")
	ENFLAMEGCU      = Name("enflame")
	ENFLAMECaps     = Name("enflame-caps")
	ENFLAMEUVM      = Name("enflame-uvm")

	procDevicesPath    = "/proc/devices"
	enflameDevicePrefix = "enflame"
)

// Name represents the name of a device as specified under /proc/devices
type Name string

// Major represents a device major as specified under /proc/devices
type Major int

// Devices represents the set of devices under /proc/devices
//
//go:generate moq -stub -out devices_mock.go . Devices
type Devices interface {
	Exists(Name) bool
	Get(Name) (Major, bool)
	Count() int
}

type devices map[Name]Major

var _ Devices = devices(nil)

// Count returns the number of devices defined.
func (d devices) Count() int {
	return len(d)
}

// Exists checks if a Device with a given name exists or not
func (d devices) Exists(name Name) bool {
	_, exists := d.Get(name)
	return exists
}

// Get a Device from Devices. It also has fallback logic to ensure device name changes in /proc/devices are handled
// For e.g:- For GCU drivers 550.40.x or greater, the gcu device has been renamed from  "enflame-frontend" to "enflame".
func (d devices) Get(name Name) (Major, bool) {
	for _, n := range name.getWithFallback() {
		device, exists := d[n]
		if exists {
			return device, true
		}
	}
	return 0, false
}

// getWithFallback returns a prioritised list of device names for a specific name.
// This allows multiple names to be associated with a single name to support various driver versions.
func (n Name) getWithFallback() []Name {
	if n == ENFLAMEGCU || n == ENFLAMEFrontend {
		return []Name{ENFLAMEGCU, ENFLAMEFrontend}
	}
	return []Name{n}
}

// GetENFLAMEDevices returns the set of ENFLAME Devices on the machine
func GetENFLAMEDevices() (Devices, error) {
	return enflameDevices(procDevicesPath)
}

// enflameDevices returns the set of ENFLAME Devices from the specified devices file.
// This is useful for testing since we may be testing on a system where `/proc/devices` does
// contain a reference to ENFLAME devices.
func enflameDevices(devicesPath string) (Devices, error) {
	devicesFile, err := os.Open(devicesPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error opening devices file: %v", err)
	}
	defer devicesFile.Close()

	return enflameDeviceFrom(devicesFile)
}

var errNoEnflameDevices = errors.New("no ENFLAME devices found")

func enflameDeviceFrom(reader io.Reader) (Devices, error) {
	allDevices := devicesFrom(reader)

	enflameDevices := New(
		WithDeviceToMajor(allDevices),
		WithFilter(func(n string) bool {
			return !strings.HasPrefix(n, enflameDevicePrefix)
		}),
	)
	if enflameDevices.Count() == 0 {
		return nil, errNoEnflameDevices
	}
	return enflameDevices, nil
}

func devicesFrom(reader io.Reader) map[string]int {
	allDevices := make(map[string]int)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		device, major, err := processProcDeviceLine(scanner.Text())
		if err != nil {
			continue
		}
		allDevices[device] = major
	}
	return allDevices
}

func processProcDeviceLine(line string) (string, int, error) {
	trimmed := strings.TrimSpace(line)

	var name string
	var major int

	n, _ := fmt.Sscanf(trimmed, "%d %s", &major, &name)
	if n == 2 {
		return name, major, nil
	}

	return "", 0, fmt.Errorf("unparsable line: %v", line)
}
