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
	"testing"

	// "github.com/ENFLAME/go-eflib/pkg/efml"
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	testCases := []struct {
		description   string
		efml          efmlUUIDer
		expectedError error
		expecteUUID   string
	}{
		{
			description:   "empty UUIDer returns error",
			expectedError: errUUIDUnsupported,
			expecteUUID:   "",
		},
		{
			description: "efmlUUIDer returns UUID",
			efml: &efmlUUIDerMock{
				GetUUIDFunc: func() (string, efml.Return) {
					return "SOME_UUID", efml.SUCCESS
				},
			},
			expectedError: nil,
			expecteUUID:   "SOME_UUID",
		},
		{
			description: "efmlUUIDer returns error",
			efml: &efmlUUIDerMock{
				GetUUIDFunc: func() (string, efml.Return) {
					return "SOME_UUID", efml.ERROR_UNKNOWN
				},
			},
			expectedError: efml.ERROR_UNKNOWN,
			expecteUUID:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			uuid, err := convert{tc.efml}.GetUUID()
			require.ErrorIs(t, err, tc.expectedError)
			require.Equal(t, tc.expecteUUID, uuid)
		})
	}
}
