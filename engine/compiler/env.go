// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package compiler

import "os"

// default function to get environment variables.
var getenv = os.Getenv

// hostEnviron is a helper function that returns a list of host
// machine variables that should be shared with child processes.
func hostEnviron() map[string]string {
	envs := map[string]string{}
	for _, name := range hostVars {
		if value := getenv(name); value != "" {
			envs[name] = value
		}
	}
	return envs
}
