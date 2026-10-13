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
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"enflame/enflame-container-toolkit/internal/discover"

	testlog "github.com/sirupsen/logrus/hooks/test"
)

func TestEnflameSMISymlinkHook(t *testing.T) {
	logger, _ := testlog.NewNullLogger()

	errMounts := errors.New("mounts error")

	testCases := []struct {
		description   string
		mounts        discover.Discover
		expectedError error
		expectedHooks []discover.Hook
	}{
		{
			description: "mounts error is returned",
			mounts: &discover.DiscoverMock{
				MountsFunc: func() ([]discover.Mount, error) {
					return nil, errMounts
				},
			},
			expectedError: errMounts,
		},
		{
			description: "no mounts returns no hooks",
			mounts: &discover.DiscoverMock{
				MountsFunc: func() ([]discover.Mount, error) {
					return nil, nil
				},
			},
		},
		{
			description: "no enflame-smi returns no hooks",
			mounts: &discover.DiscoverMock{
				MountsFunc: func() ([]discover.Mount, error) {
					mounts := []discover.Mount{
						{Path: "/not-enflame-smi"},
						{Path: "/also-not-enflame-smi"},
					}
					return mounts, nil
				},
			},
		},
		{
			description: "enflame-smi must be in path",
			mounts: &discover.DiscoverMock{
				MountsFunc: func() ([]discover.Mount, error) {
					mounts := []discover.Mount{
						{Path: "/not-enflame-smi", HostPath: "enflame-smi"},
						{Path: "/also-not-enflame-smi", HostPath: "not-enflame-smi"},
					}
					return mounts, nil
				},
			},
		},
		{
			description: "enflame-smi returns hook",
			mounts: &discover.DiscoverMock{
				MountsFunc: func() ([]discover.Mount, error) {
					mounts := []discover.Mount{
						{Path: "enflame-smi"},
					}
					return mounts, nil
				},
			},
			expectedHooks: []discover.Hook{
				{
					Lifecycle: "createContainer",
					Path:      "enflame-ctk",
					Args: []string{"enflame-ctk", "hook", "create-symlinks",
						"--link", "enflame-smi::/usr/bin/enflame-smi"},
				},
			},
		},
		{
			description: "checks basename",
			mounts: &discover.DiscoverMock{
				MountsFunc: func() ([]discover.Mount, error) {
					mounts := []discover.Mount{
						{Path: "/some/path/enflame-smi"},
						{Path: "/enflame-smi/but-not"},
					}
					return mounts, nil
				},
			},
			expectedHooks: []discover.Hook{
				{
					Lifecycle: "createContainer",
					Path:      "enflame-ctk",
					Args: []string{"enflame-ctk", "hook", "create-symlinks",
						"--link", "/some/path/enflame-smi::/usr/bin/enflame-smi"},
				},
			},
		},
		{
			description: "returns first match",
			mounts: &discover.DiscoverMock{
				MountsFunc: func() ([]discover.Mount, error) {
					mounts := []discover.Mount{
						{Path: "/some/path/enflame-smi"},
						{Path: "/another/path/enflame-smi"},
					}
					return mounts, nil
				},
			},
			expectedHooks: []discover.Hook{
				{
					Lifecycle: "createContainer",
					Path:      "enflame-ctk",
					Args: []string{"enflame-ctk", "hook", "create-symlinks",
						"--link", "/some/path/enflame-smi::/usr/bin/enflame-smi"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			m := enflameSMISimlinkHook{
				logger:        logger,
				mountsFrom:    tc.mounts,
				enflameCTKPath: "enflame-ctk",
			}

			devices, err := m.Devices()
			require.NoError(t, err)
			require.Empty(t, devices)

			mounts, err := m.Mounts()
			require.NoError(t, err)
			require.Empty(t, mounts)

			hooks, err := m.Hooks()
			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expectedHooks, hooks)
		})
	}
}
