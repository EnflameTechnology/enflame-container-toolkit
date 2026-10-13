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

package operator

import "path/filepath"

const (
	defaultRuntimeName = "enflame"

	defaultRoot = "/usr/bin"
)

// Runtime defines a runtime to be configured.
// The path and whether the runtime is the default runtime can be specfied
type Runtime struct {
	name         string
	Path         string
	SetAsDefault bool
}

// Runtimes defines a set of runtimes to be configure for use in the GCU Operator
type Runtimes map[string]Runtime

type config struct {
	root              string
	enflameRuntimeName string
	setAsDefault      bool
}

// GetRuntimes returns the set of runtimes to be configured for use with the GCU Operator.
func GetRuntimes(opts ...Option) Runtimes {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}

	if c.root == "" {
		c.root = defaultRoot
	}
	if c.enflameRuntimeName == "" {
		c.enflameRuntimeName = defaultRuntimeName
	}

	runtimes := make(Runtimes)
	runtimes.add(c.enflameRuntime())

	modes := []string{"cdi", "legacy"}
	for _, mode := range modes {
		runtimes.add(c.modeRuntime(mode))
	}
	return runtimes
}

// DefaultRuntimeName returns the name of the default runtime.
func (r Runtimes) DefaultRuntimeName() string {
	for _, runtime := range r {
		if runtime.SetAsDefault {
			return runtime.name
		}
	}
	return ""
}

// Add a runtime to the set of runtimes.
func (r *Runtimes) add(runtime Runtime) {
	(*r)[runtime.name] = runtime
}

// enflameRuntime creates a runtime that corresponds to the enflame runtime.
// If name is equal to one of the predefined runtimes, `enflame` is used as the runtime name instead.
func (c config) enflameRuntime() Runtime {
	predefinedRuntimes := map[string]struct{}{
		"enflame-cdi":    {},
		"enflame-legacy": {},
	}
	name := c.enflameRuntimeName
	if _, isPredefinedRuntime := predefinedRuntimes[name]; isPredefinedRuntime {
		name = defaultRuntimeName
	}
	return c.newRuntime(name, "enflame-container-runtime")
}

// modeRuntime creates a runtime for the specified mode.
func (c config) modeRuntime(mode string) Runtime {
	return c.newRuntime("enflame-"+mode, "enflame-container-runtime."+mode)
}

// newRuntime creates a runtime based on the configuration
func (c config) newRuntime(name string, binary string) Runtime {
	return Runtime{
		name:         name,
		Path:         filepath.Join(c.root, binary),
		SetAsDefault: c.setAsDefault && name == c.enflameRuntimeName,
	}
}

// Option is a functional option for configuring set of runtimes.
type Option func(*config)

// WithRoot sets the root directory for the runtime binaries.
func WithRoot(root string) Option {
	return func(c *config) {
		c.root = root
	}
}

// WithEnflameRuntimeName sets the name of the enflame runtime.
func WithEnflameRuntimeName(name string) Option {
	return func(c *config) {
		c.enflameRuntimeName = name
	}
}

// WithSetAsDefault sets the default runtime to the enflame runtime.
func WithSetAsDefault(set bool) Option {
	return func(c *config) {
		c.setAsDefault = set
	}
}
