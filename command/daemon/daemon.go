// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package daemon

import (
	"context"
	"time"

	"github.com/drone-runners/drone-runner-exec/command/daemon/config"
	"github.com/drone-runners/drone-runner-exec/engine"
	"github.com/drone-runners/drone-runner-exec/engine/resource"
	"github.com/drone-runners/drone-runner-exec/runtime"

	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/handler/router"
	"github.com/drone/runner-go/logger"
	"github.com/drone/runner-go/pipeline/history"
	"github.com/drone/runner-go/pipeline/remote"
	"github.com/drone/runner-go/secret"
	"github.com/drone/runner-go/server"
	"github.com/drone/signal"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gopkg.in/alecthomas/kingpin.v2"
)

var nocontext = context.Background()

func run(*kingpin.ParseContext) error {
	config, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("cannot load configuration")
		return err
	}

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if config.Trace {
		logrus.SetLevel(logrus.TraceLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = signal.WithContextFunc(ctx, func() {
		println("received signal, terminating process")
		cancel()
	})

	cli := client.New(
		config.Client.Address,
		config.Client.Secret,
		config.Client.SkipVerify,
	)
	if config.Client.Dump {
		cli.Dumper = logger.StandardDumper(
			config.Client.DumpBody,
		)
	}
	cli.Logger = logger.Logrus(
		logrus.StandardLogger(),
	)

	// Ping the server and block until a successful connection
	// to the server has been established.
	for {
		err := cli.Ping(ctx, config.Runner.Name)
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		if ctx.Err() != nil {
			break
		}
		if err != nil {
			logrus.WithError(err).
				Errorln("cannot ping the server")
			time.Sleep(time.Second)
		} else {
			logrus.Debugln("successfully pinged the server")
			break
		}
	}

	engine := engine.New()
	remote := remote.New(cli)
	tracer := history.New(remote)

	poller := &runtime.Poller{
		Client: cli,
		Runner: &runtime.Runner{
			Client:   cli,
			Environ:  config.Runner.Environ,
			Machine:  config.Runner.Name,
			Reporter: tracer,
			Secret: secret.External(
				config.Secret.Endpoint,
				config.Secret.Token,
				config.Secret.SkipVerify,
			),
			Execer: runtime.NewExecer(
				tracer,
				remote,
				engine,
				config.Runner.Procs,
			),
		},
		Filter: &client.Filter{
			Kind:    resource.Kind,
			Type:    resource.Type,
			OS:      config.Platform.OS,
			Arch:    config.Platform.Arch,
			Variant: config.Platform.Variant,
			Kernel:  config.Platform.Kernel,
			Labels:  config.Runner.Labels,
		},
	}

	var g errgroup.Group
	if config.Dashboard.Disabled == false {
		server := server.Server{
			Addr: config.Server.Port,
			Handler: router.New(tracer, router.Config{
				Username: config.Dashboard.Username,
				Password: config.Dashboard.Password,
				Realm:    config.Dashboard.Realm,
			}),
		}

		logrus.WithField("addr", config.Server.Port).
			Debugln("starting the server")

		g.Go(func() error {
			return server.ListenAndServe(ctx)
		})
	}

	g.Go(func() error {
		logrus.WithField("capacity", config.Runner.Capacity).
			WithField("endpoint", config.Client.Address).
			WithField("kind", resource.Kind).
			WithField("type", resource.Type).
			Debugln("starting the poller")

		poller.Poll(ctx, config.Runner.Capacity)
		return nil
	})

	err = g.Wait()
	if err != nil {
		logrus.WithError(err).
			Errorln("shutting down the server")
	}
	return err
}

// Register registers the command.
func Register(app *kingpin.Application) {
	app.Command("daemon", "starts the runner daemon").
		Default().Action(run)
}
