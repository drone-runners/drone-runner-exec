// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// +build windows

package compiler

// netrc filename
const netrc = "_netrc"

// parameters that may be useful to child processes. some
// parameters are required by powershell, without which, will
// return 8009001d powershell error.
var hostVars = []string{
	"ALLUSERSPROFILE",
	"APPDATA",
	"CLIENTNAME",
	"CommonProgramFiles",
	"CommonProgramFiles(x86)",
	"CommonProgramW6432",
	"COMPUTERNAME",
	"ComSpec",
	"DriverData",
	"OS",
	"Path",    // required by powershell
	"PATHEXT", // required by powershell
	"PROCESSOR_ARCHITECTURE",
	"PROCESSOR_IDENTIFIER",
	"PROCESSOR_LEVEL",
	"PROCESSOR_REVISION",
	"ProgramData",
	"ProgramFiles",
	"ProgramFiles(x86)",
	"ProgramW6432",
	"PUBLIC",
	"SystemDrive", // required by powershell
	"SystemRoot",  // required by powershell
	"TEMP",
	"TMP",
	"USERNAME",
	"windir",
}
