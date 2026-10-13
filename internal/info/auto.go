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

package info

import (
	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	// "github.com/ENFLAME/go-eflib/pkg/eflib/info"
	// "github.com/ENFLAME/go-eflib/pkg/efml"

	"enflame/enflame-container-toolkit/internal/config/image"
	"enflame/enflame-container-toolkit/internal/logger"
)

// infoInterface provides an alias for mocking.
//
//go:generate moq -stub -out info-interface_mock.go . infoInterface
type infoInterface interface {
	// EF_CLOSE
	// info.Interface

	// UsesEFGCUModule indicates whether the system is using the efgcu kernel module
	UsesEFGCUModule() (bool, string)
}

type resolver struct {
	logger logger.Interface
	// EF_CLOSE
	// info   infoInterface
}

// ResolveAutoMode determines the correct mode for the platform if set to "auto"
func ResolveAutoMode(logger logger.Interface, mode string, image image.TOPS) (rmode string) {
// EF_CLOSE
/*
	efinfo := info.New()
	efmllib := efml.New()
	devicelib := device.New(
		device.WithEfml(efmllib),
	)

	info := additionalInfo{
		Interface: efinfo,
		efmllib:   efmllib,
		devicelib: devicelib,
	}
*/
	r := resolver{
		logger: logger,
		// EF_CLOSE
		//info:   info,
	}
	return r.resolveMode(mode, image)
}

// resolveMode determines the correct mode for the platform if set to "auto"
func (r resolver) resolveMode(mode string, image image.TOPS) (rmode string) {
	if mode != "auto" {
		r.logger.Infof("Using requested mode '%s'", mode)
		return mode
	}
	defer func() {
		r.logger.Infof("Auto-detected mode as '%v'", rmode)
	}()

	if image.OnlyFullyQualifiedCDIDevices() {
		return "cdi"
	}
//EF_CLOSE
/*
	isTegra, reason := r.info.IsTegraSystem()
	r.logger.Debugf("Is Tegra-based system? %v: %v", isTegra, reason)

	hasEFML, reason := r.info.HasEfml()
	r.logger.Debugf("Has EFML? %v: %v", hasEFML, reason)

	usesEFGCUModule, reason := r.info.UsesEFGCUModule()
	r.logger.Debugf("Uses efgcu kernel module? %v: %v", usesEFGCUModule, reason)

	if (isTegra && !hasEFML) || usesEFGCUModule {
		return "csv"
	}
*/
	return "legacy"
}
