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

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfigWithCustomConfig(t *testing.T) {
	testDir := t.TempDir()
	t.Setenv(configOverride, testDir)

	filename := filepath.Join(testDir, configFilePath)

	// By default debug is disabled
	contents := []byte("[enflame-container-runtime]\ndebug = \"/enflame-container-toolkit.log\"")

	require.NoError(t, os.MkdirAll(filepath.Dir(filename), 0766))
	require.NoError(t, os.WriteFile(filename, contents, 0600))

	cfg, err := GetConfig()
	require.NoError(t, err)
	require.Equal(t, "/enflame-container-toolkit.log", cfg.ENFLAMEContainerRuntimeConfig.DebugFilePath)
}

func TestGetConfig(t *testing.T) {
	testCases := []struct {
		description     string
		contents        []string
		expectedError   error
		inspectLdconfig bool
		distIdsLike     []string
		expectedConfig  *Config
	}{
		{
			description:     "empty config is default",
			inspectLdconfig: true,
			expectedConfig: &Config{
				AcceptEnvvarUnprivileged:    true,
				SupportedDriverCapabilities: "compat32,compute,display,graphics,ngx,utility,video",
				ENFLAMEContainerCLIConfig: ContainerCLIConfig{
					Root:      "",
					LoadKmods: true,
					Ldconfig:  "WAS_CHECKED",
				},
				ENFLAMEContainerRuntimeConfig: RuntimeConfig{
					DebugFilePath: "/dev/null",
					LogLevel:      "info",
					Runtimes:      []string{"docker-runc", "runc", "crun"},
					Mode:          "auto",
					Modes: modesConfig{
						CSV: csvModeConfig{
							MountSpecPath: "/etc/enflame-container-runtime/host-files-for-container.d",
						},
						CDI: cdiModeConfig{
							DefaultKind:        "enflame.com/gcu",
							AnnotationPrefixes: []string{"cdi.k8s.io/"},
							SpecDirs:           []string{"/etc/cdi", "/var/run/cdi"},
						},
					},
				},
				ENFLAMEContainerRuntimeHookConfig: RuntimeHookConfig{
					Path: "enflame-container-runtime-hook",
				},
				ENFLAMECTKConfig: CTKConfig{
					Path: "enflame-ctk",
				},
			},
		},
		{
			description: "config options set inline",
			contents: []string{
				"accept-enflame-visible-devices-envvar-when-unprivileged = false",
				"supported-driver-capabilities = \"compute,utility\"",
				"enflame-container-cli.root = \"/bar/baz\"",
				"enflame-container-cli.load-kmods = false",
				"enflame-container-cli.ldconfig = \"/foo/bar/ldconfig\"",
				"enflame-container-cli.user = \"foo:bar\"",
				"enflame-container-runtime.debug = \"/foo/bar\"",
				"enflame-container-runtime.discover-mode = \"not-legacy\"",
				"enflame-container-runtime.log-level = \"debug\"",
				"enflame-container-runtime.runtimes = [\"/some/runtime\",]",
				"enflame-container-runtime.mode = \"not-auto\"",
				"enflame-container-runtime.modes.cdi.default-kind = \"example.vendor.com/device\"",
				"enflame-container-runtime.modes.cdi.annotation-prefixes = [\"cdi.k8s.io/\", \"example.vendor.com/\",]",
				"enflame-container-runtime.modes.cdi.spec-dirs = [\"/except/etc/cdi\", \"/not/var/run/cdi\",]",
				"enflame-container-runtime.modes.csv.mount-spec-path = \"/not/etc/enflame-container-runtime/host-files-for-container.d\"",
				"enflame-container-runtime-hook.path = \"/foo/bar/enflame-container-runtime-hook\"",
				"enflame-ctk.path = \"/foo/bar/enflame-ctk\"",
			},
			expectedConfig: &Config{
				AcceptEnvvarUnprivileged:    false,
				SupportedDriverCapabilities: "compute,utility",
				ENFLAMEContainerCLIConfig: ContainerCLIConfig{
					Root:      "/bar/baz",
					LoadKmods: false,
					Ldconfig:  "/foo/bar/ldconfig",
					User:      "foo:bar",
				},
				ENFLAMEContainerRuntimeConfig: RuntimeConfig{
					DebugFilePath: "/foo/bar",
					LogLevel:      "debug",
					Runtimes:      []string{"/some/runtime"},
					Mode:          "not-auto",
					Modes: modesConfig{
						CSV: csvModeConfig{
							MountSpecPath: "/not/etc/enflame-container-runtime/host-files-for-container.d",
						},
						CDI: cdiModeConfig{
							DefaultKind: "example.vendor.com/device",
							AnnotationPrefixes: []string{
								"cdi.k8s.io/",
								"example.vendor.com/",
							},
							SpecDirs: []string{
								"/except/etc/cdi",
								"/not/var/run/cdi",
							},
						},
					},
				},
				ENFLAMEContainerRuntimeHookConfig: RuntimeHookConfig{
					Path: "/foo/bar/enflame-container-runtime-hook",
				},
				ENFLAMECTKConfig: CTKConfig{
					Path: "/foo/bar/enflame-ctk",
				},
			},
		},
		{
			description: "config options set in section",
			contents: []string{
				"accept-enflame-visible-devices-envvar-when-unprivileged = false",
				"supported-driver-capabilities = \"compute,utility\"",
				"[enflame-container-cli]",
				"root = \"/bar/baz\"",
				"load-kmods = false",
				"ldconfig = \"/foo/bar/ldconfig\"",
				"user = \"foo:bar\"",
				"[enflame-container-runtime]",
				"debug = \"/foo/bar\"",
				"discover-mode = \"not-legacy\"",
				"log-level = \"debug\"",
				"runtimes = [\"/some/runtime\",]",
				"mode = \"not-auto\"",
				"[enflame-container-runtime.modes.cdi]",
				"default-kind = \"example.vendor.com/device\"",
				"annotation-prefixes = [\"cdi.k8s.io/\", \"example.vendor.com/\",]",
				"spec-dirs = [\"/except/etc/cdi\", \"/not/var/run/cdi\",]",
				"[enflame-container-runtime.modes.csv]",
				"mount-spec-path = \"/not/etc/enflame-container-runtime/host-files-for-container.d\"",
				"[enflame-container-runtime-hook]",
				"path = \"/foo/bar/enflame-container-runtime-hook\"",
				"[enflame-ctk]",
				"path = \"/foo/bar/enflame-ctk\"",
			},
			expectedConfig: &Config{
				AcceptEnvvarUnprivileged:    false,
				SupportedDriverCapabilities: "compute,utility",
				ENFLAMEContainerCLIConfig: ContainerCLIConfig{
					Root:      "/bar/baz",
					LoadKmods: false,
					Ldconfig:  "/foo/bar/ldconfig",
					User:      "foo:bar",
				},
				ENFLAMEContainerRuntimeConfig: RuntimeConfig{
					DebugFilePath: "/foo/bar",
					LogLevel:      "debug",
					Runtimes:      []string{"/some/runtime"},
					Mode:          "not-auto",
					Modes: modesConfig{
						CSV: csvModeConfig{
							MountSpecPath: "/not/etc/enflame-container-runtime/host-files-for-container.d",
						},
						CDI: cdiModeConfig{
							DefaultKind: "example.vendor.com/device",
							AnnotationPrefixes: []string{
								"cdi.k8s.io/",
								"example.vendor.com/",
							},
							SpecDirs: []string{
								"/except/etc/cdi",
								"/not/var/run/cdi",
							},
						},
					},
				},
				ENFLAMEContainerRuntimeHookConfig: RuntimeHookConfig{
					Path: "/foo/bar/enflame-container-runtime-hook",
				},
				ENFLAMECTKConfig: CTKConfig{
					Path: "/foo/bar/enflame-ctk",
				},
			},
		},
		{
			description:     "suse config",
			distIdsLike:     []string{"suse", "opensuse"},
			inspectLdconfig: true,
			expectedConfig: &Config{
				AcceptEnvvarUnprivileged:    true,
				SupportedDriverCapabilities: "compat32,compute,display,graphics,ngx,utility,video",
				ENFLAMEContainerCLIConfig: ContainerCLIConfig{
					Root:      "",
					LoadKmods: true,
					Ldconfig:  "WAS_CHECKED",
					User:      "root:video",
				},
				ENFLAMEContainerRuntimeConfig: RuntimeConfig{
					DebugFilePath: "/dev/null",
					LogLevel:      "info",
					Runtimes:      []string{"docker-runc", "runc", "crun"},
					Mode:          "auto",
					Modes: modesConfig{
						CSV: csvModeConfig{
							MountSpecPath: "/etc/enflame-container-runtime/host-files-for-container.d",
						},
						CDI: cdiModeConfig{
							DefaultKind:        "enflame.com/gcu",
							AnnotationPrefixes: []string{"cdi.k8s.io/"},
							SpecDirs:           []string{"/etc/cdi", "/var/run/cdi"},
						},
					},
				},
				ENFLAMEContainerRuntimeHookConfig: RuntimeHookConfig{
					Path: "enflame-container-runtime-hook",
				},
				ENFLAMECTKConfig: CTKConfig{
					Path: "enflame-ctk",
				},
			},
		},
		{
			description:     "suse config overrides user",
			distIdsLike:     []string{"suse", "opensuse"},
			inspectLdconfig: true,
			contents: []string{
				"enflame-container-cli.user = \"foo:bar\"",
			},
			expectedConfig: &Config{
				AcceptEnvvarUnprivileged:    true,
				SupportedDriverCapabilities: "compat32,compute,display,graphics,ngx,utility,video",
				ENFLAMEContainerCLIConfig: ContainerCLIConfig{
					Root:      "",
					LoadKmods: true,
					Ldconfig:  "WAS_CHECKED",
					User:      "foo:bar",
				},
				ENFLAMEContainerRuntimeConfig: RuntimeConfig{
					DebugFilePath: "/dev/null",
					LogLevel:      "info",
					Runtimes:      []string{"docker-runc", "runc", "crun"},
					Mode:          "auto",
					Modes: modesConfig{
						CSV: csvModeConfig{
							MountSpecPath: "/etc/enflame-container-runtime/host-files-for-container.d",
						},
						CDI: cdiModeConfig{
							DefaultKind:        "enflame.com/gcu",
							AnnotationPrefixes: []string{"cdi.k8s.io/"},
							SpecDirs:           []string{"/etc/cdi", "/var/run/cdi"},
						},
					},
				},
				ENFLAMEContainerRuntimeHookConfig: RuntimeHookConfig{
					Path: "enflame-container-runtime-hook",
				},
				ENFLAMECTKConfig: CTKConfig{
					Path: "enflame-ctk",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			defer setGetDistIDLikeForTest(tc.distIdsLike)()
			reader := strings.NewReader(strings.Join(tc.contents, "\n"))

			tomlCfg, err := loadConfigTomlFrom(reader)
			if tc.expectedError != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			cfg, err := tomlCfg.Config()
			require.NoError(t, err)

			// We first handle the ldconfig path since this is currently system-dependent.
			if tc.inspectLdconfig {
				ldconfig := cfg.ENFLAMEContainerCLIConfig.Ldconfig
				require.True(t, strings.HasPrefix(ldconfig, "@/sbin/ldconfig"))
				remaining := strings.TrimPrefix(ldconfig, "@/sbin/ldconfig")
				require.True(t, remaining == ".real" || remaining == "")

				cfg.ENFLAMEContainerCLIConfig.Ldconfig = "WAS_CHECKED"
			}

			require.EqualValues(t, tc.expectedConfig, cfg)
		})
	}
}

// setGetDistIDsLikeForTest overrides the distribution IDs that would normally be read from the /etc/os-release file.
func setGetDistIDLikeForTest(ids []string) func() {
	if ids == nil {
		return func() {}
	}
	original := getDistIDLike

	getDistIDLike = func() []string {
		return ids
	}

	return func() {
		getDistIDLike = original
	}
}
