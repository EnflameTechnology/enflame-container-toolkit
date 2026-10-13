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
	"testing"

	// "github.com/ENFLAME/go-eflib/pkg/eflib/device"
	// "github.com/ENFLAME/go-eflib/pkg/efml"
	"github.com/stretchr/testify/require"
)

func TestUsesEFGCUModule(t *testing.T) {
	testCases := []struct {
		description string
		efmllib     efml.Interface
		expected    bool
	}{
		{
			description: "init failure returns false",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.ERROR_LIBRARY_NOT_FOUND
				},
			},
			expected: false,
		},
		{
			description: "no devices returns false",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 0, efml.SUCCESS
				},
			},
			expected: false,
		},
		{
			description: "DeviceGetCount error returns false",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 0, efml.ERROR_UNKNOWN
				},
			},
			expected: false,
		},
		{
			description: "Failure to get device name returns false",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 1, efml.SUCCESS
				},
				DeviceGetHandleByIndexFunc: func(index int) (efml.Device, efml.Return) {
					device := &efml.DeviceMock{
						GetNameFunc: func() (string, efml.Return) {
							return "", efml.ERROR_UNKNOWN
						},
					}
					return device, efml.SUCCESS
				},
			},
			expected: false,
		},
		{
			description: "nested panic returns false",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 1, efml.SUCCESS
				},
				DeviceGetHandleByIndexFunc: func(index int) (efml.Device, efml.Return) {
					device := &efml.DeviceMock{
						GetNameFunc: func() (string, efml.Return) {
							panic("deep panic")
						},
					}
					return device, efml.SUCCESS
				},
			},
			expected: false,
		},
		{
			description: "Single device name with no efgcu",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 1, efml.SUCCESS
				},
				DeviceGetHandleByIndexFunc: func(index int) (efml.Device, efml.Return) {
					device := &efml.DeviceMock{
						GetNameFunc: func() (string, efml.Return) {
							return "ENFLAME A100-SXM4-40GB", efml.SUCCESS
						},
					}
					return device, efml.SUCCESS
				},
			},
			expected: false,
		},
		{
			description: "Single device name with efgcu",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 1, efml.SUCCESS
				},
				DeviceGetHandleByIndexFunc: func(index int) (efml.Device, efml.Return) {
					device := &efml.DeviceMock{
						GetNameFunc: func() (string, efml.Return) {
							return "Orin (efgcu)", efml.SUCCESS
						},
					}
					return device, efml.SUCCESS
				},
			},
			expected: true,
		},
		{
			description: "Multiple device names with no efgcu",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 2, efml.SUCCESS
				},
				DeviceGetHandleByIndexFunc: func(index int) (efml.Device, efml.Return) {
					device := &efml.DeviceMock{
						GetNameFunc: func() (string, efml.Return) {
							return "ENFLAME A100-SXM4-40GB", efml.SUCCESS
						},
					}
					return device, efml.SUCCESS
				},
			},
			expected: false,
		},
		{
			description: "Multiple device names with efgcu",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 2, efml.SUCCESS
				},
				DeviceGetHandleByIndexFunc: func(index int) (efml.Device, efml.Return) {
					device := &efml.DeviceMock{
						GetNameFunc: func() (string, efml.Return) {
							return "Orin (efgcu)", efml.SUCCESS
						},
					}
					return device, efml.SUCCESS
				},
			},
			expected: true,
		},
		{
			description: "Mixed device names",
			efmllib: &efml.InterfaceMock{
				InitFunc: func() efml.Return {
					return efml.SUCCESS
				},
				ShutdownFunc: func() efml.Return {
					return efml.SUCCESS
				},
				DeviceGetCountFunc: func() (int, efml.Return) {
					return 2, efml.SUCCESS
				},
				DeviceGetHandleByIndexFunc: func(index int) (efml.Device, efml.Return) {
					var deviceName string
					if index == 0 {
						deviceName = "ENFLAME A100-SXM4-40GB"
					} else {
						deviceName = "Orin (efgcu)"
					}
					device := &efml.DeviceMock{
						GetNameFunc: func() (string, efml.Return) {
							return deviceName, efml.SUCCESS
						},
					}
					return device, efml.SUCCESS
				},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			sut := additionalInfo{
				efmllib:   tc.efmllib,
				devicelib: device.New(device.WithEfml(tc.efmllib)),
			}

			flag, _ := sut.UsesEFGCUModule()
			require.Equal(t, tc.expected, flag)
		})
	}
}
