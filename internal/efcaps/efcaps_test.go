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

package efcaps

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessMinorsFile(t *testing.T) {
	testCases := []struct {
		lines    []string
		expected MigCaps
	}{
		{[]string{}, MigCaps{}},
		{[]string{"invalidLine"}, MigCaps{}},
		{[]string{"config 1"}, MigCaps{"config": 1}},
		{[]string{"gcu0/gi0/ci0/access 4"}, MigCaps{"gcu0/gi0/ci0/access": 4}},
		{[]string{"config 1", "invalidLine"}, MigCaps{"config": 1}},
		{[]string{"config 1", "gcu0/gi0/ci0/access 4"}, MigCaps{"config": 1, "gcu0/gi0/ci0/access": 4}},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("testcase %d", i), func(t *testing.T) {
			contents := strings.NewReader(strings.Join(tc.lines, "\n"))
			d := processMinorsFile(contents)
			require.Equal(t, tc.expected, d)
		})
	}
}

func TestProcessMigMinorsLine(t *testing.T) {
	testCases := []struct {
		line  string
		cap   MigCap
		minor MigMinor
		err   bool
	}{
		{"config 1", "config", 1, false},
		{"monitor 2", "monitor", 2, false},
		{"gcu0/gi0/access 3", "gcu0/gi0/access", 3, false},
		{"gcu0/gi0/ci0/access 4", "gcu0/gi0/ci0/access", 4, false},
		{"notconfig 99", "", 0, true},
		{"config notanint", "", 0, true},
		{"", "", 0, true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("testcase %d", i), func(t *testing.T) {
			cap, minor, err := processMigMinorsLine(tc.line)

			require.Equal(t, tc.cap, cap)
			require.Equal(t, tc.minor, minor)
			if tc.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMigCapProcPaths(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"config", "/proc/driver/enflame/capabilities/mig/config"},
		{"monitor", "/proc/driver/enflame/capabilities/mig/monitor"},
		{"gcu0/gi0/access", "/proc/driver/enflame/capabilities/gcu0/mig/gi0/access"},
		{"gcu0/gi0/ci0/access", "/proc/driver/enflame/capabilities/gcu0/mig/gi0/ci0/access"},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("testcase %d", i), func(t *testing.T) {
			m := MigCap(tc.input)
			require.Equal(t, tc.expected, m.ProcPath())
		})
	}
}

func TestMigMinorDevicePath(t *testing.T) {
	m := MigMinor(0)
	require.Equal(t, "/dev/enflame-caps/enflame-cap0", m.DevicePath())
}
