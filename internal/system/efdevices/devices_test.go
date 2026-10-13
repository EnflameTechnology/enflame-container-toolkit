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

package efdevices

import (
	"errors"
	"testing"

	testlog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"

	"enflame/enflame-container-toolkit/internal/info/proc/devices"
)

func TestCreateControlDevices(t *testing.T) {
	logger, _ := testlog.NewNullLogger()

	enflameDevices := devices.New(
		devices.WithDeviceToMajor(map[string]int{
			"enflame-frontend": 195,
			"enflame-uvm":      243,
		}),
	)
	enflame550Devices := devices.New(
		devices.WithDeviceToMajor(map[string]int{
			"enflame":     195,
			"enflame-uvm": 243,
		}),
	)

	mknodeError := errors.New("mknode error")

	testCases := []struct {
		description   string
		root          string
		devices       devices.Devices
		mknodeError   error
		expectedError error
		expectedCalls []struct {
			S  string
			N1 int
			N2 int
		}
	}{
		{
			description: "no root specified; pre 550 driver",
			root:        "",
			devices:     enflameDevices,
			mknodeError: nil,
			expectedCalls: []struct {
				S  string
				N1 int
				N2 int
			}{
				{"/dev/enflamectl", 195, 255},
				{"/dev/enflame-modeset", 195, 254},
				{"/dev/enflame-uvm", 243, 0},
				{"/dev/enflame-uvm-tools", 243, 1},
			},
		},
		{
			description: "no root specified; 550 driver",
			root:        "",
			devices:     enflame550Devices,
			mknodeError: nil,
			expectedCalls: []struct {
				S  string
				N1 int
				N2 int
			}{
				{"/dev/enflamectl", 195, 255},
				{"/dev/enflame-modeset", 195, 254},
				{"/dev/enflame-uvm", 243, 0},
				{"/dev/enflame-uvm-tools", 243, 1},
			},
		},
		{
			description: "some root specified",
			root:        "/some/root",
			devices:     enflameDevices,
			mknodeError: nil,
			expectedCalls: []struct {
				S  string
				N1 int
				N2 int
			}{
				{"/some/root/dev/enflamectl", 195, 255},
				{"/some/root/dev/enflame-modeset", 195, 254},
				{"/some/root/dev/enflame-uvm", 243, 0},
				{"/some/root/dev/enflame-uvm-tools", 243, 1},
			},
		},
		{
			description:   "mknod error returns error",
			devices:       enflameDevices,
			mknodeError:   mknodeError,
			expectedError: mknodeError,
			// We expect the first call to this to fail, and the rest to be skipped
			expectedCalls: []struct {
				S  string
				N1 int
				N2 int
			}{
				{"/dev/enflamectl", 195, 255},
			},
		},
		{
			description: "missing major returns error",
			devices: &devices.DevicesMock{
				GetFunc: func(name devices.Name) (devices.Major, bool) {
					return 0, false
				},
			},
			expectedError: errInvalidDeviceNode,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			mknode := &mknoderMock{
				MknodeFunc: func(string, int, int) error {
					return tc.mknodeError
				},
			}

			d, _ := New(
				WithLogger(logger),
				WithDevRoot(tc.root),
				WithDevices(tc.devices),
			)
			d.mknoder = mknode

			err := d.CreateENFLAMEControlDevices()
			require.ErrorIs(t, err, tc.expectedError)
			require.EqualValues(t, tc.expectedCalls, mknode.MknodeCalls())
		})
	}
}
