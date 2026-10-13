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
	"path/filepath"

	"github.com/opencontainers/runtime-spec/specs-go"

	"enflame/enflame-container-toolkit/internal/config"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/oci"
)

// enflameContainerRuntimeHookRemover is a spec modifer that detects and removes inserted enflame-container-runtime hooks
type enflameContainerRuntimeHookRemover struct {
	logger logger.Interface
}

var _ oci.SpecModifier = (*enflameContainerRuntimeHookRemover)(nil)

// Modify removes any ENFLAME Container Runtime hooks from the provided spec
func (m enflameContainerRuntimeHookRemover) Modify(spec *specs.Spec) error {
	if spec == nil {
		return nil
	}

	if spec.Hooks == nil {
		return nil
	}

	if len(spec.Hooks.Prestart) == 0 {
		return nil
	}

	var newPrestart []specs.Hook

	for _, hook := range spec.Hooks.Prestart {
		hook := hook
		if isENFLAMEContainerRuntimeHook(&hook) {
			m.logger.Debugf("Removing hook %v", hook)
			continue
		}
		newPrestart = append(newPrestart, hook)
	}

	if len(newPrestart) != len(spec.Hooks.Prestart) {
		m.logger.Debugf("Updating 'prestart' hooks to %v", newPrestart)
		spec.Hooks.Prestart = newPrestart
	}

	return nil
}

// isENFLAMEContainerRuntimeHook checks if the provided hook is an enflame-container-runtime-hook
// or enflame-container-toolkit hook. These are included, for example, by the non-experimental
// enflame-container-runtime or docker when specifying the --gcus flag.
func isENFLAMEContainerRuntimeHook(hook *specs.Hook) bool {
	bins := map[string]struct{}{
		config.ENFLAMEContainerRuntimeHookExecutable: {},
		config.ENFLAMEContainerToolkitExecutable:     {},
	}

	_, exists := bins[filepath.Base(hook.Path)]

	return exists
}
