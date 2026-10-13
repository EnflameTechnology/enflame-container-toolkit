package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"

	"enflame/enflame-container-toolkit/internal/config"
	"enflame/enflame-container-toolkit/internal/info"
	"enflame/enflame-container-toolkit/internal/logger"
	"enflame/enflame-container-toolkit/internal/lookup"
)

var (
	debugflag   = flag.Bool("debug", false, "enable debug output")
	versionflag = flag.Bool("version", false, "enable version output")
	configflag  = flag.String("config", "", "configuration file")
)

func exit() {
	if err := recover(); err != nil {
		if _, ok := err.(runtime.Error); ok {
			log.Println(err)
		}
		if *debugflag {
			log.Printf("%s", debug.Stack())
		}
		os.Exit(1)
	}
	os.Exit(0)
}

func getCLIPath(config config.ContainerCLIConfig) string {
	if config.Path != "" {
		return config.Path
	}

	if err := os.Setenv("PATH", lookup.GetPath(config.Root)); err != nil {
		log.Panicln("couldn't set PATH variable:", err)
	}

	path, err := exec.LookPath("enflame-container-cli")
	if err != nil {
		log.Panicln("couldn't find binary enflame-container-cli in", os.Getenv("PATH"), ":", err)
	}
	return path
}

// getRootfsPath returns an absolute path. We don't need to resolve symlinks for now.
func getRootfsPath(config containerConfig) string {
	rootfs, err := filepath.Abs(config.Rootfs)
	if err != nil {
		log.Panicln(err)
	}
	return rootfs
}

func doPrestart() {
	var err error

	defer exit()
	log.SetFlags(0)

	hook, err := getHookConfig()
	if err != nil || hook == nil {
		log.Panicln("error getting hook config:", err)
	}
	cli := hook.ENFLAMEContainerCLIConfig

	container := getContainerConfig(*hook)
	enflame := container.Enflame
	if enflame == nil {
		// Not a GCU container, nothing to do.
		return
	}

	if !hook.ENFLAMEContainerRuntimeHookConfig.SkipModeDetection && info.ResolveAutoMode(&logInterceptor{}, hook.ENFLAMEContainerRuntimeConfig.Mode, container.Image) != "legacy" {
		log.Panicln("invoking the ENFLAME Container Runtime Hook directly (e.g. specifying the docker --gcus flag) is not supported. Please use the ENFLAME Container Runtime (e.g. specify the --runtime=enflame flag) instead.")
	}

	rootfs := getRootfsPath(container)

	args := []string{getCLIPath(cli)}
	if cli.Root != "" {
		args = append(args, fmt.Sprintf("--root=%s", cli.Root))
	}
	if cli.LoadKmods {
		args = append(args, "--load-kmods")
	}
	if cli.NoPivot {
		args = append(args, "--no-pivot")
	}
	if *debugflag {
		args = append(args, "--debug=/dev/stderr")
	} else if cli.Debug != "" {
		args = append(args, fmt.Sprintf("--debug=%s", cli.Debug))
	}
	if cli.Ldcache != "" {
		args = append(args, fmt.Sprintf("--ldcache=%s", cli.Ldcache))
	}
	if cli.User != "" {
		args = append(args, fmt.Sprintf("--user=%s", cli.User))
	}
	args = append(args, "configure")

	if ldconfigPath := cli.NormalizeLDConfigPath(); ldconfigPath != "" {
		args = append(args, fmt.Sprintf("--ldconfig=%s", ldconfigPath))
	}
	if cli.NoCgroups {
		args = append(args, "--no-cgroups")
	}
	if len(enflame.Devices) > 0 {
		args = append(args, fmt.Sprintf("--device=%s", enflame.Devices))
	}
// EF_CLOSE
/*
	if len(enflame.MigConfigDevices) > 0 {
		args = append(args, fmt.Sprintf("--mig-config=%s", enflame.MigConfigDevices))
	}
	if len(enflame.MigMonitorDevices) > 0 {
		args = append(args, fmt.Sprintf("--mig-monitor=%s", enflame.MigMonitorDevices))
	}
	if len(enflame.ImexChannels) > 0 {
		args = append(args, fmt.Sprintf("--imex-channel=%s", enflame.ImexChannels))
	}
*/
	for _, cap := range strings.Split(enflame.DriverCapabilities, ",") {
		if len(cap) == 0 {
			break
		}
		args = append(args, capabilityToCLI(cap))
	}

	for _, req := range enflame.Requirements {
		args = append(args, fmt.Sprintf("--require=%s", req))
	}

	args = append(args, fmt.Sprintf("--pid=%s", strconv.FormatUint(uint64(container.Pid), 10)))
	args = append(args, rootfs)

	env := append(os.Environ(), cli.Environment...)
	//nolint:gosec // TODO: Can we harden this so that there is less risk of command injection?
	err = syscall.Exec(args[0], args, env)
	log.Panicln("exec failed:", err)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nCommands:\n")
	fmt.Fprintf(os.Stderr, "  prestart\n        run the prestart hook\n")
	fmt.Fprintf(os.Stderr, "  poststart\n        no-op\n")
	fmt.Fprintf(os.Stderr, "  poststop\n        no-op\n")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *versionflag {
		fmt.Printf("%v version %v\n", "ENFLAME Container Runtime Hook", info.GetVersionString())
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	switch args[0] {
	case "prestart":
		doPrestart()
		os.Exit(0)
	case "poststart":
		fallthrough
	case "poststop":
		os.Exit(0)
	default:
		flag.Usage()
		os.Exit(2)
	}
}

// logInterceptor implements the logger.Interface to allow for logging from executable.
type logInterceptor struct {
	logger.NullLogger
}

func (l *logInterceptor) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}
