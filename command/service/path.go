// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// +build !windows
// +build !darwin

package service

import (
	"os"
	"os/user"
)

// helper function returns the default configuration path
// for the drone configuration.
func configPath() string {
	u, err := user.Current()
	if err != nil || os.Getuid() == 0 {
		return "/etc/drone-runner-exec/config"
	}
	return u.HomeDir + "/.drone-runner-exec/config"
}
