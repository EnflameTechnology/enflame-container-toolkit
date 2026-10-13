/*
# Copyright (c) 2021-2022, NVIDIA CORPORATION.  All rights reserved.
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

package runtime

import (
	"fmt"

	"enflame/enflame-container-toolkit/internal/config"
	"enflame/enflame-container-toolkit/internal/config/image"
	"enflame/enflame-container-toolkit/internal/info"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/lookup/root"
	"enflame/enflame-container-toolkit/internal/modifier"
	"enflame/enflame-container-toolkit/internal/oci"
)

// newENFLAMEContainerRuntime is a factory method that constructs a runtime based on the selected configuration and specified logger
func newENFLAMEContainerRuntime(logger logger.Interface, cfg *config.Config, argv []string, driver *root.Driver) (oci.Runtime, error) {
	lowLevelRuntime, err := oci.NewLowLevelRuntime(logger, cfg.ENFLAMEContainerRuntimeConfig.Runtimes)
	if err != nil {
		return nil, fmt.Errorf("error constructing low-level runtime: %v", err)
	}

	if !oci.HasCreateSubcommand(argv) {
		logger.Debugf("Skipping modifier for non-create subcommand")
		return lowLevelRuntime, nil
	}

	ociSpec, err := oci.NewSpec(logger, argv)
	if err != nil {
		return nil, fmt.Errorf("error constructing OCI specification: %v", err)
	}

	specModifier, err := newSpecModifier(logger, cfg, ociSpec, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to construct OCI spec modifier: %v", err)
	}

	// Create the wrapping runtime with the specified modifier
	r := oci.NewModifyingRuntimeWrapper(
		logger,
		lowLevelRuntime,
		ociSpec,
		specModifier,
	)

	return r, nil
}

// newSpecModifier is a factory method that creates constructs an OCI spec modifer based on the provided config.
func newSpecModifier(logger logger.Interface, cfg *config.Config, ociSpec oci.Spec, driver *root.Driver) (oci.SpecModifier, error) {
	rawSpec, err := ociSpec.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load OCI spec: %v", err)
	}

	image, err := image.NewTOPSImageFromSpec(rawSpec)
	if err != nil {
		return nil, err
	}

	mode := info.ResolveAutoMode(logger, cfg.ENFLAMEContainerRuntimeConfig.Mode, image)
	modeModifier, err := newModeModifier(logger, mode, cfg, ociSpec, image)
	if err != nil {
		return nil, err
	}
	// For CDI mode we make no additional modifications.
	if mode == "cdi" {
		return modeModifier, nil
	}

	graphicsModifier, err := modifier.NewGraphicsModifier(logger, cfg, image, driver)
	if err != nil {
		return nil, err
	}

	featureModifier, err := modifier.NewFeatureGatedModifier(logger, cfg, image)
	if err != nil {
		return nil, err
	}

	modifiers := modifier.Merge(
		modeModifier,
		graphicsModifier,
		featureModifier,
	)
	return modifiers, nil
}

func newModeModifier(logger logger.Interface, mode string, cfg *config.Config, ociSpec oci.Spec, image image.TOPS) (oci.SpecModifier, error) {
	switch mode {
	case "legacy":
		return modifier.NewStableRuntimeModifier(logger, cfg.ENFLAMEContainerRuntimeHookConfig.Path), nil
	case "csv":
		return modifier.NewCSVModifier(logger, cfg, image)
	case "cdi":
		return modifier.NewCDIModifier(logger, cfg, ociSpec)
	}

	return nil, fmt.Errorf("invalid runtime mode: %v", cfg.ENFLAMEContainerRuntimeConfig.Mode)
}
