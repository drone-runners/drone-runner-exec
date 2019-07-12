// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import "gopkg.in/alecthomas/kingpin.v2"

// Register registers the command.
func Register(app *kingpin.Application) {
	cmd := app.Command("service", "manages the runner service")
	registerInstall(cmd)
	registerStart(cmd)
	registerStop(cmd)
	registerUninstall(cmd)
	registerRun(cmd)
}
