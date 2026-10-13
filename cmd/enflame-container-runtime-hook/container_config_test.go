package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"enflame/enflame-container-toolkit/internal/config/image"
)

func TestGetEnflameConfig(t *testing.T) {
	var tests = []struct {
		description    string
		env            map[string]string
		privileged     bool
		hookConfig     *HookConfig
		expectedConfig *enflameConfig
		expectedPanic  bool
	}{
		{
			description:    "No environment, unprivileged",
			env:            map[string]string{},
			privileged:     false,
			expectedConfig: nil,
		},
		{
			description:    "No environment, privileged",
			env:            map[string]string{},
			privileged:     true,
			expectedConfig: nil,
		},
		{
			description: "Legacy image, no devices, no capabilities, no requirements",
			env: map[string]string{
				envTOPSVersion: "9.0",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "all",
				DriverCapabilities: image.SupportedDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Legacy image, devices 'all', no capabilities, no requirements",
			env: map[string]string{
				envTOPSVersion:      "9.0",
				envEFVisibleDevices: "all",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "all",
				DriverCapabilities: image.SupportedDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Legacy image, devices 'empty', no capabilities, no requirements",
			env: map[string]string{
				envTOPSVersion:      "9.0",
				envEFVisibleDevices: "",
			},
			privileged:     false,
			expectedConfig: nil,
		},
		{
			description: "Legacy image, devices 'void', no capabilities, no requirements",
			env: map[string]string{
				envTOPSVersion:      "9.0",
				envEFVisibleDevices: "void",
			},
			privileged:     false,
			expectedConfig: nil,
		},
		{
			description: "Legacy image, devices 'none', no capabilities, no requirements",
			env: map[string]string{
				envTOPSVersion:      "9.0",
				envEFVisibleDevices: "none",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "",
				DriverCapabilities: image.SupportedDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Legacy image, devices set, no capabilities, no requirements",
			env: map[string]string{
				envTOPSVersion:      "9.0",
				envEFVisibleDevices: "gcu0,gcu1",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: image.SupportedDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Legacy image, devices set, capabilities 'empty', no requirements",
			env: map[string]string{
				envTOPSVersion:          "9.0",
				envEFVisibleDevices:     "gcu0,gcu1",
				envEFDriverCapabilities: "",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Legacy image, devices set, capabilities 'all', no requirements",
			env: map[string]string{
				envTOPSVersion:          "9.0",
				envEFVisibleDevices:     "gcu0,gcu1",
				envEFDriverCapabilities: "all",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: image.SupportedDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Legacy image, devices set, capabilities set, no requirements",
			env: map[string]string{
				envTOPSVersion:          "9.0",
				envEFVisibleDevices:     "gcu0,gcu1",
				envEFDriverCapabilities: "video,display",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: "display,video",
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Legacy image, devices set, capabilities set, requirements set",
			env: map[string]string{
				envTOPSVersion:              "9.0",
				envEFVisibleDevices:         "gcu0,gcu1",
				envEFDriverCapabilities:     "video,display",
				envEFRequirePrefix + "REQ0": "req0=true",
				envEFRequirePrefix + "REQ1": "req1=false",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: "display,video",
				Requirements:       []string{"tops>=9.0", "req0=true", "req1=false"},
			},
		},
		{
			description: "Legacy image, devices set, capabilities set, requirements set, disable requirements",
			env: map[string]string{
				envTOPSVersion:              "9.0",
				envEFVisibleDevices:         "gcu0,gcu1",
				envEFDriverCapabilities:     "video,display",
				envEFRequirePrefix + "REQ0": "req0=true",
				envEFRequirePrefix + "REQ1": "req1=false",
				envEFDisableRequire:         "true",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: "display,video",
				Requirements:       []string{},
			},
		},
		{
			description: "Modern image, no devices, no capabilities, no requirements, no envTOPSVersion",
			env: map[string]string{
				envEFRequireTOPS: "tops>=9.0",
			},
			privileged:     false,
			expectedConfig: nil,
		},
		{
			description: "Modern image, no devices, no capabilities, no requirement, envTOPSVersion set",
			env: map[string]string{
				envTOPSVersion:   "9.0",
				envEFRequireTOPS: "tops>=9.0",
			},
			privileged:     false,
			expectedConfig: nil,
		},
		{
			description: "Modern image, devices 'all', no capabilities, no requirements",
			env: map[string]string{
				envEFRequireTOPS:    "tops>=9.0",
				envEFVisibleDevices: "all",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "all",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Modern image, devices 'empty', no capabilities, no requirements",
			env: map[string]string{
				envEFRequireTOPS:    "tops>=9.0",
				envEFVisibleDevices: "",
			},
			privileged:     false,
			expectedConfig: nil,
		},
		{
			description: "Modern image, devices 'void', no capabilities, no requirements",
			env: map[string]string{
				envEFRequireTOPS:    "tops>=9.0",
				envEFVisibleDevices: "void",
			},
			privileged:     false,
			expectedConfig: nil,
		},
		{
			description: "Modern image, devices 'none', no capabilities, no requirements",
			env: map[string]string{
				envEFRequireTOPS:    "tops>=9.0",
				envEFVisibleDevices: "none",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Modern image, devices set, no capabilities, no requirements",
			env: map[string]string{
				envEFRequireTOPS:    "tops>=9.0",
				envEFVisibleDevices: "gcu0,gcu1",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Modern image, devices set, capabilities 'empty', no requirements",
			env: map[string]string{
				envEFRequireTOPS:        "tops>=9.0",
				envEFVisibleDevices:     "gcu0,gcu1",
				envEFDriverCapabilities: "",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Modern image, devices set, capabilities 'all', no requirements",
			env: map[string]string{
				envEFRequireTOPS:        "tops>=9.0",
				envEFVisibleDevices:     "gcu0,gcu1",
				envEFDriverCapabilities: "all",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: image.SupportedDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Modern image, devices set, capabilities set, no requirements",
			env: map[string]string{
				envEFRequireTOPS:        "tops>=9.0",
				envEFVisibleDevices:     "gcu0,gcu1",
				envEFDriverCapabilities: "video,display",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: "display,video",
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Modern image, devices set, capabilities set, requirements set",
			env: map[string]string{
				envEFRequireTOPS:            "tops>=9.0",
				envEFVisibleDevices:         "gcu0,gcu1",
				envEFDriverCapabilities:     "video,display",
				envEFRequirePrefix + "REQ0": "req0=true",
				envEFRequirePrefix + "REQ1": "req1=false",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: "display,video",
				Requirements:       []string{"tops>=9.0", "req0=true", "req1=false"},
			},
		},
		{
			description: "Modern image, devices set, capabilities set, requirements set, disable requirements",
			env: map[string]string{
				envEFRequireTOPS:            "tops>=9.0",
				envEFVisibleDevices:         "gcu0,gcu1",
				envEFDriverCapabilities:     "video,display",
				envEFRequirePrefix + "REQ0": "req0=true",
				envEFRequirePrefix + "REQ1": "req1=false",
				envEFDisableRequire:         "true",
			},
			privileged: false,
			expectedConfig: &enflameConfig{
				Devices:            "gcu0,gcu1",
				DriverCapabilities: "display,video",
				Requirements:       []string{},
			},
		},
		{
			description: "No tops envs, devices 'all'",
			env: map[string]string{
				envEFVisibleDevices: "all",
			},
			privileged: false,

			expectedConfig: &enflameConfig{
				Devices:            "all",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
				Requirements:       []string{},
			},
		},
		{
			description: "Modern image, devices 'all', migConfig set, privileged",
			env: map[string]string{
				envEFRequireTOPS:      "tops>=9.0",
				envEFVisibleDevices:   "all",
				envEFMigConfigDevices: "mig0,mig1",
			},
			privileged: true,
			expectedConfig: &enflameConfig{
				Devices:            "all",
				MigConfigDevices:   "mig0,mig1",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Modern image, devices 'all', migConfig set, unprivileged",
			env: map[string]string{
				envEFRequireTOPS:      "tops>=9.0",
				envEFVisibleDevices:   "all",
				envEFMigConfigDevices: "mig0,mig1",
			},
			privileged:    false,
			expectedPanic: true,
		},
		{
			description: "Modern image, devices 'all', migMonitor set, privileged",
			env: map[string]string{
				envEFRequireTOPS:       "tops>=9.0",
				envEFVisibleDevices:    "all",
				envEFMigMonitorDevices: "mig0,mig1",
			},
			privileged: true,
			expectedConfig: &enflameConfig{
				Devices:            "all",
				MigMonitorDevices:  "mig0,mig1",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
				Requirements:       []string{"tops>=9.0"},
			},
		},
		{
			description: "Modern image, devices 'all', migMonitor set, unprivileged",
			env: map[string]string{
				envEFRequireTOPS:       "tops>=9.0",
				envEFVisibleDevices:    "all",
				envEFMigMonitorDevices: "mig0,mig1",
			},
			privileged:    false,
			expectedPanic: true,
		},
		{
			description: "Hook config set as driver-capabilities-all",
			env: map[string]string{
				envEFVisibleDevices:     "all",
				envEFDriverCapabilities: "all",
			},
			privileged: true,
			hookConfig: &HookConfig{
				SupportedDriverCapabilities: "video,display",
			},
			expectedConfig: &enflameConfig{
				Devices:            "all",
				DriverCapabilities: "display,video",
			},
		},
		{
			description: "Hook config set, envvar sets driver-capabilities",
			env: map[string]string{
				envEFVisibleDevices:     "all",
				envEFDriverCapabilities: "video,display",
			},
			privileged: true,
			hookConfig: &HookConfig{
				SupportedDriverCapabilities: "video,display,compute,utility",
			},
			expectedConfig: &enflameConfig{
				Devices:            "all",
				DriverCapabilities: "display,video",
			},
		},
		{
			description: "Hook config set, envvar unset sets default driver-capabilities",
			env: map[string]string{
				envEFVisibleDevices: "all",
			},
			privileged: true,
			hookConfig: &HookConfig{
				SupportedDriverCapabilities: "video,display,utility,compute",
			},
			expectedConfig: &enflameConfig{
				Devices:            "all",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
			},
		},
		{
			description: "Hook config set, swarmResource overrides device selection",
			env: map[string]string{
				envEFVisibleDevices:     "all",
				"DOCKER_SWARM_RESOURCE": "GCU1,GCU2",
			},
			privileged: true,
			hookConfig: &HookConfig{
				SwarmResource:               "DOCKER_SWARM_RESOURCE",
				SupportedDriverCapabilities: "video,display,utility,compute",
			},
			expectedConfig: &enflameConfig{
				Devices:            "GCU1,GCU2",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
			},
		},
		{
			description: "Hook config set, comma separated swarmResource is split and overrides device selection",
			env: map[string]string{
				envEFVisibleDevices:     "all",
				"DOCKER_SWARM_RESOURCE": "GCU1,GCU2",
			},
			privileged: true,
			hookConfig: &HookConfig{
				SwarmResource:               "NOT_DOCKER_SWARM_RESOURCE,DOCKER_SWARM_RESOURCE",
				SupportedDriverCapabilities: "video,display,utility,compute",
			},
			expectedConfig: &enflameConfig{
				Devices:            "GCU1,GCU2",
				DriverCapabilities: image.DefaultDriverCapabilities.String(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			image, _ := image.New(
				image.WithEnvMap(tc.env),
			)
			// Wrap the call to getEnflameConfig() in a closure.
			var config *enflameConfig
			getConfig := func() {
				hookConfig := tc.hookConfig
				if hookConfig == nil {
					defaultConfig, _ := getDefaultHookConfig()
					hookConfig = &defaultConfig
				}
				config = getEnflameConfig(hookConfig, image, nil, tc.privileged)
			}

			// For any tests that are expected to panic, make sure they do.
			if tc.expectedPanic {
				require.Panics(t, getConfig)
				return
			}

			// For all other tests, just grab the config
			getConfig()

			// And start comparing the test results to the expected results.
			if tc.expectedConfig == nil {
				require.Nil(t, config, tc.description)
				return
			}

			require.NotNil(t, config, tc.description)

			require.Equal(t, tc.expectedConfig.Devices, config.Devices)
			require.Equal(t, tc.expectedConfig.MigConfigDevices, config.MigConfigDevices)
			require.Equal(t, tc.expectedConfig.MigMonitorDevices, config.MigMonitorDevices)
			require.Equal(t, tc.expectedConfig.DriverCapabilities, config.DriverCapabilities)

			require.ElementsMatch(t, tc.expectedConfig.Requirements, config.Requirements)
		})
	}
}

func TestGetDevicesFromMounts(t *testing.T) {
	var tests = []struct {
		description     string
		mounts          []Mount
		expectedDevices *string
	}{
		{
			description:     "No mounts",
			mounts:          nil,
			expectedDevices: nil,
		},
		{
			description: "Host path is not /dev/null",
			mounts: []Mount{
				{
					Source:      "/not/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU0"),
				},
			},
			expectedDevices: nil,
		},
		{
			description: "Container path is not prefixed by 'root'",
			mounts: []Mount{
				{
					Source:      "/dev/null",
					Destination: filepath.Join("/other/prefix", "GCU0"),
				},
			},
			expectedDevices: nil,
		},
		{
			description: "Container path is only 'root'",
			mounts: []Mount{
				{
					Source:      "/dev/null",
					Destination: deviceListAsVolumeMountsRoot,
				},
			},
			expectedDevices: nil,
		},
		{
			description: "Discover 2 devices",
			mounts: []Mount{
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU0"),
				},
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU1"),
				},
			},
			expectedDevices: &[]string{"GCU0,GCU1"}[0],
		},
		{
			description: "Discover 2 devices with slashes in the name",
			mounts: []Mount{
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU0-MIG0/0/1"),
				},
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU1-MIG0/0/1"),
				},
			},
			expectedDevices: &[]string{"GCU0-MIG0/0/1,GCU1-MIG0/0/1"}[0],
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			devices := getDevicesFromMounts(tc.mounts)
			require.Equal(t, tc.expectedDevices, devices)
		})
	}
}

func TestDeviceListSourcePriority(t *testing.T) {
	var tests = []struct {
		description        string
		mountDevices       []Mount
		envvarDevices      string
		privileged         bool
		acceptUnprivileged bool
		acceptMounts       bool
		expectedDevices    *string
	}{
		{
			description: "Mount devices, unprivileged, no accept unprivileged",
			mountDevices: []Mount{
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU0"),
				},
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU1"),
				},
			},
			envvarDevices:      "GCU2,GCU3",
			privileged:         false,
			acceptUnprivileged: false,
			acceptMounts:       true,
			expectedDevices:    &[]string{"GCU0,GCU1"}[0],
		},
		{
			description:        "No mount devices, unprivileged, no accept unprivileged",
			mountDevices:       nil,
			envvarDevices:      "GCU0,GCU1",
			privileged:         false,
			acceptUnprivileged: false,
			acceptMounts:       true,
			expectedDevices:    nil,
		},
		{
			description:        "No mount devices, privileged, no accept unprivileged",
			mountDevices:       nil,
			envvarDevices:      "GCU0,GCU1",
			privileged:         true,
			acceptUnprivileged: false,
			acceptMounts:       true,
			expectedDevices:    &[]string{"GCU0,GCU1"}[0],
		},
		{
			description:        "No mount devices, unprivileged, accept unprivileged",
			mountDevices:       nil,
			envvarDevices:      "GCU0,GCU1",
			privileged:         false,
			acceptUnprivileged: true,
			acceptMounts:       true,
			expectedDevices:    &[]string{"GCU0,GCU1"}[0],
		},
		{
			description: "Mount devices, unprivileged, accept unprivileged, no accept mounts",
			mountDevices: []Mount{
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU0"),
				},
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU1"),
				},
			},
			envvarDevices:      "GCU2,GCU3",
			privileged:         false,
			acceptUnprivileged: true,
			acceptMounts:       false,
			expectedDevices:    &[]string{"GCU2,GCU3"}[0],
		},
		{
			description: "Mount devices, unprivileged, no accept unprivileged, no accept mounts",
			mountDevices: []Mount{
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU0"),
				},
				{
					Source:      "/dev/null",
					Destination: filepath.Join(deviceListAsVolumeMountsRoot, "GCU1"),
				},
			},
			envvarDevices:      "GCU2,GCU3",
			privileged:         false,
			acceptUnprivileged: false,
			acceptMounts:       false,
			expectedDevices:    nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			// Wrap the call to getDevices() in a closure.
			var devices *string
			getDevices := func() {
				image, _ := image.New(
					image.WithEnvMap(
						map[string]string{
							envEFVisibleDevices: tc.envvarDevices,
						},
					),
				)
				hookConfig, _ := getDefaultHookConfig()
				hookConfig.AcceptEnvvarUnprivileged = tc.acceptUnprivileged
				hookConfig.AcceptDeviceListAsVolumeMounts = tc.acceptMounts
				devices = getDevices(&hookConfig, image, tc.mountDevices, tc.privileged)
			}

			// For all other tests, just grab the devices and check the results
			getDevices()

			require.Equal(t, tc.expectedDevices, devices)
		})
	}
}

func TestGetDevicesFromEnvvar(t *testing.T) {
	all := "all"
	empty := ""
	envDockerResourceGCUs := "DOCKER_RESOURCE_GCUS"
	gcuID := "GCU-12345"
	anotherGCUID := "GCU-67890"
	thirdGCUID := "MIG-12345"

	var tests = []struct {
		description          string
		swarmResourceEnvvars []string
		env                  map[string]string
		expectedDevices      *string
	}{
		{
			description: "empty env returns nil for non-legacy image",
		},
		{
			description: "blank ENFLAME_VISIBLE_DEVICES returns nil for non-legacy image",
			env: map[string]string{
				envEFVisibleDevices: "",
			},
		},
		{
			description: "'void' ENFLAME_VISIBLE_DEVICES returns nil for non-legacy image",
			env: map[string]string{
				envEFVisibleDevices: "void",
			},
		},
		{
			description: "'none' ENFLAME_VISIBLE_DEVICES returns empty for non-legacy image",
			env: map[string]string{
				envEFVisibleDevices: "none",
			},
			expectedDevices: &empty,
		},
		{
			description: "ENFLAME_VISIBLE_DEVICES set returns value for non-legacy image",
			env: map[string]string{
				envEFVisibleDevices: gcuID,
			},
			expectedDevices: &gcuID,
		},
		{
			description: "ENFLAME_VISIBLE_DEVICES set returns value for legacy image",
			env: map[string]string{
				envEFVisibleDevices: gcuID,
				envTOPSVersion:      "legacy",
			},
			expectedDevices: &gcuID,
		},
		{
			description: "empty env returns all for legacy image",
			env: map[string]string{
				envTOPSVersion: "legacy",
			},
			expectedDevices: &all,
		},
		// Add the `DOCKER_RESOURCE_GCUS` envvar and ensure that this is ignored when
		// not enabled
		{
			description: "missing ENFLAME_VISIBLE_DEVICES returns nil for non-legacy image",
			env: map[string]string{
				envDockerResourceGCUs: anotherGCUID,
			},
		},
		{
			description: "blank ENFLAME_VISIBLE_DEVICES returns nil for non-legacy image",
			env: map[string]string{
				envEFVisibleDevices:   "",
				envDockerResourceGCUs: anotherGCUID,
			},
		},
		{
			description: "'void' ENFLAME_VISIBLE_DEVICES returns nil for non-legacy image",
			env: map[string]string{
				envEFVisibleDevices:   "void",
				envDockerResourceGCUs: anotherGCUID,
			},
		},
		{
			description: "'none' ENFLAME_VISIBLE_DEVICES returns empty for non-legacy image",
			env: map[string]string{
				envEFVisibleDevices:   "none",
				envDockerResourceGCUs: anotherGCUID,
			},
			expectedDevices: &empty,
		},
		{
			description: "ENFLAME_VISIBLE_DEVICES set returns value for non-legacy image",
			env: map[string]string{
				envEFVisibleDevices:   gcuID,
				envDockerResourceGCUs: anotherGCUID,
			},
			expectedDevices: &gcuID,
		},
		{
			description: "ENFLAME_VISIBLE_DEVICES set returns value for legacy image",
			env: map[string]string{
				envEFVisibleDevices:   gcuID,
				envDockerResourceGCUs: anotherGCUID,
				envTOPSVersion:        "legacy",
			},
			expectedDevices: &gcuID,
		},
		{
			description: "empty env returns all for legacy image",
			env: map[string]string{
				envDockerResourceGCUs: anotherGCUID,
				envTOPSVersion:        "legacy",
			},
			expectedDevices: &all,
		},
		// Add the `DOCKER_RESOURCE_GCUS` envvar and ensure that this is selected when
		// enabled
		{
			description:          "empty env returns nil for non-legacy image",
			swarmResourceEnvvars: []string{envDockerResourceGCUs},
		},
		{
			description:          "blank DOCKER_RESOURCE_GCUS returns nil for non-legacy image",
			swarmResourceEnvvars: []string{envDockerResourceGCUs},
			env: map[string]string{
				envDockerResourceGCUs: "",
			},
		},
		{
			description:          "'void' DOCKER_RESOURCE_GCUS returns nil for non-legacy image",
			swarmResourceEnvvars: []string{envDockerResourceGCUs},
			env: map[string]string{
				envDockerResourceGCUs: "void",
			},
		},
		{
			description:          "'none' DOCKER_RESOURCE_GCUS returns empty for non-legacy image",
			swarmResourceEnvvars: []string{envDockerResourceGCUs},
			env: map[string]string{
				envDockerResourceGCUs: "none",
			},
			expectedDevices: &empty,
		},
		{
			description:          "DOCKER_RESOURCE_GCUS set returns value for non-legacy image",
			swarmResourceEnvvars: []string{envDockerResourceGCUs},
			env: map[string]string{
				envDockerResourceGCUs: gcuID,
			},
			expectedDevices: &gcuID,
		},
		{
			description:          "DOCKER_RESOURCE_GCUS set returns value for legacy image",
			swarmResourceEnvvars: []string{envDockerResourceGCUs},
			env: map[string]string{
				envDockerResourceGCUs: gcuID,
				envTOPSVersion:        "legacy",
			},
			expectedDevices: &gcuID,
		},
		{
			description:          "DOCKER_RESOURCE_GCUS is selected if present",
			swarmResourceEnvvars: []string{envDockerResourceGCUs},
			env: map[string]string{
				envDockerResourceGCUs: anotherGCUID,
			},
			expectedDevices: &anotherGCUID,
		},
		{
			description:          "DOCKER_RESOURCE_GCUS overrides ENFLAME_VISIBLE_DEVICES if present",
			swarmResourceEnvvars: []string{envDockerResourceGCUs},
			env: map[string]string{
				envEFVisibleDevices:   gcuID,
				envDockerResourceGCUs: anotherGCUID,
			},
			expectedDevices: &anotherGCUID,
		},
		{
			description:          "DOCKER_RESOURCE_GCUS_ADDITIONAL overrides ENFLAME_VISIBLE_DEVICES if present",
			swarmResourceEnvvars: []string{"DOCKER_RESOURCE_GCUS_ADDITIONAL"},
			env: map[string]string{
				envEFVisibleDevices:               gcuID,
				"DOCKER_RESOURCE_GCUS_ADDITIONAL": anotherGCUID,
			},
			expectedDevices: &anotherGCUID,
		},
		{
			description:          "All available swarm resource envvars are selected and override ENFLAME_VISIBLE_DEVICES if present",
			swarmResourceEnvvars: []string{"DOCKER_RESOURCE_GCUS", "DOCKER_RESOURCE_GCUS_ADDITIONAL"},
			env: map[string]string{
				envEFVisibleDevices:               gcuID,
				"DOCKER_RESOURCE_GCUS":            thirdGCUID,
				"DOCKER_RESOURCE_GCUS_ADDITIONAL": anotherGCUID,
			},
			expectedDevices: func() *string {
				result := fmt.Sprintf("%s,%s", thirdGCUID, anotherGCUID)
				return &result
			}(),
		},
		{
			description:          "DOCKER_RESOURCE_GCUS_ADDITIONAL or DOCKER_RESOURCE_GCUS override ENFLAME_VISIBLE_DEVICES if present",
			swarmResourceEnvvars: []string{"DOCKER_RESOURCE_GCUS", "DOCKER_RESOURCE_GCUS_ADDITIONAL"},
			env: map[string]string{
				envEFVisibleDevices:               gcuID,
				"DOCKER_RESOURCE_GCUS_ADDITIONAL": anotherGCUID,
			},
			expectedDevices: &anotherGCUID,
		},
	}

	for i, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			image, _ := image.New(
				image.WithEnvMap(tc.env),
			)
			devices := getDevicesFromEnvvar(image, tc.swarmResourceEnvvars)
			if tc.expectedDevices == nil {
				require.Nil(t, devices, "%d: %v", i, tc)
				return
			}

			require.NotNil(t, devices, "%d: %v", i, tc)
			require.Equal(t, *tc.expectedDevices, *devices, "%d: %v", i, tc)
		})
	}
}

func TestGetDriverCapabilities(t *testing.T) {

	supportedCapabilities := "compute,display,utility,video"

	testCases := []struct {
		description           string
		env                   map[string]string
		legacyImage           bool
		supportedCapabilities string
		expectedPanic         bool
		expectedCapabilities  string
	}{
		{
			description: "Env is set for legacy image",
			env: map[string]string{
				envEFDriverCapabilities: "display,video",
			},
			legacyImage:           true,
			supportedCapabilities: supportedCapabilities,
			expectedCapabilities:  "display,video",
		},
		{
			description: "Env is all for legacy image",
			env: map[string]string{
				envEFDriverCapabilities: "all",
			},
			legacyImage:           true,
			supportedCapabilities: supportedCapabilities,
			expectedCapabilities:  supportedCapabilities,
		},
		{
			description: "Env is empty for legacy image",
			env: map[string]string{
				envEFDriverCapabilities: "",
			},
			legacyImage:           true,
			supportedCapabilities: supportedCapabilities,
			expectedCapabilities:  image.DefaultDriverCapabilities.String(),
		},
		{
			description:           "Env unset for legacy image is 'all'",
			env:                   map[string]string{},
			legacyImage:           true,
			supportedCapabilities: supportedCapabilities,
			expectedCapabilities:  supportedCapabilities,
		},
		{
			description: "Env is set for modern image",
			env: map[string]string{
				envEFDriverCapabilities: "display,video",
			},
			legacyImage:           false,
			supportedCapabilities: supportedCapabilities,
			expectedCapabilities:  "display,video",
		},
		{
			description:           "Env unset for modern image is default",
			env:                   map[string]string{},
			legacyImage:           false,
			supportedCapabilities: supportedCapabilities,
			expectedCapabilities:  image.DefaultDriverCapabilities.String(),
		},
		{
			description: "Env is all for modern image",
			env: map[string]string{
				envEFDriverCapabilities: "all",
			},
			legacyImage:           false,
			supportedCapabilities: supportedCapabilities,
			expectedCapabilities:  supportedCapabilities,
		},
		{
			description: "Env is empty for modern image",
			env: map[string]string{
				envEFDriverCapabilities: "",
			},
			legacyImage:           false,
			supportedCapabilities: supportedCapabilities,
			expectedCapabilities:  image.DefaultDriverCapabilities.String(),
		},
		{
			description: "Invalid capabilities panic",
			env: map[string]string{
				envEFDriverCapabilities: "compute,utility",
			},
			supportedCapabilities: "not-compute,not-utility",
			expectedPanic:         true,
		},
		{
			description:           "Default is restricted for modern image",
			legacyImage:           false,
			supportedCapabilities: "compute",
			expectedCapabilities:  "compute",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			var capabilities string

			c := HookConfig{
				SupportedDriverCapabilities: tc.supportedCapabilities,
			}

			image, _ := image.New(
				image.WithEnvMap(tc.env),
			)
			getDriverCapabilities := func() {
				capabilities = c.getDriverCapabilities(image, tc.legacyImage).String()
			}

			if tc.expectedPanic {
				require.Panics(t, getDriverCapabilities)
				return
			}

			getDriverCapabilities()
			require.EqualValues(t, tc.expectedCapabilities, capabilities)
		})
	}
}
