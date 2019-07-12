// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/drone-runners/drone-runner-exec/daemon"

	"github.com/kardianos/service"
)

var nocontext = context.Background()

// a manager manages the service lifecycle. 
type manager struct{
	cancel context.CancelFunc
}

// Start starts the service in a separate go routine.
func (m *manager) Start(service.Service) error {
	config, err := daemon.FromEnviron()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(nocontext)
	m.cancel = cancel
	go daemon.Run(ctx, config)
	return nil
}

// Stop stops the service.
func (m *manager) Stop(service.Service) error  {
	m.cancel()
	return nil
}