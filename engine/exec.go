// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/drone/runner-go/environ"
	"github.com/drone/runner-go/logger"
)

// New returns a new engine.
func New() Engine {
	return new(engine)
}

type engine struct{}

// Setup the pipeline environment.
func (e *engine) Setup(ctx context.Context, spec *Spec) error {
	err := os.MkdirAll(spec.Root, 0777)
	if err != nil {
		return err
	}

	// creates folders
	for _, file := range spec.Files {
		if file.IsDir == false {
			continue
		}
		err = os.MkdirAll(file.Path, 0700)
		if err != nil {
			logger.FromContext(ctx).
				WithError(err).
				Error("cannot create working directory")
			return err
		}
	}

	// creates files
	for _, file := range spec.Files {
		if file.IsDir == true {
			continue
		}
		err = ioutil.WriteFile(file.Path, file.Data, os.FileMode(file.Mode))
		if err != nil {
			logger.FromContext(ctx).
				WithError(err).
				Error("cannot write file")
			return err
		}
	}

	// creates step files
	for _, step := range spec.Steps {
		for _, file := range step.Files {
			if file.IsDir == true {
				continue
			}
			err = ioutil.WriteFile(file.Path, file.Data, os.FileMode(file.Mode))
			if err != nil {
				logger.FromContext(ctx).
					WithError(err).
					Error("cannot write file")
				return err
			}
		}
	}

	return nil
}

// Destroy the pipeline environment.
func (e *engine) Destroy(ctx context.Context, spec *Spec) error {
	return os.RemoveAll(spec.Root)
}

// Run runs the pipeline step.
func (e *engine) Run(ctx context.Context, spec *Spec, step *Step, output io.Writer) (*State, error) {
	cmd := exec.CommandContext(ctx, step.Command, step.Args...)
	cmd.Env = environ.Slice(step.Envs)
	cmd.Dir = step.WorkingDir
	cmd.Stdout = output
	cmd.Stderr = output

	for _, secret := range step.Secrets {
		s := fmt.Sprintf("%s=%s", secret.Env, string(secret.Data))
		cmd.Env = append(cmd.Env, s)
	}

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	log := logger.FromContext(ctx)
	log = log.WithField("process.pid", cmd.Process.Pid)
	log.Debug("process started")

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err = <-done:
	case <-ctx.Done():
		cmd.Process.Kill()

		log.Debug("process killed")
		return nil, ctx.Err()
	}

	state := &State{
		ExitCode:  0,
		Exited:    true,
		OOMKilled: false,
	}
	if err != nil {
		state.ExitCode = 255
	}
	if exiterr, ok := err.(*exec.ExitError); ok {
		state.ExitCode = exiterr.ExitCode()
	}

	log.WithField("process.exit", state.ExitCode).
		Debug("process finished")
	return state, err
}

type nilReader struct{}

func (*nilReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

//
// Not Implemented
//

// Create creates the pipeline step.
func (e *engine) Create(ctx context.Context, spec *Spec, step *Step) error {
	return nil // no-op for bash implementation
}

// Start the pipeline step.
func (e *engine) Start(context.Context, *Spec, *Step) error {
	return nil // no-op for bash implementation
}

// Wait for the pipeline step to complete and returns the completion results.
func (e *engine) Wait(context.Context, *Spec, *Step) (*State, error) {
	return nil, nil // no-op for bash implementation
}

// Tail the pipeline step logs.
func (e *engine) Tail(context.Context, *Spec, *Step) (io.ReadCloser, error) {
	return nil, nil // no-op for bash implementation
}
