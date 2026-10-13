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

package main

import (
	"fmt"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"

	"enflame/enflame-container-toolkit/pkg/config/engine/containerd"
	"enflame/enflame-container-toolkit/tools/container"
)

func TestUpdateV1ConfigDefaultRuntime(t *testing.T) {
	const runtimeDir = "/test/runtime/dir"

	testCases := []struct {
		legacyConfig                 bool
		setAsDefault                 bool
		runtimeName                  string
		expectedDefaultRuntimeName   interface{}
		expectedDefaultRuntimeBinary interface{}
	}{
		{},
		{
			legacyConfig:                 true,
			setAsDefault:                 false,
			expectedDefaultRuntimeName:   nil,
			expectedDefaultRuntimeBinary: nil,
		},
		{
			legacyConfig:                 true,
			setAsDefault:                 true,
			expectedDefaultRuntimeName:   nil,
			expectedDefaultRuntimeBinary: "/test/runtime/dir/enflame-container-runtime",
		},
		{
			legacyConfig:                 true,
			setAsDefault:                 true,
			runtimeName:                  "NAME",
			expectedDefaultRuntimeName:   nil,
			expectedDefaultRuntimeBinary: "/test/runtime/dir/enflame-container-runtime",
		},
		{
			legacyConfig:                 false,
			setAsDefault:                 false,
			expectedDefaultRuntimeName:   nil,
			expectedDefaultRuntimeBinary: nil,
		},
		{
			legacyConfig:                 false,
			setAsDefault:                 true,
			expectedDefaultRuntimeName:   "enflame",
			expectedDefaultRuntimeBinary: nil,
		},
		{
			legacyConfig:                 false,
			setAsDefault:                 true,
			runtimeName:                  "NAME",
			expectedDefaultRuntimeName:   "NAME",
			expectedDefaultRuntimeBinary: nil,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			o := &options{
				Options: container.Options{
					RuntimeName:  tc.runtimeName,
					RuntimeDir:   runtimeDir,
					SetAsDefault: tc.setAsDefault,
				},
				runtimeType:     runtimeType,
				useLegacyConfig: tc.legacyConfig,
			}

			config, err := toml.TreeFromMap(map[string]interface{}{})
			require.NoError(t, err, "%d: %v", i, tc)

			v1 := &containerd.ConfigV1{
				Tree:                  config,
				UseDefaultRuntimeName: !tc.legacyConfig,
				RuntimeType:           runtimeType,
			}

			err = o.UpdateConfig(v1)
			require.NoError(t, err, "%d: %v", i, tc)

			defaultRuntimeName := v1.GetPath([]string{"plugins", "cri", "containerd", "default_runtime_name"})
			require.EqualValues(t, tc.expectedDefaultRuntimeName, defaultRuntimeName, "%d: %v", i, tc)

			defaultRuntime := v1.GetPath([]string{"plugins", "cri", "containerd", "default_runtime"})
			if tc.expectedDefaultRuntimeBinary == nil {
				require.Nil(t, defaultRuntime, "%d: %v", i, tc)
			} else {
				require.NotNil(t, defaultRuntime)

				expected, err := defaultRuntimeTomlConfigV1(tc.expectedDefaultRuntimeBinary.(string))
				require.NoError(t, err, "%d: %v", i, tc)

				configContents, _ := toml.Marshal(defaultRuntime.(*toml.Tree))
				expectedContents, _ := toml.Marshal(expected)

				require.Equal(t, string(expectedContents), string(configContents), "%d: %v: %v", i, tc)
			}

		})
	}
}

func TestUpdateV1Config(t *testing.T) {
	const runtimeDir = "/test/runtime/dir"

	testCases := []struct {
		runtimeName    string
		expectedConfig map[string]interface{}
	}{
		{
			runtimeName: "enflame",
			expectedConfig: map[string]interface{}{
				"version": int64(1),
				"plugins": map[string]interface{}{
					"cri": map[string]interface{}{
						"containerd": map[string]interface{}{
							"runtimes": map[string]interface{}{
								"enflame": map[string]interface{}{
									"runtime_type":                    "runtime_type",
									"runtime_root":                    "",
									"runtime_engine":                  "",
									"privileged_without_host_devices": false,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"BinaryName": "/test/runtime/dir/enflame-container-runtime",
										"Runtime":    "/test/runtime/dir/enflame-container-runtime",
									},
								},
								"enflame-cdi": map[string]interface{}{
									"runtime_type":                    "runtime_type",
									"runtime_root":                    "",
									"runtime_engine":                  "",
									"privileged_without_host_devices": false,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"BinaryName": "/test/runtime/dir/enflame-container-runtime.cdi",
										"Runtime":    "/test/runtime/dir/enflame-container-runtime.cdi",
									},
								},
								"enflame-legacy": map[string]interface{}{
									"runtime_type":                    "runtime_type",
									"runtime_root":                    "",
									"runtime_engine":                  "",
									"privileged_without_host_devices": false,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"BinaryName": "/test/runtime/dir/enflame-container-runtime.legacy",
										"Runtime":    "/test/runtime/dir/enflame-container-runtime.legacy",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			runtimeName: "NAME",
			expectedConfig: map[string]interface{}{
				"version": int64(1),
				"plugins": map[string]interface{}{
					"cri": map[string]interface{}{
						"containerd": map[string]interface{}{
							"runtimes": map[string]interface{}{
								"NAME": map[string]interface{}{
									"runtime_type":                    "runtime_type",
									"runtime_root":                    "",
									"runtime_engine":                  "",
									"privileged_without_host_devices": false,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"BinaryName": "/test/runtime/dir/enflame-container-runtime",
										"Runtime":    "/test/runtime/dir/enflame-container-runtime",
									},
								},
								"enflame-cdi": map[string]interface{}{
									"runtime_type":                    "runtime_type",
									"runtime_root":                    "",
									"runtime_engine":                  "",
									"privileged_without_host_devices": false,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"BinaryName": "/test/runtime/dir/enflame-container-runtime.cdi",
										"Runtime":    "/test/runtime/dir/enflame-container-runtime.cdi",
									},
								},
								"enflame-legacy": map[string]interface{}{
									"runtime_type":                    "runtime_type",
									"runtime_root":                    "",
									"runtime_engine":                  "",
									"privileged_without_host_devices": false,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"BinaryName": "/test/runtime/dir/enflame-container-runtime.legacy",
										"Runtime":    "/test/runtime/dir/enflame-container-runtime.legacy",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			o := &options{
				Options: container.Options{
					RuntimeName: tc.runtimeName,
					RuntimeDir:  runtimeDir,
				},
			}

			config, err := toml.TreeFromMap(map[string]interface{}{})
			require.NoError(t, err)

			v1 := &containerd.ConfigV1{
				Tree:                  config,
				UseDefaultRuntimeName: true,
				RuntimeType:           runtimeType,
				ContainerAnnotations:  []string{"cdi.k8s.io/*"},
			}

			err = o.UpdateConfig(v1)
			require.NoError(t, err)

			expected, err := toml.TreeFromMap(tc.expectedConfig)
			require.NoError(t, err)

			require.Equal(t, expected.String(), config.String())
		})
	}
}

func TestUpdateV1ConfigWithRuncPresent(t *testing.T) {
	const runtimeDir = "/test/runtime/dir"

	testCases := []struct {
		runtimeName    string
		expectedConfig map[string]interface{}
	}{
		{
			runtimeName: "enflame",
			expectedConfig: map[string]interface{}{
				"version": int64(1),
				"plugins": map[string]interface{}{
					"cri": map[string]interface{}{
						"containerd": map[string]interface{}{
							"runtimes": map[string]interface{}{
								"runc": map[string]interface{}{
									"runtime_type":                    "runc_runtime_type",
									"runtime_root":                    "runc_runtime_root",
									"runtime_engine":                  "runc_runtime_engine",
									"privileged_without_host_devices": true,
									"options": map[string]interface{}{
										"runc-option": "value",
										"BinaryName":  "/runc-binary",
									},
								},
								"enflame": map[string]interface{}{
									"runtime_type":                    "runc_runtime_type",
									"runtime_root":                    "runc_runtime_root",
									"runtime_engine":                  "runc_runtime_engine",
									"privileged_without_host_devices": true,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"runc-option": "value",
										"BinaryName":  "/test/runtime/dir/enflame-container-runtime",
										"Runtime":     "/test/runtime/dir/enflame-container-runtime",
									},
								},
								"enflame-cdi": map[string]interface{}{
									"runtime_type":                    "runc_runtime_type",
									"runtime_root":                    "runc_runtime_root",
									"runtime_engine":                  "runc_runtime_engine",
									"privileged_without_host_devices": true,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"runc-option": "value",
										"BinaryName":  "/test/runtime/dir/enflame-container-runtime.cdi",
										"Runtime":     "/test/runtime/dir/enflame-container-runtime.cdi",
									},
								},
								"enflame-legacy": map[string]interface{}{
									"runtime_type":                    "runc_runtime_type",
									"runtime_root":                    "runc_runtime_root",
									"runtime_engine":                  "runc_runtime_engine",
									"privileged_without_host_devices": true,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"runc-option": "value",
										"BinaryName":  "/test/runtime/dir/enflame-container-runtime.legacy",
										"Runtime":     "/test/runtime/dir/enflame-container-runtime.legacy",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			runtimeName: "NAME",
			expectedConfig: map[string]interface{}{
				"version": int64(1),
				"plugins": map[string]interface{}{
					"cri": map[string]interface{}{
						"containerd": map[string]interface{}{
							"runtimes": map[string]interface{}{
								"runc": map[string]interface{}{
									"runtime_type":                    "runc_runtime_type",
									"runtime_root":                    "runc_runtime_root",
									"runtime_engine":                  "runc_runtime_engine",
									"privileged_without_host_devices": true,
									"options": map[string]interface{}{
										"runc-option": "value",
										"BinaryName":  "/runc-binary",
									},
								},
								"NAME": map[string]interface{}{
									"runtime_type":                    "runc_runtime_type",
									"runtime_root":                    "runc_runtime_root",
									"runtime_engine":                  "runc_runtime_engine",
									"privileged_without_host_devices": true,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"runc-option": "value",
										"BinaryName":  "/test/runtime/dir/enflame-container-runtime",
										"Runtime":     "/test/runtime/dir/enflame-container-runtime",
									},
								},
								"enflame-cdi": map[string]interface{}{
									"runtime_type":                    "runc_runtime_type",
									"runtime_root":                    "runc_runtime_root",
									"runtime_engine":                  "runc_runtime_engine",
									"privileged_without_host_devices": true,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"runc-option": "value",
										"BinaryName":  "/test/runtime/dir/enflame-container-runtime.cdi",
										"Runtime":     "/test/runtime/dir/enflame-container-runtime.cdi",
									},
								},
								"enflame-legacy": map[string]interface{}{
									"runtime_type":                    "runc_runtime_type",
									"runtime_root":                    "runc_runtime_root",
									"runtime_engine":                  "runc_runtime_engine",
									"privileged_without_host_devices": true,
									"container_annotations":           []string{"cdi.k8s.io/*"},
									"options": map[string]interface{}{
										"runc-option": "value",
										"BinaryName":  "/test/runtime/dir/enflame-container-runtime.legacy",
										"Runtime":     "/test/runtime/dir/enflame-container-runtime.legacy",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			o := &options{
				Options: container.Options{
					RuntimeName: tc.runtimeName,
					RuntimeDir:  runtimeDir,
				},
			}

			config, err := toml.TreeFromMap(runcConfigMapV1("/runc-binary"))
			require.NoError(t, err)

			v1 := &containerd.ConfigV1{
				Tree:                  config,
				UseDefaultRuntimeName: true,
				RuntimeType:           runtimeType,
				ContainerAnnotations:  []string{"cdi.k8s.io/*"},
			}

			err = o.UpdateConfig(v1)
			require.NoError(t, err)

			expected, err := toml.TreeFromMap(tc.expectedConfig)
			require.NoError(t, err)

			require.Equal(t, expected.String(), config.String())
		})
	}
}

func TestRevertV1Config(t *testing.T) {
	testCases := []struct {
		config map[string]interface {
		}
		expected map[string]interface{}
	}{
		{},
		{
			config: map[string]interface{}{
				"version": int64(1),
			},
		},
		{
			config: map[string]interface{}{
				"version": int64(1),
				"plugins": map[string]interface{}{
					"cri": map[string]interface{}{
						"containerd": map[string]interface{}{
							"runtimes": map[string]interface{}{
								"enflame":        runtimeMapV1("/test/runtime/dir/enflame-container-runtime"),
								"enflame-cdi":    runtimeMapV1("/test/runtime/dir/enflame-container-runtime.cdi"),
								"enflame-legacy": runtimeMapV1("/test/runtime/dir/enflame-container-runtime.legacy"),
							},
						},
					},
				},
			},
		},
		{
			config: map[string]interface{}{
				"version": int64(1),
				"plugins": map[string]interface{}{
					"cri": map[string]interface{}{
						"containerd": map[string]interface{}{
							"runtimes": map[string]interface{}{
								"enflame":        runtimeMapV1("/test/runtime/dir/enflame-container-runtime"),
								"enflame-cdi":    runtimeMapV1("/test/runtime/dir/enflame-container-runtime.cdi"),
								"enflame-legacy": runtimeMapV1("/test/runtime/dir/enflame-container-runtime.legacy"),
							},
							"default_runtime":      defaultRuntimeV1("/test/runtime/dir/enflame-container-runtime"),
							"default_runtime_name": "enflame",
						},
					},
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			o := &options{
				Options: container.Options{
					RuntimeName: "enflame",
				},
			}

			config, err := toml.TreeFromMap(tc.config)
			require.NoError(t, err, "%d: %v", i, tc)

			expected, err := toml.TreeFromMap(tc.expected)
			require.NoError(t, err, "%d: %v", i, tc)

			v1 := &containerd.ConfigV1{
				Tree:                  config,
				UseDefaultRuntimeName: true,
				RuntimeType:           runtimeType,
			}

			err = o.RevertConfig(v1)
			require.NoError(t, err, "%d: %v", i, tc)

			configContents, _ := toml.Marshal(config)
			expectedContents, _ := toml.Marshal(expected)

			require.Equal(t, string(expectedContents), string(configContents), "%d: %v", i, tc)
		})
	}
}

func defaultRuntimeTomlConfigV1(binary string) (*toml.Tree, error) {
	return toml.TreeFromMap(defaultRuntimeV1(binary))
}

func defaultRuntimeV1(binary string) map[string]interface{} {
	return map[string]interface{}{
		"runtime_type":                    runtimeType,
		"runtime_root":                    "",
		"runtime_engine":                  "",
		"privileged_without_host_devices": false,
		"options": map[string]interface{}{
			"BinaryName": binary,
			"Runtime":    binary,
		},
	}
}

func runtimeMapV1(binary string) map[string]interface{} {
	return map[string]interface{}{
		"runtime_type":                    runtimeType,
		"runtime_root":                    "",
		"runtime_engine":                  "",
		"privileged_without_host_devices": false,
		"options": map[string]interface{}{
			"BinaryName": binary,
			"Runtime":    binary,
		},
	}
}

func runcConfigMapV1(binary string) map[string]interface{} {
	return map[string]interface{}{
		"plugins": map[string]interface{}{
			"cri": map[string]interface{}{
				"containerd": map[string]interface{}{
					"runtimes": map[string]interface{}{
						"runc": map[string]interface{}{
							"runtime_type":                    "runc_runtime_type",
							"runtime_root":                    "runc_runtime_root",
							"runtime_engine":                  "runc_runtime_engine",
							"privileged_without_host_devices": true,
							"options": map[string]interface{}{
								"runc-option": "value",
								"BinaryName":  binary,
							},
						},
					},
				},
			},
		},
	}
}
