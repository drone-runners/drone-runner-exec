// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// +build !windows

package compiler

// netrc filename
const netrc = ".netrc"

// parameters that may be useful or required by child processes
// to successfully execute.
var hostVars = []string{
	"PATH",
	"USER",
}
