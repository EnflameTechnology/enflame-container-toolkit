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
	"fmt"
	"path/filepath"

	"enflame/enflame-container-toolkit/internal/discover"
	"enflame/enflame-container-toolkit/internal/dxcore"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/lookup"
)

var requiredDriverStoreFiles = []string{
	"libtops.so.1.1",                /* Core library for tops support */
	"libtops_loader.so",             /* Core library for tops support on WSL */
	"libenflame-ptxjitcompiler.so.1", /* Core library for PTX Jit support */
	"libefml.so",                     /* Core library for efml */
	"libenflame-ml_loader.so",        /* Core library for efml on WSL */
	"libdxcore.so",                   /* Core library for dxcore support */
	"libefdxgdmal.so.1",              /* dxgdmal library for tops */
	"eftopsbins.bin",                  /* Binary containing GCU code for tops */
	"efsmi",                           /* enflame-smi binary*/
}

// newWSLDriverDiscoverer returns a Discoverer for WSL2 drivers.
func newWSLDriverDiscoverer(logger logger.Interface, driverRoot string, enflameCTKPath, ldconfigPath string) (discover.Discover, error) {
	err := dxcore.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dxcore: %v", err)
	}
	defer func() {
		if err := dxcore.Shutdown(); err != nil {
			logger.Warningf("failed to shutdown dxcore: %v", err)
		}
	}()

	driverStorePaths := dxcore.GetDriverStorePaths()
	if len(driverStorePaths) == 0 {
		return nil, fmt.Errorf("no driver store paths found")
	}
	logger.Infof("Using WSL driver store paths: %v", driverStorePaths)

	return newWSLDriverStoreDiscoverer(logger, driverRoot, enflameCTKPath, ldconfigPath, driverStorePaths)
}

// newWSLDriverStoreDiscoverer returns a Discoverer for WSL2 drivers in the driver store associated with a dxcore adapter.
func newWSLDriverStoreDiscoverer(logger logger.Interface, driverRoot string, enflameCTKPath string, ldconfigPath string, driverStorePaths []string) (discover.Discover, error) {
	var searchPaths []string
	seen := make(map[string]bool)
	for _, path := range driverStorePaths {
		if seen[path] {
			continue
		}
		searchPaths = append(searchPaths, path)
	}
	if len(searchPaths) > 1 {
		logger.Warningf("Found multiple driver store paths: %v", searchPaths)
	}
	searchPaths = append(searchPaths, "/usr/lib/wsl/lib")

	libraries := discover.NewMounts(
		logger,
		lookup.NewFileLocator(
			lookup.WithLogger(logger),
			lookup.WithSearchPaths(
				searchPaths...,
			),
			lookup.WithCount(1),
		),
		driverRoot,
		requiredDriverStoreFiles,
	)

	symlinkHook := enflameSMISimlinkHook{
		logger:        logger,
		mountsFrom:    libraries,
		enflameCTKPath: enflameCTKPath,
	}

	ldcacheHook, _ := discover.NewLDCacheUpdateHook(logger, libraries, enflameCTKPath, ldconfigPath)

	d := discover.Merge(
		libraries,
		symlinkHook,
		ldcacheHook,
	)

	return d, nil
}

type enflameSMISimlinkHook struct {
	discover.None
	logger        logger.Interface
	mountsFrom    discover.Discover
	enflameCTKPath string
}

// Hooks returns a hook that creates a symlink to enflame-smi in the driver store.
// On WSL2 the driver store location is used unchanged, for this reason we need
// to create a symlink from /usr/bin/enflame-smi to the enflame-smi binary in the
// driver store.
func (m enflameSMISimlinkHook) Hooks() ([]discover.Hook, error) {
	mounts, err := m.mountsFrom.Mounts()
	if err != nil {
		return nil, fmt.Errorf("failed to discover mounts: %w", err)
	}

	var target string
	for _, mount := range mounts {
		if filepath.Base(mount.Path) == "enflame-smi" {
			target = mount.Path
			break
		}
	}

	if target == "" {
		m.logger.Warningf("Failed to find enflame-smi in mounts: %v", mounts)
		return nil, nil
	}
	link := "/usr/bin/enflame-smi"
	links := []string{fmt.Sprintf("%s::%s", target, link)}
	symlinkHook := discover.CreateCreateSymlinkHook(m.enflameCTKPath, links)

	return symlinkHook.Hooks()
}
