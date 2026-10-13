# Enflame Container Toolkit (enflame-container-toolkit)

> **Language**: English | [中文](README_zh.md)

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Build from Source](#build-from-source)
- [Quick Start](#quick-start)
  - [Build the Binaries](#build-the-binaries)
  - [Configure the Container Runtime](#configure-the-container-runtime)
  - [Run a GCU Container](#run-a-gcu-container)
- [Key Features](#key-features)
- [Version Compatibility](#version-compatibility)
- [Hardware Support](#hardware-support)
- [Trouble Shooting](#trouble-shooting)
- [License](#license)

## Overview

The **Enflame Container Toolkit** enables users to build and run
GCU-accelerated containers. It provides an OCI-compatible container runtime, a
runtime hook, a command-line interface, and CDI (Container Device Interface)
support that automatically configure containers to access Enflame GCU (General
Compute Unit) devices. It is the Enflame analog of the
`nvidia-container-toolkit`, and works together with the low-level
`libenflame-container` library.

The toolkit delivers:

- **OCI-compatible runtime** — `enflame-container-runtime` wraps the underlying OCI runtime and injects GCU devices into containers at start time.
- **Runtime hook & CDI support** — a runtime hook (`enflame-container-runtime-hook`) and a Container Device Interface path (`enflame-container-runtime.cdi`) automatically prepare containers for GCU access.
- **Simple configuration** — the `enflame-ctk` CLI configures the container runtime so that containers transparently obtain GCU devices.
- **Works with common runtimes** — integrates with containerd and Docker, and interoperates with the `libenflame-container` library and the `enflame-docker` usage layer.

## Project Structure

```
enflame-container-toolkit/
├── cmd/                        # Entry points for the built binaries
│   ├── enflame-container-runtime/        # OCI-compatible container runtime
│   ├── enflame-container-runtime-hook/   # Runtime hook (prestart configuration)
│   ├── enflame-container-runtime.cdi/    # CDI-based runtime path
│   ├── enflame-container-runtime.legacy/ # Legacy runtime path
│   └── enflame-ctk/                      # Command-line configuration tool
├── internal/                   # Private application packages
├── pkg/                        # Reusable library packages
├── deployments/                # Deployment manifests and assets
├── docker/                     # Docker build assets
├── packaging/                  # OS package build files
├── scripts/                    # Helper scripts
├── test/                       # Tests
├── tools/                      # Development tooling
├── go.mod                      # Go module definition
├── Makefile                    # Builds all binaries
└── versions.mk                 # Version definitions
```

## Installation

### Prerequisites

- An Enflame GCU device with the Enflame GCU driver installed.
- A container runtime — containerd or Docker.
- [Go](https://go.dev/) and a working toolchain for building from source.
- The low-level `libenflame-container` library available on the host.

### Build from Source

```bash
# Build the runtime, hook, CDI, legacy and CLI binaries
make

# Clean build artifacts
make clean
```

Building produces the runtime, hook and CLI binaries defined under `cmd/`:
`enflame-container-runtime`, `enflame-container-runtime-hook`,
`enflame-container-runtime.cdi`, `enflame-container-runtime.legacy` and
`enflame-ctk`.

## Quick Start

### Build the Binaries

```bash
make
```

### Configure the Container Runtime

Use the `enflame-ctk` CLI to register the Enflame runtime with your container
runtime (containerd or Docker). Once configured, the container runtime routes
GCU workloads through `enflame-container-runtime`, which injects the GCU
devices via the runtime hook / CDI path.

```bash
# Configure the container runtime with the Enflame runtime (illustrative)
enflame-ctk runtime configure

# Restart the container runtime so the change takes effect
# (e.g. containerd or docker, per your environment)
```

### Run a GCU Container

After the runtime is configured, launch a container using the Enflame runtime.
Containers started this way automatically get GCU devices injected and can
access the Enflame GCUs from inside the container.

```bash
# Launch a container using the Enflame runtime (illustrative)
docker run --runtime=enflame --rm -it <your-image> <command>
```

## Key Features

| Feature | Description |
| --- | --- |
| OCI runtime | `enflame-container-runtime` is an OCI-compatible runtime that injects GCU devices into containers |
| Runtime hook | `enflame-container-runtime-hook` prepares containers for GCU access at start time |
| CDI support | Container Device Interface path (`enflame-container-runtime.cdi`) for declarative device injection |
| enflame-ctk CLI | Configures the container runtime and containers so they transparently access GCUs |
| Runtime integration | Works with containerd and Docker |
| Library interop | Builds on the low-level `libenflame-container` library and the `enflame-docker` usage layer |

## Version Compatibility

The Enflame Container Toolkit configures access to GCU devices exposed by the
Enflame GCU driver, so its behavior tracks the installed driver. Always match
the toolkit release to the GCU driver deployed on your hosts, and ensure the
`libenflame-container` library is compatible with that driver.

## Hardware Support

Supports Enflame GCU devices that provide the Enflame GCU driver. The runtime,
hook and CLI work across GCU product generations; the devices made available
inside containers depend on the driver installed on the host.

## Trouble Shooting

- **Container cannot see GCU devices** — verify the Enflame GCU driver is loaded
  on the host and that the Enflame runtime was registered with your container
  runtime via `enflame-ctk`; restart containerd/Docker after configuring.
- **`enflame` runtime not found** — confirm the runtime binaries are installed
  and on the expected path, and that the container runtime configuration
  references `enflame-container-runtime`.
- **Library errors at container start** — ensure the `libenflame-container`
  library is installed and compatible with the running GCU driver.

## License

The Enflame Container Toolkit is licensed under the **Apache License 2.0**.
Copyright (c) Enflame Technology. All Rights Reserved. See [LICENSE](LICENSE).
