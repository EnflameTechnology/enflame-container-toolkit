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

package tops
//EF_CLOSE
/*
import (
	"fmt"

	// "github.com/ENFLAME/go-efml/pkg/dl"
)

/*
#cgo linux LDFLAGS: -Wl,--export-dynamic -Wl,--unresolved-symbols=ignore-in-object-files
#cgo darwin LDFLAGS: -Wl,-undefined,dynamic_lookup

#ifdef _WIN32
#define TOPSAPI __stdcall
#else
#define TOPSAPI
#endif

typedef int TOPSdevice;

typedef enum TOPSdevice_attribute_enum {
    TOPS_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MAJOR = 75,
    TOPS_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MINOR = 76
} TOPSdevice_attribute;

typedef enum topsError_enum {
	TOPS_SUCCESS = 0
} CUresult;

CUresult TOPSAPI topsInit(unsigned int Flags);
CUresult TOPSAPI topsDriverGetVersion(int *driverVersion);
CUresult TOPSAPI topsDeviceGet(TOPSdevice *device, int ordinal);
CUresult TOPSAPI topsDeviceGetAttribute(int *pi, TOPSdevice_attribute attrib, TOPSdevice dev);
*
import "C"

const (
	libraryName      = "libtops.so.1"
	libraryLoadFlags = dl.RTLD_LAZY | dl.RTLD_GLOBAL
)

// tops stores a reference the tops dynamic library
var lib *dl.DynamicLibrary

// Version returns the TOPS version of the driver as a string or an error if this
// cannot be determined.
func Version() (string, error) {
	lib, err := load()
	if err != nil {
		return "", err
	}
	defer lib.Close()

	if err := lib.Lookup("topsDriverGetVersion"); err != nil {
		return "", fmt.Errorf("failed to lookup symbol: %v", err)
	}

	var version C.int
	if result := C.topsDriverGetVersion(&version); result != C.TOPS_SUCCESS {
		return "", fmt.Errorf("failed to get TOPS version: result=%v", result)
	}

	major := version / 1000
	minor := version % 100 / 10

	return fmt.Sprintf("%d.%d", major, minor), nil
}

// ComputeCapability returns the TOPS compute capability of a device with the specified index as a string
// or an error if this cannot be determined.
func ComputeCapability(index int) (string, error) {
	lib, err := load()
	if err != nil {
		return "", err
	}
	defer lib.Close()

	if err := lib.Lookup("topsInit"); err != nil {
		return "", fmt.Errorf("failed to lookup symbol: %v", err)
	}
	if err := lib.Lookup("topsDeviceGet"); err != nil {
		return "", fmt.Errorf("failed to lookup symbol: %v", err)
	}
	if err := lib.Lookup("topsDeviceGetAttribute"); err != nil {
		return "", fmt.Errorf("failed to lookup symbol: %v", err)
	}

	if result := C.topsInit(C.uint(0)); result != C.TOPS_SUCCESS {
		return "", fmt.Errorf("failed to initialize TOPS: result=%v", result)
	}

	var device C.TOPSdevice
	// NOTE: We only query the first device
	if result := C.topsDeviceGet(&device, C.int(index)); result != C.TOPS_SUCCESS {
		return "", fmt.Errorf("failed to get TOPS device %v: result=%v", 0, result)
	}

	var major C.int
	if result := C.topsDeviceGetAttribute(&major, C.TOPS_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MAJOR, device); result != C.TOPS_SUCCESS {
		return "", fmt.Errorf("failed to get TOPS compute capability major for device %v : result=%v", 0, result)
	}

	var minor C.int
	if result := C.topsDeviceGetAttribute(&minor, C.TOPS_DEVICE_ATTRIBUTE_COMPUTE_CAPABILITY_MINOR, device); result != C.TOPS_SUCCESS {
		return "", fmt.Errorf("failed to get TOPS compute capability minor for device %v: result=%v", 0, result)
	}

	return fmt.Sprintf("%d.%d", major, minor), nil
}

func load() (*dl.DynamicLibrary, error) {
	lib := dl.New(libraryName, libraryLoadFlags)
	if lib == nil {
		return nil, fmt.Errorf("error instantiating DynamicLibrary for TOPS")
	}
	err := lib.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening DynamicLibrary for TOPS: %v", err)
	}

	return lib, nil
}
*/