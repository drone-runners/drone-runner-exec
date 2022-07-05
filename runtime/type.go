// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"github.com/drone-runners/drone-runner-exec/engine"
	"github.com/drone-runners/drone-runner-exec/engine/resource"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/manifest"
	"github.com/drone/runner-go/secret"
)

type (
	// CompilerArgs provides compiler arguments.
	CompilerArgs struct {
		// Manifest provides the parsed manifest.
		Manifest *manifest.Manifest

		// Pipeline provides the parsed pipeline. This pipeline is
		// the compiler source and is converted to the intermediate
		// representation by the Compile method.
		Pipeline *resource.Pipeline

		// Build provides the compiler with stage information that
		// is converted to environment variable format and passed to
		// each pipeline step. It is also used to clone the commit.
		Build *drone.Build

		// Stage provides the compiler with stage information that
		// is converted to environment variable format and passed to
		// each pipeline step.
		Stage *drone.Stage

		// Repo provides the compiler with repo information. This
		// repo information is converted to environment variable
		// format and passed to each pipeline step. It is also used
		// to clone the repository.
		Repo *drone.Repo

		// System provides the compiler with system information that
		// is converted to environment variable format and passed to
		// each pipeline step.
		System *drone.System

		// Netrc provides netrc parameters that can be used by the
		// default clone step to authenticate to the remote
		// repository.
		Netrc *drone.Netrc

		// Secret returns a named secret value that can be injected
		// into the pipeline step.
		Secret secret.Provider
	}

	// Compiler compiles the Yaml configuration file to an
	// intermediate representation optimized for simple execution.
	Compiler interface {
		Compile(context.Context, CompilerArgs) *engine.Spec
	}
)
