// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package compiler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drone-runners/drone-runner-exec/engine"
	"github.com/drone-runners/drone-runner-exec/runtime"

	"github.com/drone/runner-go/clone"
	"github.com/drone/runner-go/environ"
	"github.com/drone/runner-go/environ/provider"
	"github.com/drone/runner-go/manifest"
	"github.com/drone/runner-go/secret"
	"github.com/drone/runner-go/shell"

	"github.com/dchest/uniuri"
	"github.com/gosimple/slug"
)

// random generator function
var random = uniuri.New

// temporary directory function
var tempdir = os.TempDir

// Compiler compiles the Yaml configuration file to an
// intermediate representation optimized for simple execution.
type Compiler struct {
	// Environ provides a set of environment variables that
	// should be added to each pipeline step by default.
	Environ provider.Provider

	// Secret returns a named secret value that can be injected
	// into the pipeline step.
	Secret secret.Provider

	// Root defines the optional build root path, defaults to
	// temp directory.
	Root string

	// Symlinks provides an optional list of symlinks that are
	// created and linked to the pipeline workspace.
	Symlinks map[string]string
}

// Compile compiles the configuration file.
func (c *Compiler) Compile(ctx context.Context, args runtime.CompilerArgs) *engine.Spec {
	spec := new(engine.Spec)

	if c.Root != "" {
		spec.Root = filepath.Join(
			c.Root,
			fmt.Sprintf("drone-%s", random()),
		)
	} else {
		spec.Root = filepath.Join(
			tempdir(),
			fmt.Sprintf("drone-%s", random()),
		)
	}

	pipeline := args.Pipeline
	spec.Platform.OS = pipeline.Platform.OS
	spec.Platform.Arch = pipeline.Platform.Arch
	spec.Platform.Variant = pipeline.Platform.Variant
	spec.Platform.Version = pipeline.Platform.Version

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
	if args.Netrc != nil {
		netrcpath := filepath.Join(homedir, netrc)
		netrcdata := fmt.Sprintf(
			"machine %s login %s password %s",
			args.Netrc.Machine,
			args.Netrc.Login,
			args.Netrc.Password,
		)
		spec.Files = append(spec.Files, &engine.File{
			Path: netrcpath,
			Mode: 0600,
			Data: []byte(netrcdata),
		})
	}

	// create symbolic links
	for source, target := range c.Symlinks {
		spec.Links = append(spec.Links, &engine.Link{
			Source: source,
			Target: filepath.Join(spec.Root, target),
		})
	}

	// list the global environment variables
	globals, _ := c.Environ.List(ctx, &provider.Request{
		Build: args.Build,
		Repo:  args.Repo,
	})

	// create the default environment variables.
	envs := environ.Combine(
		hostEnviron(),
		provider.ToMap(
			provider.FilterUnmasked(globals),
		),
		args.Build.Params,
		environ.Proxy(),
		environ.System(args.System),
		environ.Repo(args.Repo),
		environ.Build(args.Build),
		environ.Stage(args.Stage),
		environ.Link(args.Repo, args.Build, args.System),
		clone.Environ(clone.Config{
			SkipVerify: pipeline.Clone.SkipVerify,
			Trace:      pipeline.Clone.Trace,
			User: clone.User{
				Name:  args.Build.AuthorName,
				Email: args.Build.AuthorEmail,
			},
		}),
		// TODO(bradrydzewski) windows variable HOMEDRIVE
		// TODO(bradrydzewski) windows variable LOCALAPPDATA
		map[string]string{
			"HOME":                homedir,
			"HOMEPATH":            homedir, // for windows
			"USERPROFILE":         homedir, // for windows
			"DRONE_HOME":          sourcedir,
			"DRONE_WORKSPACE":     sourcedir,
			"GIT_TERMINAL_PROMPT": "0",
		},
	)

	// create clone step, maybe
	if pipeline.Clone.Disable == false {
		clonepath := filepath.Join(spec.Root, "opt", "clone"+shell.Suffix)
		repoUrl := args.Repo.HTTPURL
		if repoUrl == "" && args.Repo.SSHURL != "" {
			repoUrl = args.Repo.SSHURL
		}
		clonefile := shell.Script(
			clone.Commands(
				clone.Args{
					Branch: args.Build.Target,
					Commit: args.Build.After,
					Ref:    args.Build.Ref,
					Remote: repoUrl,
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
	for _, src := range pipeline.Steps {
		buildslug := slug.Make(src.Name)
		buildpath := filepath.Join(spec.Root, "opt", buildslug+shell.Suffix)
		buildfile := shell.Script(src.Commands)

		cmd, cmdArgs := shell.Command()
		dst := &engine.Step{
			Name:      src.Name,
			Args:      append(cmdArgs, buildpath),
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
			Action:   args.Build.Action,
			Cron:     args.Build.Cron,
			Ref:      args.Build.Ref,
			Repo:     args.Repo.Slug,
			Instance: args.System.Host,
			Target:   args.Build.Deploy,
			Event:    args.Build.Event,
			Branch:   args.Build.Target,
		}) {
			dst.RunPolicy = engine.RunNever
		}
	}

	if isGraph(spec) == false {
		configureSerial(spec)
	} else if pipeline.Clone.Disable == false {
		configureCloneDeps(spec)
	} else if pipeline.Clone.Disable == true {
		removeCloneDeps(spec)
	}

	// HACK: append masked global variables to secrets
	// this ensures the environment variable values are
	// masked when printed to the console.
	masked := provider.FilterMasked(globals)
	for _, step := range spec.Steps {
		for _, g := range masked {
			step.Secrets = append(step.Secrets, &engine.Secret{
				Name: g.Name,
				Data: []byte(g.Data),
				Mask: g.Mask,
				Env:  g.Name,
			})
		}
	}

	for _, step := range spec.Steps {
		for _, s := range step.Secrets {
			// source secrets from the global secret provider
			// and the repository secret provider.
			provider := secret.Combine(
				args.Secret,
				c.Secret,
			)

			found, _ := provider.Find(ctx, &secret.Request{
				Name:  s.Name,
				Build: args.Build,
				Repo:  args.Repo,
				Conf:  args.Manifest,
			})
			if found != nil {
				s.Data = []byte(found.Data)
			}
		}
	}

	return spec
}
