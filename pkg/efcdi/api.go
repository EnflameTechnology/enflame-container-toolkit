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
	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	"tags.cncf.io/container-device-interface/pkg/cdi"
	// "tags.cncf.io/container-device-interface/specs-go"

	"enflame/enflame-container-toolkit/pkg/efcdi/spec"
)

const (
	// ModeAuto configures the CDI spec generator to automatically detect the system configuration
	ModeAuto = "auto"
	// ModeEfml configures the CDI spec generator to use the EFML library.
	ModeEfml = "efml"
	// ModeWsl configures the CDI spec generator to generate a WSL spec.
	ModeWsl = "wsl"
	// ModeManagement configures the CDI spec generator to generate a management spec.
	ModeManagement = "management"
	// ModeGds configures the CDI spec generator to generate a GDS spec.
	ModeGds = "gds"
	// ModeMofed configures the CDI spec generator to generate a MOFED spec.
	ModeMofed = "mofed"
	// ModeCSV configures the CDI spec generator to generate a spec based on the contents of CSV
	// mountspec files.
	ModeCSV = "csv"
)

// Interface defines the API for the efcdi package
type Interface interface {
	GetSpec() (spec.Interface, error)
	GetCommonEdits() (*cdi.ContainerEdits, error)
// EF_CLOSE
/*
	GetAllDeviceSpecs() ([]specs.Device, error)
	GetGCUDeviceEdits(device.Device) (*cdi.ContainerEdits, error)
	GetGCUDeviceSpecs(int, device.Device) ([]specs.Device, error)
	GetMIGDeviceEdits(device.Device, device.MigDevice) (*cdi.ContainerEdits, error)
	GetMIGDeviceSpecs(int, device.Device, int, device.MigDevice) ([]specs.Device, error)
	GetDeviceSpecsByID(...string) ([]specs.Device, error)
*/
}
