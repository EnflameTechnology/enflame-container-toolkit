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
	"strconv"
	"strings"

	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	// "github.com/ENFLAME/go-eflib/pkg/efml"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	"tags.cncf.io/container-device-interface/specs-go"

	"enflame/enflame-container-toolkit/internal/edits"
	"enflame/enflame-container-toolkit/pkg/efcdi/spec"
)

type efmllib efcdilib

var _ Interface = (*efmllib)(nil)

// GetSpec should not be called for efmllib
func (l *efmllib) GetSpec() (spec.Interface, error) {
	return nil, fmt.Errorf("Unexpected call to efmllib.GetSpec()")
}

// GetAllDeviceSpecs returns the device specs for all available devices.
func (l *efmllib) GetAllDeviceSpecs() ([]specs.Device, error) {
	var deviceSpecs []specs.Device

	if r := l.efmllib.Init(); r != efml.SUCCESS {
		return nil, fmt.Errorf("failed to initialize EFML: %v", r)
	}
	defer func() {
		if r := l.efmllib.Shutdown(); r != efml.SUCCESS {
			l.logger.Warningf("failed to shutdown EFML: %v", r)
		}
	}()

	gcuDeviceSpecs, err := l.getGCUDeviceSpecs()
	if err != nil {
		return nil, err
	}
	deviceSpecs = append(deviceSpecs, gcuDeviceSpecs...)

	migDeviceSpecs, err := l.getMigDeviceSpecs()
	if err != nil {
		return nil, err
	}
	deviceSpecs = append(deviceSpecs, migDeviceSpecs...)

	return deviceSpecs, nil
}

// GetCommonEdits generates a CDI specification that can be used for ANY devices
func (l *efmllib) GetCommonEdits() (*cdi.ContainerEdits, error) {
	common, err := l.newCommonEFMLDiscoverer()
	if err != nil {
		return nil, fmt.Errorf("failed to create discoverer for common entities: %v", err)
	}

	return edits.FromDiscoverer(common)
}

// GetDeviceSpecsByID returns the CDI device specs for the GCU(s) represented by
// the provided identifiers, where an identifier is an index or UUID of a valid
// GCU device.
func (l *efmllib) GetDeviceSpecsByID(identifiers ...string) ([]specs.Device, error) {
	for _, id := range identifiers {
		if id == "all" {
			return l.GetAllDeviceSpecs()
		}
	}

	var deviceSpecs []specs.Device

	if r := l.efmllib.Init(); r != efml.SUCCESS {
		return nil, fmt.Errorf("failed to initialize EFML: %w", r)
	}
	defer func() {
		if r := l.efmllib.Shutdown(); r != efml.SUCCESS {
			l.logger.Warningf("failed to shutdown EFML: %w", r)
		}
	}()

	efmlDevices, err := l.getEFMLDevicesByID(identifiers...)
	if err != nil {
		return nil, fmt.Errorf("failed to get EFML device handles: %w", err)
	}

	for i, efmlDevice := range efmlDevices {
		deviceEdits, err := l.getEditsForDevice(efmlDevice)
		if err != nil {
			return nil, fmt.Errorf("failed to get CDI device edits for identifier %q: %w", identifiers[i], err)
		}
		deviceSpec := specs.Device{
			Name:           identifiers[i],
			ContainerEdits: *deviceEdits.ContainerEdits,
		}
		deviceSpecs = append(deviceSpecs, deviceSpec)
	}

	return deviceSpecs, nil
}

// TODO: move this to go-eflib?
func (l *efmllib) getEFMLDevicesByID(identifiers ...string) ([]efml.Device, error) {
	var devices []efml.Device
	for _, id := range identifiers {
		dev, err := l.getEFMLDeviceByID(id)
		if err != efml.SUCCESS {
			return nil, fmt.Errorf("failed to get EFML device handle for identifier %q: %w", id, err)
		}
		devices = append(devices, dev)
	}
	return devices, nil
}

func (l *efmllib) getEFMLDeviceByID(id string) (efml.Device, error) {
	var err error
	devID := device.Identifier(id)

	if devID.IsUUID() {
		return l.efmllib.DeviceGetHandleByUUID(id)
	}

	if devID.IsGcuIndex() {
		if idx, err := strconv.Atoi(id); err == nil {
			return l.efmllib.DeviceGetHandleByIndex(idx)
		}
		return nil, fmt.Errorf("failed to convert device index to an int: %w", err)
	}

	if devID.IsMigIndex() {
		var gcuIdx, migIdx int
		var parent efml.Device
		split := strings.SplitN(id, ":", 2)
		if gcuIdx, err = strconv.Atoi(split[0]); err != nil {
			return nil, fmt.Errorf("failed to convert device index to an int: %w", err)
		}
		if migIdx, err = strconv.Atoi(split[1]); err != nil {
			return nil, fmt.Errorf("failed to convert device index to an int: %w", err)
		}
		if parent, err = l.efmllib.DeviceGetHandleByIndex(gcuIdx); err != efml.SUCCESS {
			return nil, fmt.Errorf("failed to get parent device handle: %w", err)
		}
		return parent.GetMigDeviceHandleByIndex(migIdx)
	}

	return nil, fmt.Errorf("identifier is not a valid UUID or index: %q", id)
}

func (l *efmllib) getEditsForDevice(efmlDevice efml.Device) (*cdi.ContainerEdits, error) {
	mig, err := efmlDevice.IsMigDeviceHandle()
	if err != efml.SUCCESS {
		return nil, fmt.Errorf("failed to determine if device handle is a MIG device: %w", err)
	}
	if mig {
		return l.getEditsForMIGDevice(efmlDevice)
	}
	return l.getEditsForGCUDevice(efmlDevice)
}

func (l *efmllib) getEditsForGCUDevice(efmlDevice efml.Device) (*cdi.ContainerEdits, error) {
	nvlibDevice, err := l.devicelib.NewDevice(efmlDevice)
	if err != nil {
		return nil, fmt.Errorf("failed to construct device: %w", err)
	}
	deviceEdits, err := l.GetGCUDeviceEdits(nvlibDevice)
	if err != nil {
		return nil, fmt.Errorf("failed to get GCU device edits: %w", err)
	}

	return deviceEdits, nil
}

func (l *efmllib) getEditsForMIGDevice(efmlDevice efml.Device) (*cdi.ContainerEdits, error) {
	efmlParentDevice, ret := efmlDevice.GetDeviceHandleFromMigDeviceHandle()
	if ret != efml.SUCCESS {
		return nil, fmt.Errorf("failed to get parent device handle: %w", ret)
	}
	nvlibMigDevice, err := l.devicelib.NewMigDevice(efmlDevice)
	if err != nil {
		return nil, fmt.Errorf("failed to construct device: %w", err)
	}
	nvlibParentDevice, err := l.devicelib.NewDevice(efmlParentDevice)
	if err != nil {
		return nil, fmt.Errorf("failed to construct parent device: %w", err)
	}
	return l.GetMIGDeviceEdits(nvlibParentDevice, nvlibMigDevice)
}

func (l *efmllib) getGCUDeviceSpecs() ([]specs.Device, error) {
	var deviceSpecs []specs.Device
	err := l.devicelib.VisitDevices(func(i int, d device.Device) error {
		specsForDevice, err := l.GetGCUDeviceSpecs(i, d)
		if err != nil {
			return err
		}
		deviceSpecs = append(deviceSpecs, specsForDevice...)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate CDI edits for GCU devices: %v", err)
	}
	return deviceSpecs, err
}

func (l *efmllib) getMigDeviceSpecs() ([]specs.Device, error) {
	var deviceSpecs []specs.Device
	err := l.devicelib.VisitMigDevices(func(i int, d device.Device, j int, mig device.MigDevice) error {
		specsForDevice, err := l.GetMIGDeviceSpecs(i, d, j, mig)
		if err != nil {
			return err
		}
		deviceSpecs = append(deviceSpecs, specsForDevice...)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate CDI edits for GCU devices: %v", err)
	}
	return deviceSpecs, err
}
*/