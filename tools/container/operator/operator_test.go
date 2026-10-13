/**
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

package operator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	testCases := []struct {
		setAsDefault           bool
		enflameRuntimeName      string
		root                   string
		expectedDefaultRuntime string
		expectedRuntimes       Runtimes
	}{
		{
			expectedRuntimes: Runtimes{
				"enflame": Runtime{
					name: "enflame",
					Path: "/usr/bin/enflame-container-runtime",
				},
				"enflame-cdi": Runtime{
					name: "enflame-cdi",
					Path: "/usr/bin/enflame-container-runtime.cdi",
				},
				"enflame-legacy": Runtime{
					name: "enflame-legacy",
					Path: "/usr/bin/enflame-container-runtime.legacy",
				},
			},
		},
		{
			setAsDefault:           true,
			expectedDefaultRuntime: "enflame",
			expectedRuntimes: Runtimes{
				"enflame": Runtime{
					name:         "enflame",
					Path:         "/usr/bin/enflame-container-runtime",
					SetAsDefault: true,
				},
				"enflame-cdi": Runtime{
					name: "enflame-cdi",
					Path: "/usr/bin/enflame-container-runtime.cdi",
				},
				"enflame-legacy": Runtime{
					name: "enflame-legacy",
					Path: "/usr/bin/enflame-container-runtime.legacy",
				},
			},
		},
		{
			setAsDefault:           true,
			enflameRuntimeName:      "enflame",
			expectedDefaultRuntime: "enflame",
			expectedRuntimes: Runtimes{
				"enflame": Runtime{
					name:         "enflame",
					Path:         "/usr/bin/enflame-container-runtime",
					SetAsDefault: true,
				},
				"enflame-cdi": Runtime{
					name: "enflame-cdi",
					Path: "/usr/bin/enflame-container-runtime.cdi",
				},
				"enflame-legacy": Runtime{
					name: "enflame-legacy",
					Path: "/usr/bin/enflame-container-runtime.legacy",
				},
			},
		},
		{
			setAsDefault:           true,
			enflameRuntimeName:      "NAME",
			expectedDefaultRuntime: "NAME",
			expectedRuntimes: Runtimes{
				"NAME": Runtime{
					name:         "NAME",
					Path:         "/usr/bin/enflame-container-runtime",
					SetAsDefault: true,
				},
				"enflame-cdi": Runtime{
					name: "enflame-cdi",
					Path: "/usr/bin/enflame-container-runtime.cdi",
				},
				"enflame-legacy": Runtime{
					name: "enflame-legacy",
					Path: "/usr/bin/enflame-container-runtime.legacy",
				},
			},
		},
		{
			setAsDefault:      false,
			enflameRuntimeName: "NAME",
			expectedRuntimes: Runtimes{
				"NAME": Runtime{
					name: "NAME",
					Path: "/usr/bin/enflame-container-runtime",
				},
				"enflame-cdi": Runtime{
					name: "enflame-cdi",
					Path: "/usr/bin/enflame-container-runtime.cdi",
				},
				"enflame-legacy": Runtime{
					name: "enflame-legacy",
					Path: "/usr/bin/enflame-container-runtime.legacy",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			runtimes := GetRuntimes(
				WithEnflameRuntimeName(tc.enflameRuntimeName),
				WithSetAsDefault(tc.setAsDefault),
				WithRoot(tc.root),
			)

			require.EqualValues(t, tc.expectedRuntimes, runtimes)
			require.Equal(t, tc.expectedDefaultRuntime, runtimes.DefaultRuntimeName())
		})
	}
}
