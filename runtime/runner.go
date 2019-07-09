// Code generated automatically. DO NOT EDIT.

// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/drone-runners/drone-runner-exec/engine"
	"github.com/drone-runners/drone-runner-exec/engine/compiler"
	"github.com/drone-runners/drone-runner-exec/engine/resource"

	"github.com/drone/drone-go/drone"
	"github.com/drone/envsubst"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/environ"
	"github.com/drone/runner-go/manifest"
	"github.com/drone/runner-go/pipeline"
	"github.com/drone/runner-go/secret"

	"github.com/sirupsen/logrus"
)

// Runnner runs the pipeline.
type Runner struct {
	// Client is the remote client responsible for interacting
	// with the central server.
	Client client.Client

	// Execer is responsible for executing intermediate
	// representation of the pipeline and returns its results.
	Execer Execer

	// Reporter reports pipeline status back to the remote
	// server.
	Reporter pipeline.Reporter

	// Environ provides custom, global environment variables
	// that are added to every pipeline step.
	Environ map[string]string

	// Machine provides the runner with the name of the host
	// machine executing the pipeline.
	Machine string

	// Match is an optional function that returns true if the
	// repository or build match user-defined criteria. This is
	// intended as a security measure to prevent a runner from
	// processing an unwanted pipeline.
	Match func(*drone.Repo, *drone.Build) bool

	// Secret provides the compiler with secrets.
	Secret secret.Provider
}

// Run runs the pipeline stage.
func (s *Runner) Run(ctx context.Context, stage *drone.Stage) error {
	logger := logrus.WithFields(logrus.Fields{
		"stage.id":     stage.ID,
		"stage.name":   stage.Name,
		"stage.number": stage.Number,
	})

	logger.Debug("received stage")

	// delivery to a single agent is not guaranteed, which means
	// we need confirm receipt. The first agent that confirms
	// receipt of the stage can assume ownership.

	stage.Machine = s.Machine
	err := s.Client.Accept(ctx, stage)
	if err != nil {
		logger.WithError(err).Error("cannot accept stage")
		return err
	}

	logger.Debug("accepted stage")

	data, err := s.Client.Detail(ctx, stage)
	if err != nil {
		logger.WithError(err).Error("cannot get stage details")
		return err
	}

	logger = logrus.WithFields(logrus.Fields{
		"repo.id":        data.Repo.ID,
		"repo.namespace": data.Repo.Namespace,
		"repo.name":      data.Repo.Name,
		"build.id":       data.Build.ID,
		"build.number":   data.Build.Number,
	})

	logger.Debug("retrieved stage details")

	ctxdone, cancel := context.WithCancel(ctx)
	defer cancel()

	timeout := time.Duration(data.Repo.Timeout) * time.Minute
	ctxtimeout, cancel := context.WithTimeout(ctxdone, timeout)
	defer cancel()

	ctxcancel, cancel := context.WithCancel(ctxtimeout)
	defer cancel()

	// next we opens a connection to the server to watch for
	// cancellation requests. If a build is cancelled the running
	// stage should also be cancelled.
	go func() {
		done, _ := s.Client.Watch(ctxdone, data.Build.ID)
		if done {
			cancel()
			logger.Debugln("received cancellation")
		} else {
			logger.Debugln("done listening for cancellations")
		}
	}()

	envs := environ.Combine(
		s.Environ,
		environ.System(data.System),
		environ.Repo(data.Repo),
		environ.Build(data.Build),
		environ.Stage(stage),
		environ.Link(data.Repo, data.Build, data.System),
		data.Build.Params,
	)

	// string substitution function ensures that string
	// replacement variables are escaped and quoted if they
	// contain a newline character.
	subf := func(k string) string {
		v := envs[k]
		if strings.Contains(v, "\n") {
			v = fmt.Sprintf("%q", v)
		}
		return v
	}

	state := &pipeline.State{
		Build:  data.Build,
		Stage:  stage,
		Repo:   data.Repo,
		System: data.System,
	}

	// evaluates whether or not the agent can process the
	// pipeline. An agent may choose to reject a repository
	// or build for security reasons.
	if s.Match != nil && s.Match(data.Repo, data.Build) == false {
		logger.Error("cannot process stage, access denied")
		state.FailAll(errors.New("cannot process stage. access denied"))
		return s.Reporter.ReportStage(noContext, state)
	}

	// evaluates string replacement expressions and returns an
	// update configuration file string.
	config, err := envsubst.Eval(string(data.Config.Data), subf)
	if err != nil {
		logger.WithError(err).Error("cannot emulate bash substitution")
		state.FailAll(err)
		return s.Reporter.ReportStage(noContext, state)
	}

	// parse the yaml configuration file.
	manifest, err := manifest.ParseString(config)
	if err != nil {
		logger.WithError(err).Error("cannot parse configuration file")
		state.FailAll(err)
		return s.Reporter.ReportStage(noContext, state)
	}

	// find the named stage in the yaml configuration file.
	resource, err := resource.Lookup(stage.Name, manifest)
	if err != nil {
		logger.WithError(err).Error("cannot find pipeline resource")
		state.FailAll(err)
		return s.Reporter.ReportStage(noContext, state)
	}

	secrets := secret.Combine(
		secret.Static(data.Secrets),
		secret.Encrypted(),
		s.Secret,
	)

	// compile the yaml configuration file to an intermediate
	// representation, and then
	comp := &compiler.Compiler{
		Pipeline: resource,
		Manifest: manifest,
		Environ:  s.Environ,
		Build:    data.Build,
		Stage:    stage,
		Repo:     data.Repo,
		System:   data.System,
		Netrc:    data.Netrc,
		Secret:   secrets,
	}

	spec := comp.Compile(ctx)
	for _, src := range spec.Steps {
		// steps that are skipped are ignored and are not stored
		// in the drone database, nor displayed in the UI.
		if src.RunPolicy == engine.RunNever {
			continue
		}
		stage.Steps = append(stage.Steps, &drone.Step{
			Name:      src.Name,
			Number:    len(stage.Steps) + 1,
			StageID:   stage.ID,
			Status:    drone.StatusPending,
			ErrIgnore: src.IgnoreErr,
		})
	}

	stage.Started = time.Now().Unix()
	stage.Status = drone.StatusRunning
	if err := s.Client.Update(ctx, stage); err != nil {
		logger.WithError(err).Error("cannot update stage")
		return err
	}

	logger.Debug("updated stage to running")

	err = s.Execer.Exec(ctxcancel, spec, state)
	if err != nil {
		logger.WithError(err).Debug("stage failed")
		return err
	}
	logger.Debug("updated stage to complete")
	return nil
}
