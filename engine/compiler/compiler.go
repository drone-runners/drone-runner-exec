// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package compiler

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/drone-runners/drone-runner-exec/engine"
	"github.com/drone-runners/drone-runner-exec/engine/resource"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/clone"
	"github.com/drone/runner-go/environ"
	"github.com/drone/runner-go/manifest"
	"github.com/drone/runner-go/secret"
	"github.com/drone/runner-go/shell"

	"github.com/dchest/uniuri"
)

// random generator function
var random = uniuri.New

// Compiler compiles the Yaml configuration file to an
// intermediate representation optimized for simple execution.
type Compiler struct {
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

	// Environ provides a set of environment varaibles that
	// should be added to each pipeline step by default.
	Environ map[string]string

	// Netrc provides netrc parameters that can be used by the
	// default clone step to authenticate to the remote
	// repository.
	Netrc *drone.Netrc

	// Secret returns a named secret value that can be injected
	// into the pipeline step.
	Secret secret.Provider
}

// Compile compiles the configuration file.
func (c *Compiler) Compile(ctx context.Context) *engine.Spec {
	spec := new(engine.Spec)
	spec.Root = fmt.Sprintf("/tmp/drone-%s", random())

	spec.Platform.OS = c.Pipeline.Platform.OS
	spec.Platform.Arch = c.Pipeline.Platform.Arch
	spec.Platform.Variant = c.Pipeline.Platform.Variant
	spec.Platform.Version = c.Pipeline.Platform.Version

	// creates a home directory in the root.
	homedir := filepath.Join(spec.Root, "home", "drone")
	spec.Files = append(spec.Files, &engine.File{
		Path:  homedir,
		Mode:  0700,
		IsDir: true,
	})

	// creates a source directory in the root.
	sourcedir := filepath.Join(spec.Root, "drone", "src")
	spec.Files = append(spec.Files, &engine.File{
		Path:  sourcedir,
		Mode:  0700,
		IsDir: true,
	})

	// creates the opt directory to hold all scripts.
	spec.Files = append(spec.Files, &engine.File{
		Path:  filepath.Join(spec.Root, "opt"),
		Mode:  0700,
		IsDir: true,
	})

	// creates the netrc file
	if c.Netrc != nil {
		netrcpath := filepath.Join(homedir, ".netrc")
		netrcdata := fmt.Sprintf(
			"machine %s login %s password %s",
			c.Netrc.Machine,
			c.Netrc.Login,
			c.Netrc.Password,
		)
		spec.Files = append(spec.Files, &engine.File{
			Path: netrcpath,
			Mode: 0600,
			Data: []byte(netrcdata),
		})
	}

	// create the default environment variables.
	envs := environ.Combine(
		c.Environ,
		c.Build.Params,
		environ.System(c.System),
		environ.Repo(c.Repo),
		environ.Build(c.Build),
		environ.Stage(c.Stage),
		environ.Link(c.Repo, c.Build, c.System),
		clone.Environ(clone.Config{
			SkipVerify: c.Pipeline.Clone.SkipVerify,
			Trace:      c.Pipeline.Clone.Trace,
			User: clone.User{
				Name:  c.Build.AuthorName,
				Email: c.Build.AuthorEmail,
			},
		}),
		map[string]string{
			"HOME":                homedir,
			"DRONE_HOME":          sourcedir,
			"DRONE_WORKSPACE":     sourcedir,
			"GIT_TERMINAL_PROMPT": "0",
		},
	)

	// create clone step, maybe
	if c.Pipeline.Clone.Disable == false {
		clonepath := filepath.Join(spec.Root, "opt", "clone")
		clonefile := shell.Script(
			clone.Commands(
				clone.Args{
					Branch: c.Build.Target,
					Commit: c.Build.After,
					Ref:    c.Build.Ref,
					Remote: c.Repo.HTTPURL,
				},
			),
		)

		cmd, args := shell.Command()
		spec.Steps = append(spec.Steps, &engine.Step{
			Name:      "clone",
			Args:      append(args, clonepath),
			Command:   cmd,
			Envs:      envs,
			RunPolicy: engine.RunAlways,
			Files: []*engine.File{
				{
					Path: clonepath,
					Mode: 0700,
					Data: []byte(clonefile),
				},
			},
			Secrets:    []*engine.Secret{},
			WorkingDir: sourcedir,
		})
	}

	// create steps
	for _, src := range c.Pipeline.Steps {
		buildpath := filepath.Join(spec.Root, "opt", src.Name)
		buildfile := shell.Script(src.Commands)

		cmd, args := shell.Command()
		dst := &engine.Step{
			Name:      src.Name,
			Args:      append(args, buildpath),
			Command:   cmd,
			Detach:    src.Detach,
			DependsOn: src.DependsOn,
			Envs: environ.Combine(envs,
				environ.Expand(
					convertStaticEnv(src.Environment),
				),
			),
			IgnoreErr:    strings.EqualFold(src.Failure, "ignore"),
			IgnoreStdout: false,
			IgnoreStderr: false,
			RunPolicy:    engine.RunOnSuccess,
			Files: []*engine.File{
				{
					Path: buildpath,
					Mode: 0700,
					Data: []byte(buildfile),
				},
			},
			Secrets:    convertSecretEnv(src.Environment),
			WorkingDir: sourcedir,
		}
		spec.Steps = append(spec.Steps, dst)

		// set the pipeline step run policy. steps run on
		// success by default, but may be optionally configured
		// to run on failure.
		if isRunAlways(src) {
			dst.RunPolicy = engine.RunAlways
		} else if isRunOnFailure(src) {
			dst.RunPolicy = engine.RunOnFailure
		}

		// if the pipeline step has unmet conditions the step is
		// automatically skipped.
		if !src.When.Match(manifest.Match{
			Action:   c.Build.Action,
			Cron:     c.Build.Cron,
			Ref:      c.Build.Ref,
			Repo:     c.Repo.Slug,
			Instance: c.System.Host,
			Target:   c.Build.Deploy,
			Event:    c.Build.Event,
			Branch:   c.Build.Target,
		}) {
			dst.RunPolicy = engine.RunNever
		}
	}

	if isGraph(spec) == false {
		configureSerial(spec)
	} else if c.Pipeline.Clone.Disable == false {
		configureCloneDeps(spec)
	} else if c.Pipeline.Clone.Disable == true {
		removeCloneDeps(spec)
	}

	for _, step := range spec.Steps {
		for _, s := range step.Secrets {
			found, _ := c.Secret.Find(ctx, &secret.Request{
				Name:  s.Name,
				Build: c.Build,
				Repo:  c.Repo,
				Conf:  c.Manifest,
			})
			if found != nil {
				s.Data = []byte(found.Data)
			}
		}
	}

	return spec
}
