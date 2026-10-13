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
	"testing"

	testlog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestResolveMode(t *testing.T) {
	logger, _ := testlog.NewNullLogger()

	testCases := []struct {
		mode      string
		isTegra   bool
		hasDXCore bool
		hasEFML   bool
		expected  string
	}{
		{
			mode:      "auto",
			hasDXCore: true,
			expected:  "wsl",
		},
		{
			mode:      "auto",
			hasDXCore: false,
			isTegra:   true,
			hasEFML:   false,
			expected:  "csv",
		},
		{
			mode:      "auto",
			hasDXCore: false,
			isTegra:   false,
			hasEFML:   false,
			expected:  "efml",
		},
		{
			mode:      "auto",
			hasDXCore: false,
			isTegra:   true,
			hasEFML:   true,
			expected:  "efml",
		},
		{
			mode:      "auto",
			hasDXCore: false,
			isTegra:   false,
			expected:  "efml",
		},
		{
			mode:      "efml",
			hasDXCore: true,
			isTegra:   true,
			expected:  "efml",
		},
		{
			mode:      "wsl",
			hasDXCore: false,
			expected:  "wsl",
		},
		{
			mode:      "not-auto",
			hasDXCore: true,
			expected:  "not-auto",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			l := efcdilib{
				logger:  logger,
				mode:    tc.mode,
				infolib: infoMock{hasDXCore: tc.hasDXCore, isTegra: tc.isTegra, hasEFML: tc.hasEFML},
			}

			require.Equal(t, tc.expected, l.resolveMode())
		})
	}
}

type infoMock struct {
	hasDXCore bool
	isTegra   bool
	hasEFML   bool
}

func (i infoMock) HasDXCore() (bool, string) {
	return i.hasDXCore, ""
}

func (i infoMock) HasEfml() (bool, string) {
	return i.hasEFML, ""
}

func (i infoMock) IsTegraSystem() (bool, string) {
	return i.isTegra, ""
}
