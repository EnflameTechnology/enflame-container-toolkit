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
	"fmt"

	// "github.com/ENFLAME/go-eflib/pkg/efml"
)

// UUIDer is an interface for getting UUIDs.
type UUIDer interface {
	GetUUID() (string, error)
}

// DeviceNamers represents a list of device namers
type DeviceNamers []DeviceNamer

// DeviceNamer is an interface for getting device names
type DeviceNamer interface {
	GetDeviceName(int, UUIDer) (string, error)
	GetMigDeviceName(int, UUIDer, int, UUIDer) (string, error)
}

// Supported device naming strategies
const (
	// DeviceNameStrategyIndex generates devices names such as 0 or 1:0
	DeviceNameStrategyIndex = "index"
	// DeviceNameStrategyTypeIndex generates devices names such as gcu0 or mig1:0
	DeviceNameStrategyTypeIndex = "type-index"
	// DeviceNameStrategyUUID uses the device UUID as the name
	DeviceNameStrategyUUID = "uuid"
)

type deviceNameIndex struct {
	gcuPrefix string
	migPrefix string
}
type deviceNameUUID struct{}

// NewDeviceNamer creates a Device Namer based on the supplied strategy.
// This namer can be used to construct the names for MIG and GCU devices when generating the CDI spec.
func NewDeviceNamer(strategy string) (DeviceNamer, error) {
	switch strategy {
	case DeviceNameStrategyIndex:
		return deviceNameIndex{}, nil
	case DeviceNameStrategyTypeIndex:
		return deviceNameIndex{gcuPrefix: "gcu", migPrefix: "mig"}, nil
	case DeviceNameStrategyUUID:
		return deviceNameUUID{}, nil
	}

	return nil, fmt.Errorf("invalid device name strategy: %v", strategy)
}

// GetDeviceName returns the name for the specified device based on the naming strategy
func (s deviceNameIndex) GetDeviceName(i int, _ UUIDer) (string, error) {
	return fmt.Sprintf("%s%d", s.gcuPrefix, i), nil
}

// GetMigDeviceName returns the name for the specified device based on the naming strategy
func (s deviceNameIndex) GetMigDeviceName(i int, _ UUIDer, j int, _ UUIDer) (string, error) {
	return fmt.Sprintf("%s%d:%d", s.migPrefix, i, j), nil
}

// GetDeviceName returns the name for the specified device based on the naming strategy
func (s deviceNameUUID) GetDeviceName(i int, d UUIDer) (string, error) {
	uuid, err := d.GetUUID()
	if err != nil {
		return "", fmt.Errorf("failed to get device UUID: %v", err)
	}
	return uuid, nil
}

// GetMigDeviceName returns the name for the specified device based on the naming strategy
func (s deviceNameUUID) GetMigDeviceName(i int, _ UUIDer, j int, mig UUIDer) (string, error) {
	uuid, err := mig.GetUUID()
	if err != nil {
		return "", fmt.Errorf("failed to get device UUID: %v", err)
	}
	return uuid, nil
}

//go:generate moq -stub -out namer_efml_mock.go . efmlUUIDer
type efmlUUIDer interface {
// EF_CLOSE
	// GetUUID() (string, efml.Return)
	GetUUID() (string, error)
}

type convert struct {
	efmlUUIDer
}

type uuidIgnored struct{}

func (m uuidIgnored) GetUUID() (string, error) {
	return "", nil
}

type uuidUnsupported struct{}

func (m convert) GetUUID() (string, error) {
	if m.efmlUUIDer == nil {
		return uuidUnsupported{}.GetUUID()
	}

	// EF_CLOSE
	// uuid, ret := m.efmlUUIDer.GetUUID()
	// if ret != efml.SUCCESS {
	// 	return "", ret
	// }
	uuid := "N/A"

	return uuid, nil
}

var errUUIDUnsupported = errors.New("GetUUID is not supported")

func (m uuidUnsupported) GetUUID() (string, error) {
	return "", errUUIDUnsupported
}

func (l DeviceNamers) GetDeviceNames(i int, d UUIDer) ([]string, error) {
	var names []string
	for _, namer := range l {
		name, err := namer.GetDeviceName(i, d)
		if err != nil {
			return nil, err
		}
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, errors.New("no names defined")
	}
	return names, nil
}

func (l DeviceNamers) GetMigDeviceNames(i int, d UUIDer, j int, mig UUIDer) ([]string, error) {
	var names []string
	for _, namer := range l {
		name, err := namer.GetMigDeviceName(i, d, j, mig)
		if err != nil {
			return nil, err
		}
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, errors.New("no names defined")
	}
	return names, nil
}
