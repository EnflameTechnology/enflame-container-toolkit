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

package constraints

import (
	"testing"

	testlog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	logger, _ := testlog.NewNullLogger()

	tops := NewVersionProperty("tops", "")

	f := factory{
		logger: logger,
		properties: map[string]Property{
			"tops": tops,
		},
	}

	testCases := []struct {
		description   string
		condition     string
		expectedError bool
		expected      Constraint
	}{
		{
			description: "empty is nil",
			condition:   "",
			expected:    nil,
		},
		{
			description:   "missing operator is invalid",
			condition:     "notvalid",
			expectedError: true,
		},
		{
			description: "invalid property is invalid",
			condition:   "foo=45",
			expected:    nil,
		},
		{
			description:   "tops must be semver",
			condition:     "tops=foo",
			expectedError: true,
			expected:      binary{tops, equal, "foo"},
		},
		{
			description: "tops greater than equal",
			condition:   "tops>=11.6",
			expected:    binary{tops, greaterEqual, "11.6"},
		},
		{
			description: "tops greater than",
			condition:   "tops>11.6",
			expected:    binary{tops, greater, "11.6"},
		},
		{
			description: "tops less than equal",
			condition:   "tops<=11.6",
			expected:    binary{tops, lessEqual, "11.6"},
		},
		{
			description: "tops less than",
			condition:   "tops<11.6",
			expected:    binary{tops, less, "11.6"},
		},
		{
			description: "tops equal",
			condition:   "tops=11.6",
			expected:    binary{tops, equal, "11.6"},
		},
		{
			description: "tops not equal",
			condition:   "tops!=11.6",
			expected:    binary{tops, notEqual, "11.6"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			c, err := f.parse(tc.condition)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, tc.expected, c)
		})
	}
}

func TestNewConstraintFromRequirement(t *testing.T) {
	logger, _ := testlog.NewNullLogger()

	tops := &PropertyMock{}
	arch := &PropertyMock{}

	f := factory{
		logger: logger,
		properties: map[string]Property{
			"tops": tops,
			"arch": arch,
		},
	}

	testCases := []struct {
		description   string
		requirement   string
		expectedError bool
		expected      Constraint
	}{
		{
			description: "empty is nil",
			requirement: "",
			expected:    nil,
		},
		{
			description:   "malformed constraint is invalid",
			requirement:   "notvalid",
			expectedError: true,
			expected:      nil,
		},
		{
			description: "unsupported property is ignored",
			requirement: "tops>=11.6 foo=bar",
			expected:    binary{tops, greaterEqual, "11.6"},
		},
		{
			description: "space-separated is and",
			requirement: "tops>=11.6 arch=5.3",
			expected: and([]Constraint{
				binary{tops, greaterEqual, "11.6"},
				binary{arch, equal, "5.3"},
			}),
		},
		{
			description: "comma-separated is or",
			requirement: "tops>=11.6,arch=5.3",
			expected: or([]Constraint{
				binary{tops, greaterEqual, "11.6"},
				binary{arch, equal, "5.3"},
			}),
		},
		{
			description: "and takes precedence",
			requirement: "tops<13.6 tops>=11.6,arch=5.3",
			expected: or([]Constraint{
				binary{tops, less, "13.6"},
				and([]Constraint{
					binary{tops, greaterEqual, "11.6"},
					binary{arch, equal, "5.3"},
				}),
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			c, err := f.newConstraintFromRequirement(tc.requirement)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.EqualValues(t, tc.expected, c)
		})
	}

}
