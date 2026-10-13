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
	"testing"

	"github.com/opencontainers/runtime-spec/specs-go"
	testlog "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"

	"enflame/enflame-container-toolkit/internal/config/image"
)

func TestResolveAutoMode(t *testing.T) {
	logger, _ := testlog.NewNullLogger()

	testCases := []struct {
		description  string
		mode         string
		expectedMode string
		info         map[string]bool
		envmap       map[string]string
		mounts       []string
	}{
		{
			description:  "non-auto resolves to input",
			mode:         "not-auto",
			expectedMode: "not-auto",
		},
		{
			description:  "no info defaults to legacy",
			mode:         "auto",
			info:         map[string]bool{},
			expectedMode: "legacy",
		},
		{
			description: "non-efml, non-tegra, efgcu resolves to csv",
			mode:        "auto",
			info: map[string]bool{
				"efml":  false,
				"tegra": false,
				"efgcu": true,
			},
			expectedMode: "csv",
		},
		{
			description: "non-efml, tegra, non-efgcu resolves to csv",
			mode:        "auto",
			info: map[string]bool{
				"efml":  false,
				"tegra": true,
				"efgcu": false,
			},
			expectedMode: "csv",
		},
		{
			description: "non-efml, tegra, efgcu resolves to csv",
			mode:        "auto",
			info: map[string]bool{
				"efml":  false,
				"tegra": true,
				"efgcu": true,
			},
			expectedMode: "csv",
		},
		{
			description: "efml, non-tegra, non-efgcu resolves to legacy",
			mode:        "auto",
			info: map[string]bool{
				"efml":  true,
				"tegra": false,
				"efgcu": false,
			},
			expectedMode: "legacy",
		},
		{
			description: "efml, non-tegra, efgcu resolves to csv",
			mode:        "auto",
			info: map[string]bool{
				"efml":  true,
				"tegra": false,
				"efgcu": true,
			},
			expectedMode: "csv",
		},
		{
			description: "efml, tegra, non-efgcu resolves to legacy",
			mode:        "auto",
			info: map[string]bool{
				"efml":  true,
				"tegra": true,
				"efgcu": false,
			},
			expectedMode: "legacy",
		},
		{
			description: "efml, tegra, efgcu resolves to csv",
			mode:        "auto",
			info: map[string]bool{
				"efml":  true,
				"tegra": true,
				"efgcu": true,
			},
			expectedMode: "csv",
		},
		{
			description:  "cdi devices resolves to cdi",
			mode:         "auto",
			expectedMode: "cdi",
			envmap: map[string]string{
				"ENFLAME_VISIBLE_DEVICES": "enflame.com/gcu=all",
			},
		},
		{
			description:  "multiple cdi devices resolves to cdi",
			mode:         "auto",
			expectedMode: "cdi",
			envmap: map[string]string{
				"ENFLAME_VISIBLE_DEVICES": "enflame.com/gcu=0,enflame.com/gcu=1",
			},
		},
		{
			description: "at least one non-cdi device resolves to legacy",
			mode:        "auto",
			envmap: map[string]string{
				"ENFLAME_VISIBLE_DEVICES": "enflame.com/gcu=0,0",
			},
			info: map[string]bool{
				"efml":  true,
				"tegra": false,
				"efgcu": false,
			},
			expectedMode: "legacy",
		},
		{
			description: "at least one non-cdi device resolves to csv",
			mode:        "auto",
			envmap: map[string]string{
				"ENFLAME_VISIBLE_DEVICES": "enflame.com/gcu=0,0",
			},
			info: map[string]bool{
				"efml":  false,
				"tegra": true,
				"efgcu": false,
			},
			expectedMode: "csv",
		},
		{
			description: "cdi mount devices resolves to CDI",
			mode:        "auto",
			mounts: []string{
				"/var/run/enflame-container-devices/cdi/enflame.com/gcu/0",
			},
			expectedMode: "cdi",
		},
		{
			description: "cdi mount and non-CDI devices resolves to legacy",
			mode:        "auto",
			mounts: []string{
				"/var/run/enflame-container-devices/cdi/enflame.com/gcu/0",
				"/var/run/enflame-container-devices/all",
			},
			info: map[string]bool{
				"efml":  true,
				"tegra": false,
				"efgcu": false,
			},
			expectedMode: "legacy",
		},
		{
			description: "cdi mount and non-CDI envvar resolves to legacy",
			mode:        "auto",
			envmap: map[string]string{
				"ENFLAME_VISIBLE_DEVICES": "0",
			},
			mounts: []string{
				"/var/run/enflame-container-devices/cdi/enflame.com/gcu/0",
			},
			info: map[string]bool{
				"efml":  true,
				"tegra": false,
				"efgcu": false,
			},
			expectedMode: "legacy",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			info := &infoInterfaceMock{
				HasEfmlFunc: func() (bool, string) {
					return tc.info["efml"], "efml"
				},
				IsTegraSystemFunc: func() (bool, string) {
					return tc.info["tegra"], "tegra"
				},
				UsesEFGCUModuleFunc: func() (bool, string) {
					return tc.info["efgcu"], "efgcu"
				},
			}

			r := resolver{
				logger: logger,
				info:   info,
			}

			var mounts []specs.Mount
			for _, d := range tc.mounts {
				mount := specs.Mount{
					Source:      "/dev/null",
					Destination: d,
				}
				mounts = append(mounts, mount)
			}
			image, _ := image.New(
				image.WithEnvMap(tc.envmap),
				image.WithMounts(mounts),
			)
			mode := r.resolveMode(tc.mode, image)
			require.EqualValues(t, tc.expectedMode, mode)
		})
	}
}
