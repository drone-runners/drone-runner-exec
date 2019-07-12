// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package command

import (
	"context"

	"github.com/drone-runners/drone-runner-exec/daemon"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/joho/godotenv"
	"github.com/drone/signal"
)

type daemonCommand struct {
	envfile string
}

func (c *daemonCommand) run(*kingpin.ParseContext) error {
	// load environment variables from file.
	godotenv.Load(c.envfile)

	// load the configuration from the environment.
	config, err := daemon.FromEnviron()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(nocontext)
	defer cancel()

	// listen for termination signals to gracefully shutdown
	// the runner daemon.
	ctx = signal.WithContextFunc(ctx, func() {
		println("received signal, terminating process")
		cancel()
	})

	return daemon.Run(ctx, config)
}

func registerDaemon(app *kingpin.Application) {
	c := new(daemonCommand)

	cmd := app.Command("daemon", "starts the runner daemon").
		Default().
		Action(c.run)

	cmd.Arg("envfile", "load the environment variable file").
		Default("").
		StringVar(&c.envfile)
}
