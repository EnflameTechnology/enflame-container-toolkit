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
	// "enflame/enflame-container-toolkit/internal/tops"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/modifier/cdi"
	"enflame/enflame-container-toolkit/internal/oci"
	"enflame/enflame-container-toolkit/internal/platform-support/tegra/csv"
	"enflame/enflame-container-toolkit/internal/requirements"
	"enflame/enflame-container-toolkit/pkg/efcdi"
)

const (
	visibleDevicesEnvvar = "ENFLAME_VISIBLE_DEVICES"
	visibleDevicesVoid   = "void"

	enflameRequireJetpackEnvvar = "ENFLAME_REQUIRE_JETPACK"
)

// NewCSVModifier creates a modifier that applies modications to an OCI spec if required by the runtime wrapper.
// The modifications are defined by CSV MountSpecs.
func NewCSVModifier(logger logger.Interface, cfg *config.Config, image image.TOPS) (oci.SpecModifier, error) {
	if devices := image.DevicesFromEnvvars(visibleDevicesEnvvar); len(devices.List()) == 0 {
		logger.Infof("No modification required; no devices requested")
		return nil, nil
	}
	logger.Infof("Constructing modifier from config: %+v", *cfg)

	if err := checkRequirements(logger, image); err != nil {
		return nil, fmt.Errorf("requirements not met: %v", err)
	}

	csvFiles, err := csv.GetFileList(cfg.ENFLAMEContainerRuntimeConfig.Modes.CSV.MountSpecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of CSV files: %v", err)
	}

	if image.Getenv(enflameRequireJetpackEnvvar) != "csv-mounts=all" {
		csvFiles = csv.BaseFilesOnly(csvFiles)
	}

	cdilib, err := efcdi.New(
		efcdi.WithLogger(logger),
		efcdi.WithDriverRoot(cfg.ENFLAMEContainerCLIConfig.Root),
		efcdi.WithENFLAMECTKPath(cfg.ENFLAMECTKConfig.Path),
		efcdi.WithMode(efcdi.ModeCSV),
		efcdi.WithCSVFiles(csvFiles),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct CDI library: %v", err)
	}

	spec, err := cdilib.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("failed to get CDI spec: %v", err)
	}

	cdiModifier, err := cdi.New(
		cdi.WithLogger(logger),
		cdi.WithSpec(spec.Raw()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to construct CDI modifier: %v", err)
	}

	modifiers := Merge(
		enflameContainerRuntimeHookRemover{logger},
		cdiModifier,
	)

	return modifiers, nil
}

func checkRequirements(logger logger.Interface, image image.TOPS) error {
	if image.HasDisableRequire() {
		// TODO: We could print the real value here instead
		logger.Debugf("ENFLAME_DISABLE_REQUIRE=%v; skipping requirement checks", true)
		return nil
	}

	imageRequirements, err := image.GetRequirements()
	if err != nil {
		//  TODO: Should we treat this as a failure, or just issue a warning?
		return fmt.Errorf("failed to get image requirements: %v", err)
	}

	r := requirements.New(logger, imageRequirements)

// EF_CLOSE
/*
	topsVersion, err := tops.Version()
	if err != nil {
		logger.Warningf("Failed to get TOPS version: %v", err)
	} else {
		r.AddVersionProperty(requirements.TOPS, topsVersion)
	}

	compteCapability, err := tops.ComputeCapability(0)
	if err != nil {
		logger.Warningf("Failed to get TOPS Compute Capability: %v", err)
	} else {
		r.AddVersionProperty(requirements.ARCH, compteCapability)
	}
*/
	return r.Assert()
}
