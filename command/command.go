// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package command

import (
	"os"

	"github.com/drone-runners/drone-runner-exec/command/compile"
	"github.com/drone-runners/drone-runner-exec/command/daemon"
	"github.com/drone-runners/drone-runner-exec/command/exec"

	"gopkg.in/alecthomas/kingpin.v2"
)

// program version
var version = "0.0.0"

// Command parses the command line arguments and then executes a
// subcommand program.
func Command() {
	app := kingpin.New("drone", "drone exec runner")
	compile.Register(app)
	daemon.Register(app)
	exec.Register(app)

	kingpin.Version(version)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
