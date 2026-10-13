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

package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/opencontainers/runtime-spec/specs-go"

	"enflame/enflame-container-toolkit/internal/config"
	"enflame/enflame-container-toolkit/internal/info"
	"enflame/enflame-container-toolkit/internal/lookup/root"
)

// Run is an entry point that allows for idiomatic handling of errors
// when calling from the main function.
func (r rt) Run(argv []string) (rerr error) {
	defer func() {
		if rerr != nil {
			r.logger.Errorf("%v", rerr)
		}
	}()

	printVersion := hasVersionFlag(argv)
	if printVersion {
		fmt.Printf("%v version %v\n", "ENFLAME Container Runtime", info.GetVersionString(fmt.Sprintf("spec: %v", specs.Version)))
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}
	r.logger.Update(
		cfg.ENFLAMEContainerRuntimeConfig.DebugFilePath,
		cfg.ENFLAMEContainerRuntimeConfig.LogLevel,
		argv,
	)
	defer func() {
		if rerr != nil {
			r.logger.Errorf("%v", rerr)
		}
		if err := r.logger.Reset(); err != nil {
			rerr = errors.Join(rerr, fmt.Errorf("failed to reset logger: %v", err))
		}
	}()

	// We apply some config updates here to ensure that the config is valid in
	// all cases.
	if r.modeOverride != "" {
		cfg.ENFLAMEContainerRuntimeConfig.Mode = r.modeOverride
	}
	cfg.ENFLAMECTKConfig.Path = config.ResolveENFLAMECTKPath(r.logger, cfg.ENFLAMECTKConfig.Path)
	cfg.ENFLAMEContainerRuntimeHookConfig.Path = config.ResolveENFLAMEContainerRuntimeHookPath(r.logger, cfg.ENFLAMEContainerRuntimeHookConfig.Path)

	// Print the config to the output.
	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err == nil {
		r.logger.Infof("Running with config:\n%v", string(configJSON))
	} else {
		r.logger.Infof("Running with config:\n%+v", cfg)
	}

	driver := root.New(
		root.WithLogger(r.logger),
		root.WithDriverRoot(cfg.ENFLAMEContainerCLIConfig.Root),
	)

	r.logger.Debugf("Command line arguments: %v", argv)
	runtime, err := newENFLAMEContainerRuntime(r.logger, cfg, argv, driver)
	if err != nil {
		return fmt.Errorf("failed to create ENFLAME Container Runtime: %v", err)
	}

	if printVersion {
		fmt.Print("\n")
	}
	return runtime.Exec(argv)
}

func (r rt) Errorf(format string, args ...interface{}) {
	r.logger.Errorf(format, args...)
}

// TODO: This should be refactored / combined with parseArgs in logger.
func hasVersionFlag(args []string) bool {
	for i := 0; i < len(args); i++ {
		param := args[i]

		parts := strings.SplitN(param, "=", 2)
		trimmed := strings.TrimLeft(parts[0], "-")
		// If this is not a flag we continue
		if parts[0] == trimmed {
			continue
		}

		// Check the version flag
		if trimmed == "version" {
			return true
		}
	}

	return false
}
