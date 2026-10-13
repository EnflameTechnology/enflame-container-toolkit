# ENFLAME Container Toolkit CLI

The ENFLAME Container Toolkit CLI `enflame-ctk` provides a number of utilities that are useful for working with the ENFLAME Container Toolkit.

## Functionality

### Configure runtimes

The `runtime` command of the `enflame-ctk` CLI provides a set of utilities to related to the configuration
and management of supported container engines.

For example, running the following command:
```bash
enflame-ctk runtime configure --set-as-default
```
will ensure that the ENFLAME Container Runtime is added as the default runtime to the default container
engine.

## Configure the ENFLAME Container Toolkit

The `config` command of the `enflame-ctk` CLI allows a user to display and manipulate the ENFLAME Container Toolkit
configuration.

For example, running the following command:
```bash
enflame-ctk config default
```
will display the default config for the detected platform.

Whereas
```bash
enflame-ctk config
```
will display the effective ENFLAME Container Toolkit config using the configured config file, and running:

Individual config options can be set by specifying these are key-value pairs to the `--set` argument:

```bash
enflame-ctk config --set enflame-container-cli.no-cgroups=true
```

By default, all commands output to `STDOUT`, but specifying the `--output` flag writes the config to the specified file.

### Generate CDI specifications

The [Container Device Interface (CDI)](https://tags.cncf.io/container-device-interface) provides
a vendor-agnostic mechanism to make arbitrary devices accessible in containerized environments. To allow ENFLAME devices to be
used in these environments, the ENFLAME Container Toolkit CLI includes functionality to generate a CDI specification for the
available ENFLAME GCUs in a system.

In order to generate the CDI specification for the available devices, run the following command:\
```bash
enflame-ctk cdi generate
```

The default is to print the specification to STDOUT and a filename can be specified using the `--output` flag.

The specification will contain a device entries as follows (where applicable):
* An `enflame.com/gcu=gcu{INDEX}` device for each non-MIG-enabled full GCU in the system
* An `enflame.com/gcu=mig{GTOPS_INDEX}:{MIG_INDEX}` device for each MIG-device in the system
* A special device called `enflame.com/gcu=all` which represents all available devices.

For example, to generate the CDI specification in the default location where CDI-enabled tools such as `podman`, `containerd`, `cri-o`, or the ENFLAME Container Runtime can be configured to load it, the following command can be run:

```bash
sudo enflame-ctk cdi generate --output=/etc/cdi/enflame.yaml
```
(Note that `sudo` is used to ensure the correct permissions to write to the `/etc/cdi` folder)

With the specification generated, a GCU can be requested by specifying the fully-qualified CDI device name. With `podman` as an exmaple:
```bash
podman run --rm -ti --device=enflame.com/gcu=gcu0 ubuntu enflame-smi -L
```
