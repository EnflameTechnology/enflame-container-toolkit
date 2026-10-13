# ENFLAME Container Toolkit Changelog

## v1.15.0

* Remove `enflame-container-runtime` and `enflame-docker2` packages.
* Use `XDG_DATA_DIRS` environment variable when locating config files such as graphics config files.
* Add support for v0.7.0 Container Device Interface (CDI) specification.
* Add `--config-search-path` option to `enflame-ctk cdi generate` command. These paths are used when locating driver files such as graphics config files.
* Use D3DKMTEnumAdapters3 to enumerate adpaters on WSL2 if available.
* Add support for v1.2.0 OCI Runtime specification.
* Explicitly set `ENFLAME_VISIBLE_DEVICES=void` in generated CDI specifications. This prevents the ENFLAME Container Runtime from making additional modifications.

* [libenflame-container] Use D3DKMTEnumAdapters3 to enumerate adpaters on WSL2 if available.

* [toolkit-container] Bump TOPS base image version to 12.4.1

## v1.15.0-rc.4
* Add a `--spec-dir` option to the `enflame-ctk cdi generate` command. This allows specs outside of `/etc/cdi` and `/var/run/cdi` to be processed.
* Add support for extracting device major number from `/proc/devices` if `enflame` is used as a device name over `enflame-frontend`.
* Allow multiple device naming strategies for `enflame-ctk cdi generate` command. This allows a single
  CDI spec to be generated that includes GCUs by index and UUID.
* Set the default `--device-name-strategy` for the `enflame-ctk cdi generate` command to `[index, uuid]`.
* Remove `libenflame-container0` jetpack dependency included for legacy Tegra-based systems.
* Add `ENFLAME_VISIBLE_DEVICES=void` to generated CDI specifications.

* [toolkit-container] Remove centos7 image. The ubi8 image can be used on all RPM-based platforms.
* [toolkit-container] Bump TOPS base image version to 12.3.2

## v1.15.0-rc.3
* Fix bug in `enflame-ctk hook update-ldcache` where default `--ldconfig-path` value was not applied.

## v1.15.0-rc.2
* Extend the `runtime.enflame.com/gcu` CDI kind to support full-GCUs and MIG devices specified by index or UUID.
* Fix bug when specifying `--dev-root` for Tegra-based systems.
* Log explicitly requested runtime mode.
* Remove package dependency on libseccomp.
* Added detection of libefdxgdmal.so.1 on WSL2
* Use devRoot to resolve MIG device nodes.
* Fix bug in determining default enflame-container-runtime.user config value on SUSE-based systems.
* Add `crun` to the list of configured low-level runtimes.
* Added support for `--ldconfig-path` to `enflame-ctk cdi generate` command.
* Fix `enflame-ctk runtime configure --cdi.enabled` for Docker.
* Add discovery of the GDRCopy device (`gdrdrv`) if the `ENFLAME_GDRCOPY` environment variable of the container is set to `enabled`

* [toolkit-container] Bump TOPS base image version to 12.3.1.

## v1.15.0-rc.1
* Skip update of ldcache in containers without ldconfig. The .so.SONAME symlinks are still created.
* Normalize ldconfig path on use. This automatically adjust the ldconfig setting applied to ldconfig.real on systems where this exists.
* Include `enflame/nvoptix.bin` in list of graphics mounts.
* Include `vulkan/icd.d/enflame_layers.json` in list of graphics mounts.
* Add support for `--library-search-paths` to `enflame-ctk cdi generate` command.
* Add support for injecting /dev/enflame-efswitch* devices if the ENFLAME_EFSWITCH=enabled envvar is specified.
* Added support for `enflame-ctk runtime configure --enable-cdi` for the `docker` runtime. Note that this requires Docker >= 25.
* Fixed bug in `enflame-ctk config` command when using `--set`. The types of applied config options are now applied correctly.
* Add `--relative-to` option to `enflame-ctk transform root` command. This controls whether the root transformation is applied to host or container paths.
* Added automatic CDI spec generation when the `runtime.enflame.com/gcu=all` device is requested by a container.

* [libenflame-container] Fix device permission check when using cgroupv2 (fixes #227)

## v1.14.3
* [toolkit-container] Bump TOPS base image version to 12.2.2.

## v1.14.2
* Fix bug on Tegra-based systems where symlinks were not created in containers.
* Add --csv.ignore-pattern command line option to enflame-ctk cdi generate command.

## v1.14.1
* Fixed bug where contents of `/etc/enflame-container-runtime/config.toml` is ignored by the ENFLAME Container Runtime Hook.

* [libenflame-container] Use libelf.so on RPM-based systems due to removed mageia repositories hosting pmake and bmake.

## v1.14.0
* Promote v1.14.0-rc.3 to v1.14.0

## v1.14.0-rc.3
* Added support for generating OCI hook JSON file to `enflame-ctk runtime configure` command.
* Remove installation of OCI hook JSON from RPM package.
* Refactored config for `enflame-container-runtime-hook`.
* Added a `enflame-ctk config` command which supports setting config options using a `--set` flag.
* Added `--library-search-path` option to `enflame-ctk cdi generate` command in `csv` mode. This allows folders where
  libraries are located to be specified explicitly.
* Updated go-eflib to support devices which are not present in the PCI device database. This allows the creation of dev/char symlinks on systems with such devices installed.
* Added `UsesEFGCUModule` info function for more robust platform detection. This is required on Tegra-based systems where libenflame-ml.so is also supported.

* [toolkit-container] Set `ENFLAME_VISIBLE_DEVICES=void` to prevent injection of ENFLAME devices and drivers into the ENFLAME Container Toolkit container.

## v1.14.0-rc.2
* Fix bug causing incorrect enflame-smi symlink to be created on WSL2 systems with multiple driver roots.
* Remove dependency on coreutils when installing package on RPM-based systems.
* Create ouput folders if required when running `enflame-ctk runtime configure`
* Generate default config as post-install step.
* Added support for detecting GSP firmware at custom paths when generating CDI specifications.
* Added logic to skip the extraction of image requirements if `ENFLAME_DISABLE_REQUIRES` is set to `true`.

* [libenflame-container] Include Shared Compiler Library (libenflame-gcucomp.so) in the list of compute libaries.

* [toolkit-container] Ensure that common envvars have higher priority when configuring the container engines.
* [toolkit-container] Bump TOPS base image version to 12.2.0.
* [toolkit-container] Remove installation of enflame-experimental runtime. This is superceded by the ENFLAME Container Runtime in CDI mode.

## v1.14.0-rc.1

* Add support for updating containerd configs to the `enflame-ctk runtime configure` command.
* Create file in `etc/ld.so.conf.d` with permissions `644` to support non-root containers.
* Generate CDI specification files with `644` permissions to allow rootless applications (e.g. podman)
* Add `enflame-ctk cdi list` command to show the known CDI devices.
* Add support for generating merged devices (e.g. `all` device) to the efcdi API.
* Use *.* pattern to locate libtops.so when generating a CDI specification to support platforms where a patch version is not specified.
* Update go-eflib to skip devices that are not MIG capable when generating CDI specifications.
* Add `enflame-container-runtime-hook.path` config option to specify ENFLAME Container Runtime Hook path explicitly.
* Fix bug in creation of `/dev/char` symlinks by failing operation if kernel modules are not loaded.
* Add option to load kernel modules when creating device nodes
* Add option to create device nodes when creating `/dev/char` symlinks

* [libenflame-container] Support OpenSSL 3 with the Encrypt/Decrypt library

* [toolkit-container] Allow same envars for all runtime configs

## v1.13.1

* Update `update-ldcache` hook to only update ldcache if it exists.
* Update `update-ldcache` hook to create `/etc/ld.so.conf.d` folder if it doesn't exist.
* Fix failure when libtops cannot be located during XOrg library discovery.
* Fix CDI spec generation on systems that use `/etc/alternatives` (e.g. Debian)

## v1.13.0

* Promote 1.13.0-rc.3 to 1.13.0

## v1.13.0-rc.3

* Only initialize EFML for modes that require it when runing `enflame-ctk cdi generate`.
* Prefer /run over /var/run when locating enflame-persistenced and enflame-fabricmanager sockets.
* Fix the generation of CDI specifications for management containers when the driver libraries are not in the LDCache.
* Add transformers to deduplicate and simplify CDI specifications.
* Generate a simplified CDI specification by default. This means that entities in the common edits in a spec are not included in device definitions.
* Also return an error from the efcdi.New constructor instead of panicing.
* Detect XOrg libraries for injection and CDI spec generation.
* Add `enflame-ctk system create-device-nodes` command to create control devices.
* Add `enflame-ctk cdi transform` command to apply transforms to CDI specifications.
* Add `--vendor` and `--class` options to `enflame-ctk cdi generate`

* [libenflame-container] Fix segmentation fault when RPC initialization fails.
* [libenflame-container] Build centos variants of the ENFLAME Container Library with static libtirpc v1.3.2.
* [libenflame-container] Remove make targets for fedora35 as the centos8 packages are compatible.

* [toolkit-container] Add `enflame-container-runtime.modes.cdi.annotation-prefixes` config option that allows the CDI annotation prefixes that are read to be overridden.
* [toolkit-container] Create device nodes when generating CDI specification for management containers.
* [toolkit-container] Add `enflame-container-runtime.runtimes` config option to set the low-level runtime for the ENFLAME Container Runtime

## v1.13.0-rc.2

* Don't fail chmod hook if paths are not injected
* Only create `by-path` symlinks if CDI devices are actually requested.
* Fix possible blank `enflame-ctk` path in generated CDI specifications
* Fix error in postun scriplet on RPM-based systems
* Only check `ENFLAME_VISIBLE_DEVICES` for environment variables if no annotations are specified.
* Add `cdi.default-kind` config option for constructing fully-qualified CDI device names in CDI mode
* Add support for `accept-enflame-visible-devices-envvar-unprivileged` config setting in CDI mode
* Add `enflame-container-runtime-hook.skip-mode-detection` config option to bypass mode detection. This allows `legacy` and `cdi` mode, for example, to be used at the same time.
* Add support for generating CDI specifications for GDS and MOFED devices
* Ensure CDI specification is validated on save when generating a spec
* Rename `--discovery-mode` argument to `--mode` for `enflame-ctk cdi generate`
* [libenflame-container] Fix segfault on WSL2 systems
* [toolkit-container] Add `--cdi-enabled` flag to toolkit config
* [toolkit-container] Install `enflame-ctk` from toolkit container
* [toolkit-container] Use installed `enflame-ctk` path in ENFLAME Container Toolkit config
* [toolkit-container] Bump TOPS base images to 12.1.0
* [toolkit-container] Set `enflame-ctk` path in the
* [toolkit-container] Add `cdi.k8s.io/*` to set of allowed annotations in containerd config
* [toolkit-container] Generate CDI specification for use in management containers
* [toolkit-container] Install experimental runtime as `enflame-container-runtime.experimental` instead of `enflame-container-runtime-experimental`
* [toolkit-container] Install and configure mode-specific runtimes for `cdi` and `legacy` modes

## v1.13.0-rc.1

* Include MIG-enabled devices as GCUs when generating CDI specification
* Fix missing EFML symbols when running `enflame-ctk` on some platforms [#49]
* Add CDI spec generation for WSL2-based systems to `enflame-ctk cdi generate` command
* Add `auto` mode to `enflame-ctk cdi generate` command to automatically detect a WSL2-based system over a standard EFML-based system.
* Add mode-specific (`.cdi` and `.legacy`) ENFLAME Container Runtime binaries for use in the GCU Operator
* Discover all `gsb*.bin` GSP firmware files when generating CDI specification.
* Align `.deb` and `.rpm` release candidate package versions
* Remove `fedora35` packaging targets
* [libenflame-container] Include all `gsp*.bin` firmware files if present
* [libenflame-container] Align `.deb` and `.rpm` release candidate package versions
* [libenflame-container] Remove `fedora35` packaging targets
* [toolkit-container] Install `enflame-container-toolkit-operator-extensions` package for mode-specific executables.
* [toolkit-container] Allow `enflame-container-runtime.mode` to be set when configuring the ENFLAME Container Toolkit

## v1.12.0

* Promote `v1.12.0-rc.5` to `v1.12.0`
* Rename `enflame cdi generate` `--root` flag to `--driver-root` to better indicate intent
* [libenflame-container] Add eftopsbins.bin to DriverStore components under WSL2
* [toolkit-container] Bump TOPS base images to 12.0.1

## v1.12.0-rc.5

* Fix bug here the `enflame-ctk` path was not properly resolved. This causes failures to run containers when the runtime is configured in `csv` mode or if the `ENFLAME_DRIVER_CAPABILITIES` includes `graphics` or `display` (e.g. `all`).

## v1.12.0-rc.4

* Generate a minimum CDI spec version for improved compatibility.
* Add `--device-name-strategy` options to the `enflame-ctk cdi generate` command that can be used to control how device names are constructed.
* Set default for CDI device name generation to `index` to generate device names such as `enflame.com/gcu=0` or `enflame.com/gcu=1:0` by default.

## v1.12.0-rc.3

* Don't fail if by-path symlinks for DRM devices do not exist
* Replace the --json flag with a --format [json|yaml] flag for the enflame-ctk cdi generate command
* Ensure that the CDI output folder is created if required
* When generating a CDI specification use a blank host path for devices to ensure compatibility with the v0.4.0 CDI specification
* Add injection of Wayland JSON files
* Add GSP firmware paths to generated CDI specification
* Add --root flag to enflame-ctk cdi generate command

## v1.12.0-rc.2

* Inject Direct Rendering Manager (DRM) devices into a container using the ENFLAME Container Runtime
* Improve logging of errors from the ENFLAME Container Runtime
* Improve CDI specification generation to support rootless podman
* Use `enflame-ctk cdi generate` to generate CDI specifications instead of `enflame-ctk info generate-cdi`
* [libenflame-container] Skip creation of existing files when these are already mounted

## v1.12.0-rc.1

* Add support for multiple Docker Swarm resources
* Improve injection of Vulkan configurations and libraries
* Add `enflame-ctk info generate-cdi` command to generated CDI specification for available devices
* [libenflame-container] Include NVVM compiler library in compute libs

## v1.11.0

* Promote v1.11.0-rc.3 to v1.11.0

## v1.11.0-rc.3

* Build fedora35 packages
* Introduce an `enflame-container-toolkit-base` package for better dependency management
* Fix removal of `enflame-container-runtime-hook` on RPM-based systems
* Inject platform files into container on Tegra-based systems
* [toolkit container] Update TOPS base images to 11.7.1
* [libenflame-container] Preload libgcc_s.so.1 on arm64 systems

## v1.11.0-rc.2

* Allow `accept-enflame-visible-devices-*` config options to be set by toolkit container
* [libenflame-container] Fix bug where LDCache was not updated when the `--no-pivot-root` option was specified

## v1.11.0-rc.1

* Add discovery of GCUDirect Storage (`enflame-fs*`) devices if the `ENFLAME_GDS` environment variable of the container is set to `enabled`
* Add discovery of MOFED Infiniband devices if the `ENFLAME_MOFED` environment variable of the container is set to `enabled`
* Fix bug in CSV mode where libraries listed as `sym` entries in mount specification are not added to the LDCache.
* Rename `enflame-container-toolkit` executable to `enflame-container-runtime-hook` and create `enflame-container-toolkit` as a symlink to `enflame-container-runtime-hook` instead.
* Add `enflame-ctk runtime configure` command to configure the Docker config file (e.g. `/etc/docker/daemon.json`) for use with the ENFLAME Container Runtime.

## v1.10.0

* Promote v1.10.0-rc.3 to v1.10.0

## v1.10.0-rc.3

* Use default config instead of raising an error if config file cannot be found
* Ignore ENFLAME_REQUIRE_JETPACK* environment variables for requirement checks
* Fix bug in detection of Tegra systems where `/sys/devices/soc0/family` is ignored
* Fix bug where links to devices were detected as devices
* [libenflame-container] Fix bug introduced when adding libtopsdebugger.so to list of libraries

## v1.10.0-rc.2

* Add support for ENFLAME_REQUIRE_* checks for tops version and arch to csv mode
* Switch to debug logging to reduce log verbosity
* Support logging to logs requested in command line
* Fix bug when launching containers with relative root path (e.g. using containerd)
* Allow low-level runtime path to be set explicitly as enflame-container-runtime.runtimes option
* Fix failure to locate low-level runtime if PATH envvar is unset
* Replace experimental option for ENFLAME Container Runtime with enflame-container-runtime.mode = csv option
* Use csv as default mode on Tegra systems without EFML
* Add --version flag to all CLIs
* [libenflame-container] Bump libtirpc to 1.3.2
* [libenflame-container] Fix bug when running host ldconfig using glibc compiled with a non-standard prefix
* [libenflame-container] Add libtopsdebugger.so to list of compute libraries

## v1.10.0-rc.1

* Include enflame-ctk CLI in installed binaries
* Add experimental option to ENFLAME Container Runtime

## v1.9.0

* [libenflame-container] Add additional check for Tegra in /sys/.../family file in CLI
* [libenflame-container] Update jetpack-specific CLI option to only load Base CSV files by default
* [libenflame-container] Fix bug (from 1.8.0) when mounting GSP firmware into containers without /lib to /usr/lib symlinks
* [libenflame-container] Update efml.h to TOPS 11.6.1 efml_DEV 11.6.55
* [libenflame-container] Update switch statement to include new brands from latest efml.h
* [libenflame-container] Process all --require flags on Jetson platforms
* [libenflame-container] Fix long-standing issue with running ldconfig on Debian systems

## v1.8.1

* [libenflame-container] Fix bug in determining cgroup root when running in nested containers
* [libenflame-container] Fix permission issue when determining cgroup version

## v1.8.0

* Promote 1.8.0-rc.2-1 to 1.8.0

## v1.8.0-rc.2

* Remove support for building amazonlinux1 packages

## v1.8.0-rc.1

* [libenflame-container] Add support for cgroupv2
* Release toolkit-container images from enflame-container-toolkit repository

## v1.7.0

* Promote 1.7.0-rc.1-1 to 1.7.0
 * Bump Golang version to 1.16.4

## v1.7.0-rc.1

* Specify containerd runtime type as string in config tools to remove dependency on containerd package
* Add supported-driver-capabilities config option to allow for a subset of all driver capabilities to be specified

## v1.6.0

* Promote 1.6.0-rc.3-1 to 1.6.0
 * Fix unnecessary logging to stderr instead of configured enflame-container-runtime log file

## v1.6.0-rc.3

* Add supported-driver-capabilities config option to the enflame-container-toolkit
* Move OCI and command line checks for runtime to internal oci package

## v1.6.0-rc.2

* Use relative path to OCI specification file (config.json) if bundle path is not specified as an argument to the enflame-container-runtime

## v1.6.0-rc.1

* Add AARCH64 package for Amazon Linux 2
* Include enflame-container-runtime into enflame-container-toolkit package

## v1.5.1

* Fix bug where Docker Swarm device selection is ignored if ENFLAME_VISIBLE_DEVICES is also set
* Improve unit testing by using require package and adding coverage reports
* Remove unneeded go dependencies by running go mod tidy
* Move contents of pkg directory to cmd for CLI tools
* Ensure make binary target explicitly sets GOOS

## v1.5.0

* Add dependence on libenflame-container-tools >= 1.4.0
* Add golang check targets to Makefile
* Add Jenkinsfile definition for build targets
* Move docker.mk to docker folder

## v1.4.2

* Add dependence on libenflame-container-tools >= 1.3.3

## v1.4.1

* Ignore ENFLAME_VISIBLE_DEVICES for containers with insufficent privileges
* Add dependence on libenflame-container-tools >= 1.3.2

## v1.4.0

* Add 'compute' capability to list of defaults
* Add dependence on libenflame-container-tools >= 1.3.1

## v1.3.0

* Promote 1.3.0-rc.2-1 to 1.3.0
* Add dependence on libenflame-container-tools >= 1.3.0

## v1.3.0-rc.2

* 2c180947 Add more tests for new semantics with device list from volume mounts
* 7c003857 Refactor accepting device lists from volume mounts as a boolean

## v1.3.0-rc.1

* b50d86c1 Update build system to accept a TAG variable for things like rc.x
* fe65573b Add common CI tests for things like golint, gofmt, unit tests, etc.
* da6fbb34 Revert "Add ability to merge envars of the form ENFLAME_VISIBLE_DEVICES_*"
* a7fb3330 Flip build-all targets to run automatically on merge requests
* 8b248b66 Rename github.com/ENFLAME/container-toolkit to enflame-container-toolkit
* da36874e Add new config options to pull device list from mounted files instead of ENVVAR

## v1.2.1

* 4e6e0ed4 Add 'ngx' to list of*all* driver capabilities
* 2f4af743 List config.toml as a config file in the RPM SPEC

## v1.2.0

*  8e0aab46 Fix repo listed in changelog for debian distributions
*  320bb6e4 Update dependence on libenflame-container to 1.2.0
*  6cfc8097 Update package license to match source license
*  e7dc3cbb Fix debian copyright file
*  d3aee3e0 Add the 'ngx' driver capability

## v1.1.2

* c32237f3 Add support for parsing Linux Capabilities for older OCI specs

## v1.1.1

* d202aded Update dependence to libenflame-container 1.1.1

## v1.1.0

* 4e4de762 Update build system to support multi-arch builds
* fcc1d116 Add support for MIG (Multi-Instance GCUs)
* d4ff0416 Add ability to merge envars of the form ENFLAME_VISIBLE_DEVICES_*
* 60f165ad Add no-pivot option to toolkit

## v1.0.5

* Initial release. Replaces older package enflame-container-runtime-hook. (Closes: #XXXXXX)
