// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"context"
	"io"
)

// Engine is the interface that must be implemented by a
// pipeline execution engine.
type Engine interface {
	// Setup the pipeline environment.
	Setup(context.Context, *Spec) error

	// Run runs the pipeine step.
	Run(context.Context, *Spec, *Step, io.Writer) (*State, error)

	// Create creates the pipeline state.
	Create(context.Context, *Spec, *Step) error

	// Start the pipeline step.
	Start(context.Context, *Spec, *Step) error

	// Wait for the pipeline step to complete and returns the completion results.
	Wait(context.Context, *Spec, *Step) (*State, error)

	// Tail the pipeline step logs.
	Tail(context.Context, *Spec, *Step) (io.ReadCloser, error)

	// Destroy the pipeline environment.
	Destroy(context.Context, *Spec) error
}
