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
	// "github.com/ENFLAME/go-eflib/pkg/efml"

	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/pkg/efcdi/transform"
)

// Option is a function that configures the efcdilib
type Option func(*efcdilib)

// EF_CLOSE
/*
// WithDeviceLib sets the device library for the library
func WithDeviceLib(devicelib device.Interface) Option {
	return func(l *efcdilib) {
		l.devicelib = devicelib
	}
}
*/
// WithDeviceNamers sets the device namer for the library
func WithDeviceNamers(namers ...DeviceNamer) Option {
	return func(l *efcdilib) {
		l.deviceNamers = namers
	}
}

// WithDriverRoot sets the driver root for the library
func WithDriverRoot(root string) Option {
	return func(l *efcdilib) {
		l.driverRoot = root
	}
}

// WithDevRoot sets the root where /dev is located.
func WithDevRoot(root string) Option {
	return func(l *efcdilib) {
		l.devRoot = root
	}
}

// WithLogger sets the logger for the library
func WithLogger(logger logger.Interface) Option {
	return func(l *efcdilib) {
		l.logger = logger
	}
}

// WithENFLAMECTKPath sets the path to the ENFLAME Container Toolkit CLI path for the library
func WithENFLAMECTKPath(path string) Option {
	return func(l *efcdilib) {
		l.enflameCTKPath = path
	}
}

// WithLdconfigPath sets the path to the ldconfig program
func WithLdconfigPath(path string) Option {
	return func(l *efcdilib) {
		l.ldconfigPath = path
	}
}

// EF_CLOSE
/*
// WithEfmlLib sets the efml library for the library
func WithEfmlLib(efmllib efml.Interface) Option {
	return func(l *efcdilib) {
		l.efmllib = efmllib
	}
}
*/

// WithMode sets the discovery mode for the library
func WithMode(mode string) Option {
	return func(l *efcdilib) {
		l.mode = mode
	}
}

// WithVendor sets the vendor for the library
func WithVendor(vendor string) Option {
	return func(o *efcdilib) {
		o.vendor = vendor
	}
}

// WithClass sets the class for the library
func WithClass(class string) Option {
	return func(o *efcdilib) {
		o.class = class
	}
}

// WithMergedDeviceOptions sets the merged device options for the library
// If these are not set, no merged device will be generated.
func WithMergedDeviceOptions(opts ...transform.MergedDeviceOption) Option {
	return func(o *efcdilib) {
		o.mergedDeviceOptions = opts
	}
}

// WithCSVFiles sets the CSV files for the library
func WithCSVFiles(csvFiles []string) Option {
	return func(o *efcdilib) {
		o.csvFiles = csvFiles
	}
}

// WithCSVIgnorePatterns sets the ignore patterns for entries in the CSV files.
func WithCSVIgnorePatterns(csvIgnorePatterns []string) Option {
	return func(o *efcdilib) {
		o.csvIgnorePatterns = csvIgnorePatterns
	}
}

// WithConfigSearchPaths sets the search paths for config files.
func WithConfigSearchPaths(paths []string) Option {
	return func(o *efcdilib) {
		o.configSearchPaths = paths
	}
}

// WithLibrarySearchPaths sets the library search paths.
// This is currently only used for CSV-mode.
func WithLibrarySearchPaths(paths []string) Option {
	return func(o *efcdilib) {
		o.librarySearchPaths = paths
	}
}
