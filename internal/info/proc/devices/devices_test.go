/*
# Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
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

package devices

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnflameDevices(t *testing.T) {
	perDriverDeviceMaps := map[string]map[string]int{
		"pre550": {
			"enflame-frontend": 195,
			"enflame-eflink":   234,
			"enflame-caps":     235,
			"enflame-uvm":      510,
			"enflame-efswitch": 511,
		},
		"post550": {
			"enflame":          195,
			"enflame-eflink":   234,
			"enflame-caps":     235,
			"enflame-uvm":      510,
			"enflame-efswitch": 511,
		},
	}

	for k, devices := range perDriverDeviceMaps {
		enflameDevices := New(WithDeviceToMajor(devices))
		t.Run(k, func(t *testing.T) {
			// Each of the expected devices needs to exist.
			for name, major := range devices {
				device, exists := enflameDevices.Get(Name(name))
				require.True(t, exists)
				require.Equal(t, device, Major(major))
			}
			// An unexpected device cannot exist
			_, exists := enflameDevices.Get("bogus")
			require.False(t, exists)

			// Regardles of the driver version, the enflame and enflame-frontend
			// names are supported and have the same value.
			enflame, exists := enflameDevices.Get(ENFLAMEGCU)
			require.True(t, exists)
			enflameFrontend, exists := enflameDevices.Get(ENFLAMEFrontend)
			require.True(t, exists)
			require.Equal(t, enflame, enflameFrontend)
		})

	}
}

func TestProcessDeviceFile(t *testing.T) {
	testCases := []struct {
		lines         []string
		expected      devices
		expectedError error
	}{
		{lines: []string{}, expectedError: errNoEnflameDevices},
		{lines: []string{"Not a valid line:"}, expectedError: errNoEnflameDevices},
		{lines: []string{"195 enflame-frontend"}, expected: devices{"enflame-frontend": 195}},
		{lines: []string{"195 enflame"}, expected: devices{"enflame": 195}},
		{lines: []string{"195 enflame-frontend", "235 enflame-caps"}, expected: devices{"enflame-frontend": 195, "enflame-caps": 235}},
		{lines: []string{"  195 enflame-frontend"}, expected: devices{"enflame-frontend": 195}},
		{lines: []string{"Not a valid line:", "", "195 enflame-frontend"}, expected: devices{"enflame-frontend": 195}},
		{lines: []string{"195 not-enflame-frontend"}, expectedError: errNoEnflameDevices},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("testcase %d", i), func(t *testing.T) {
			contents := strings.NewReader(strings.Join(tc.lines, "\n"))
			d, err := enflameDeviceFrom(contents)
			require.ErrorIs(t, err, tc.expectedError)

			if tc.expectedError == nil {
				require.EqualValues(t, tc.expected, d.(devices))
			}

		})
	}
}

func TestProcessDeviceFileLine(t *testing.T) {
	testCases := []struct {
		line  string
		name  string
		major int
		err   bool
	}{
		{"", "", 0, true},
		{"0", "", 0, true},
		{"notint enflame-frontend", "", 0, true},
		{"195 enflame-frontend", "enflame-frontend", 195, false},
		{"   195 enflame-frontend", "enflame-frontend", 195, false},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("testcase %d", i), func(t *testing.T) {
			name, major, err := processProcDeviceLine(tc.line)

			require.Equal(t, tc.name, name)
			require.Equal(t, tc.major, major)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

		})
	}
}
