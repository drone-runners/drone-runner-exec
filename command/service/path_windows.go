// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// +build windows

package service

// helper function returns the default configuration path
// for the drone configuration.
func configPath() string {
	return "C:\\Drone\\drone-runner-exec\\config.env"
}
