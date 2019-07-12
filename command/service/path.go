// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// +build !windows

package service

import (
	"os"
	"os/user"
)

// function returns the current user
var getuser = user.Current

// function returns the current uid
var getuid = os.Getuid

// helper function returns the default configuration path
// for the drone configuration.
func configPath() string {
	u, err := getuser()
	if err != nil || getuid() == 0 {
		return "/etc/drone-runner-exec/config"
	}
	return u.HomeDir + "/.drone-runner-exec/config"
}
