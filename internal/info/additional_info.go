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

package info

import (
	"fmt"
	"strings"

	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	// "github.com/ENFLAME/go-eflib/pkg/eflib/info"
	// "github.com/ENFLAME/go-eflib/pkg/efml"
)

// additionalInfo allows for the info.Interface to be extened to implement the infoInterface.
type additionalInfo struct {
// EF_CLOSE
/*
	info.Interface
	efmllib   efml.Interface
	devicelib device.Interface
*/
}

// UsesEFGCUModule checks whether the efgcu module is used.
// We use the device name to signal this, since devices that use the efgcu module have their device
// names as:
//
//	GCU 0: Orin (efgcu) (UUID: 54d0709b-558d-5a59-9c65-0c5fc14a21a4)
//
// This function returns true if ALL devices use the efgcu module.
func (i additionalInfo) UsesEFGCUModule() (uses bool, reason string) {
// EF_CLOSE
/*
	// We ensure that this function never panics
	defer func() {
		if err := recover(); err != nil {
			uses = false
			reason = fmt.Sprintf("panic: %v", err)
		}
	}()


	ret := i.efmllib.Init()
	if ret != efml.SUCCESS {
		return false, fmt.Sprintf("failed to initialize efml: %v", ret)
	}
	defer func() {
		_ = i.efmllib.Shutdown()
	}()
*/
	var names []string

// EF_CLOSE
/*
	err := i.devicelib.VisitDevices(func(i int, d device.Device) error {
		name, ret := d.GetName()
		if ret != efml.SUCCESS {
			return fmt.Errorf("device %v: %v", i, ret)
		}
		names = append(names, name)
		return nil
	})
	if err != nil {
		return false, fmt.Sprintf("failed to get device names: %v", err)
	}
*/
	if len(names) == 0 {
		return false, "no devices found"
	}

	for _, name := range names {
		if !strings.Contains(name, "(efgcu)") {
			return false, fmt.Sprintf("device %q does not use efgcu module", name)
		}
	}
	return true, "all devices use efgcu module"
}
