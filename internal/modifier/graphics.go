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

package modifier

import (
	"fmt"

	"enflame/enflame-container-toolkit/internal/config"
	"enflame/enflame-container-toolkit/internal/config/image"
	"enflame/enflame-container-toolkit/internal/discover"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/lookup/root"
	"enflame/enflame-container-toolkit/internal/oci"
)

// NewGraphicsModifier constructs a modifier that injects graphics-related modifications into an OCI runtime specification.
// The value of the ENFLAME_DRIVER_CAPABILITIES environment variable is checked to determine if this modification should be made.
func NewGraphicsModifier(logger logger.Interface, cfg *config.Config, image image.TOPS, driver *root.Driver) (oci.SpecModifier, error) {
	if required, reason := requiresGraphicsModifier(image); !required {
		logger.Infof("No graphics modifier required: %v", reason)
		return nil, nil
	}

	enflameCTKPath := cfg.ENFLAMECTKConfig.Path

	mounts, err := discover.NewGraphicsMountsDiscoverer(
		logger,
		driver,
		enflameCTKPath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mounts discoverer: %v", err)
	}

	// In standard usage, the devRoot is the same as the driver.Root.
	devRoot := driver.Root
	drmNodes, err := discover.NewDRMNodesDiscoverer(
		logger,
		image.DevicesFromEnvvars(visibleDevicesEnvvar),
		devRoot,
		enflameCTKPath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct discoverer: %v", err)
	}

	d := discover.Merge(
		drmNodes,
		mounts,
	)
	return NewModifierFromDiscoverer(logger, d)
}

// requiresGraphicsModifier determines whether a graphics modifier is required.
func requiresGraphicsModifier(topsImage image.TOPS) (bool, string) {
	if devices := topsImage.DevicesFromEnvvars(visibleDevicesEnvvar); len(devices.List()) == 0 {
		return false, "no devices requested"
	}

	if !topsImage.GetDriverCapabilities().Any(image.DriverCapabilityGraphics, image.DriverCapabilityDisplay) {
		return false, "no required capabilities requested"
	}

	return true, ""
}
