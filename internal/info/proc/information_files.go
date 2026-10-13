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

package proc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GCUInfoField represents the field name for information specified in a GCU's information file
type GCUInfoField string

// The following constants define the fields of interest from the GCU information file
const (
	GCUInfoModel       = GCUInfoField("Model")
	GCUInfoGCUUUID     = GCUInfoField("GCU UUID")
	GCUInfoBusLocation = GCUInfoField("Bus Location")
	GCUInfoDeviceMinor = GCUInfoField("Device Minor")
)

// GCUInfo stores the information for a GCU as determined from its associated information file
type GCUInfo map[GCUInfoField]string

// GetInformationFilePaths returns the list of information files associated with ENFLAME GCUs.
func GetInformationFilePaths(root string) ([]string, error) {
	return filepath.Glob(filepath.Join(root, "/proc/driver/enflame/gcus/*/information"))
}

// ParseGCUInformationFile parses the specified GCU information file and constructs a GCUInfo structure
func ParseGCUInformationFile(path string) (GCUInfo, error) {
	infoFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %v: %v", path, err)
	}
	defer infoFile.Close()

	return gcuInfoFrom(infoFile), nil
}

// gcuInfoFrom parses a GCUInfo struct from the specified reader
// An information file has the following structure:
// $ cat /proc/driver/enflame/gcus/0000\:06\:00.0/information
// Model:           Tesla V100-SXM2-16GB
// IRQ:             408
// GCU UUID:        GCU-edfee158-11c1-52b8-0517-92f30e7fac88
// Video BIOS:      88.00.41.00.01
// Bus Type:        PCIe
// DMA Size:        47 bits
// DMA Mask:        0x7fffffffffff
// Bus Location:    0000:06:00.0
// Device Minor:    0
// GCU Excluded:    No
func gcuInfoFrom(reader io.Reader) GCUInfo {
	info := make(GCUInfo)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		field := GCUInfoField(parts[0])
		value := strings.TrimSpace(parts[1])

		info[field] = value
	}

	return info
}
