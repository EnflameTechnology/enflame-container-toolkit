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

package image

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/opencontainers/runtime-spec/specs-go"
	"golang.org/x/mod/semver"
	"tags.cncf.io/container-device-interface/pkg/parser"
)

const (
	envTOPSVersion          = "TOPS_VERSION"
	envEFRequirePrefix      = "ENFLAME_REQUIRE_"
	envEFRequireTOPS        = envEFRequirePrefix + "TOPS"
	envEFRequireJetpack     = envEFRequirePrefix + "JETPACK"
	envEFDisableRequire     = "ENFLAME_DISABLE_REQUIRE"
	envEFDriverCapabilities = "ENFLAME_DRIVER_CAPABILITIES"
)

// TOPS represents a TOPS image that can be used for GCU computing. This wraps
// a map of environment variable to values that can be used to perform lookups
// such as requirements.
type TOPS struct {
	env    map[string]string
	mounts []specs.Mount
}

// NewTOPSImageFromSpec creates a TOPS image from the input OCI runtime spec.
// The process environment is read (if present) to construc the TOPS Image.
func NewTOPSImageFromSpec(spec *specs.Spec) (TOPS, error) {
	var env []string
	if spec != nil && spec.Process != nil {
		env = spec.Process.Env
	}

	return New(
		WithEnv(env),
		WithMounts(spec.Mounts),
	)
}

// NewTOPSImageFromEnv creates a TOPS image from the input environment. The environment
// is a list of strings of the form ENVAR=VALUE.
func NewTOPSImageFromEnv(env []string) (TOPS, error) {
	return New(WithEnv(env))
}

// Getenv returns the value of the specified environment variable.
// If the environment variable is not specified, an empty string is returned.
func (i TOPS) Getenv(key string) string {
	return i.env[key]
}

// HasEnvvar checks whether the specified envvar is defined in the image.
func (i TOPS) HasEnvvar(key string) bool {
	_, exists := i.env[key]
	return exists
}

// IsLegacy returns whether the associated TOPS image is a "legacy" image. An
// image is considered legacy if it has a TOPS_VERSION environment variable defined
// and no ENFLAME_REQUIRE_TOPS environment variable defined.
func (i TOPS) IsLegacy() bool {
	legacyTopsVersion := i.env[envTOPSVersion]
	topsRequire := i.env[envEFRequireTOPS]
	return len(legacyTopsVersion) > 0 && len(topsRequire) == 0
}

// GetRequirements returns the requirements from all ENFLAME_REQUIRE_ environment
// variables.
func (i TOPS) GetRequirements() ([]string, error) {
	if i.HasDisableRequire() {
		return nil, nil
	}

	// All variables with the "ENFLAME_REQUIRE_" prefix are passed to enflame-container-cli
	var requirements []string
	for name, value := range i.env {
		if strings.HasPrefix(name, envEFRequirePrefix) && !strings.HasPrefix(name, envEFRequireJetpack) {
			requirements = append(requirements, value)
		}
	}
	if i.IsLegacy() {
		v, err := i.legacyVersion()
		if err != nil {
			return nil, fmt.Errorf("failed to get version: %v", err)
		}
		topsRequire := fmt.Sprintf("tops>=%s", v)
		requirements = append(requirements, topsRequire)
	}
	return requirements, nil
}

// HasDisableRequire checks for the value of the ENFLAME_DISABLE_REQUIRE. If set
// to a valid (true) boolean value this can be used to disable the requirement checks
func (i TOPS) HasDisableRequire() bool {
	if disable, exists := i.env[envEFDisableRequire]; exists {
		// i.logger.Debugf("ENFLAME_DISABLE_REQUIRE=%v; skipping requirement checks", disable)
		d, _ := strconv.ParseBool(disable)
		return d
	}

	return false
}

// DevicesFromEnvvars returns the devices requested by the image through environment variables
func (i TOPS) DevicesFromEnvvars(envVars ...string) VisibleDevices {
	// We concantenate all the devices from the specified env.
	var isSet bool
	var devices []string
	requested := make(map[string]bool)
	for _, envVar := range envVars {
		if devs, ok := i.env[envVar]; ok {
			isSet = true
			for _, d := range strings.Split(devs, ",") {
				trimmed := strings.TrimSpace(d)
				if len(trimmed) == 0 {
					continue
				}
				devices = append(devices, trimmed)
				requested[trimmed] = true
			}
		}
	}

	// Environment variable unset with legacy image: default to "all".
	if !isSet && len(devices) == 0 && i.IsLegacy() {
		return NewVisibleDevices("all")
	}

	// Environment variable unset or empty or "void": return nil
	if len(devices) == 0 || requested["void"] {
		return NewVisibleDevices("void")
	}

	return NewVisibleDevices(devices...)
}

// GetDriverCapabilities returns the requested driver capabilities.
func (i TOPS) GetDriverCapabilities() DriverCapabilities {
	env := i.env[envEFDriverCapabilities]

	capabilities := make(DriverCapabilities)
	for _, c := range strings.Split(env, ",") {
		capabilities[DriverCapability(c)] = true
	}

	return capabilities
}

func (i TOPS) legacyVersion() (string, error) {
	topsVersion := i.env[envTOPSVersion]
	majorMinor, err := parseMajorMinorVersion(topsVersion)
	if err != nil {
		return "", fmt.Errorf("invalid TOPS version %v: %v", topsVersion, err)
	}

	return majorMinor, nil
}

func parseMajorMinorVersion(version string) (string, error) {
	vVersion := "v" + strings.TrimPrefix(version, "v")

	if !semver.IsValid(vVersion) {
		return "", fmt.Errorf("invalid version string")
	}

	majorMinor := strings.TrimPrefix(semver.MajorMinor(vVersion), "v")
	parts := strings.Split(majorMinor, ".")

	var err error
	_, err = strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid major version")
	}
	_, err = strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid minor version")
	}
	return majorMinor, nil
}

// OnlyFullyQualifiedCDIDevices returns true if all devices requested in the image are requested as CDI devices/
func (i TOPS) OnlyFullyQualifiedCDIDevices() bool {
	var hasCDIdevice bool
	for _, device := range i.DevicesFromEnvvars("ENFLAME_VISIBLE_DEVICES").List() {
		if !parser.IsQualifiedName(device) {
			return false
		}
		hasCDIdevice = true
	}

	for _, device := range i.DevicesFromMounts() {
		if !strings.HasPrefix(device, "cdi/") {
			return false
		}
		hasCDIdevice = true
	}
	return hasCDIdevice
}

const (
	deviceListAsVolumeMountsRoot = "/var/run/enflame-container-devices"
)

// DevicesFromMounts returns a list of device specified as mounts.
// TODO: This should be merged with getDevicesFromMounts used in the ENFLAME Container Runtime
func (i TOPS) DevicesFromMounts() []string {
	root := filepath.Clean(deviceListAsVolumeMountsRoot)
	seen := make(map[string]bool)
	var devices []string
	for _, m := range i.mounts {
		source := filepath.Clean(m.Source)
		// Only consider mounts who's host volume is /dev/null
		if source != "/dev/null" {
			continue
		}

		destination := filepath.Clean(m.Destination)
		if seen[destination] {
			continue
		}
		seen[destination] = true

		// Only consider container mount points that begin with 'root'
		if !strings.HasPrefix(destination, root) {
			continue
		}

		// Grab the full path beyond 'root' and add it to the list of devices
		device := strings.Trim(strings.TrimPrefix(destination, root), "/")
		if len(device) == 0 {
			continue
		}
		devices = append(devices, device)
	}
	return devices
}

// CDIDevicesFromMounts returns a list of CDI devices specified as mounts on the image.
func (i TOPS) CDIDevicesFromMounts() []string {
	var devices []string
	for _, mountDevice := range i.DevicesFromMounts() {
		if !strings.HasPrefix(mountDevice, "cdi/") {
			continue
		}
		parts := strings.SplitN(strings.TrimPrefix(mountDevice, "cdi/"), "/", 3)
		if len(parts) != 3 {
			continue
		}
		vendor := parts[0]
		class := parts[1]
		device := parts[2]
		devices = append(devices, fmt.Sprintf("%s/%s=%s", vendor, class, device))
	}
	return devices
}
