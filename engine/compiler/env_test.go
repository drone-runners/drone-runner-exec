// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// +build !windows

package compiler

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHostEnviron(t *testing.T) {
	defer func() {
		getenv = os.Getenv
	}()

	source := map[string]string{
		"USER": "drone",
		"PATH": "/usr/ucb/which",
		"FOO":  "BAR",
	}

	getenv = func(s string) string {
		return source[s]
	}

	want := map[string]string{
		"USER": "drone",
		"PATH": "/usr/ucb/which",
	}

	got := hostEnviron()
	if diff := cmp.Diff(got, want); len(diff) != 0 {
		t.Errorf("Unexpected host environment variables")
		t.Log(diff)
	}
}
