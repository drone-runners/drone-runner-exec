// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kardianos/service"
)

const (
	// DefaultName is the default service name.
	DefaultName = "drone-runner-exec"

	// DefaultDesc is the default service description.
	DefaultDesc = "Drone Exec Runner"
)

// Config configures the service.
type Config struct {
	Name       string // service name
	Desc       string // service description
	Username   string // service username (windows only)
	Password   string // service password (windows only)
	ConfigFile string // service configuration file path
}

// New creates and configures a new service.
func New(conf Config) (service.Service, error) {
	config := &service.Config{
		Name:        conf.Name,
		DisplayName: conf.Name,
		Description: conf.Desc,
		Arguments:   []string{"service", "run", "--config", conf.ConfigFile},
	}

	switch runtime.GOOS {
	case "darwin":
		// In Mac OS, it is impossible to reliably set the PATH
		// of a LaunchAgent outside the plist file. DRONE_RUNNER_ENVIRON
		// and DRONE_RUNNER_ENVFILE will NOT work. So we use a custom service template.
		nonRootUser := os.Getuid() != 0
		config.Option = service.KeyValue{
			"KeepAlive":   true,
			"RunAtLoad":   true,
			"UserService": nonRootUser,
		}
		if nonRootUser {
			config.Option["LaunchdConfig"] = fmt.Sprintf(launchdConfig, os.Getenv("PATH"))
		}
	case "windows":
		if conf.Username != "" {
			config.UserName = conf.Username
			config.Option = service.KeyValue{
				"Password": conf.Password,
			}
		}
	}

	m := new(manager)
	return service.New(m, config)
}
// launchdConfig is our custom service template.
const launchdConfig = `<?xml version='1.0' encoding='UTF-8'?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN"
"http://www.apple.com/DTDs/PropertyList-1.0.dtd" >
<plist version='1.0'>
<dict>
<key>Label</key><string>{{html .Name}}</string>
<key>EnvironmentVariables</key>
<dict>
	<key>PATH</key>
	<string>%s</string>
</dict>
<key>ProgramArguments</key>
<array>
        <string>{{html .Path}}</string>
{{range .Config.Arguments}}
        <string>{{html .}}</string>
{{end}}
</array>
{{if .UserName}}<key>UserName</key><string>{{html .UserName}}</string>{{end}}
{{if .ChRoot}}<key>RootDirectory</key><string>{{html .ChRoot}}</string>{{end}}
{{if .WorkingDirectory}}<key>WorkingDirectory</key><string>{{html .WorkingDirectory}}</string>{{end}}
<key>SessionCreate</key><{{bool .SessionCreate}}/>
<key>KeepAlive</key><{{bool .KeepAlive}}/>
<key>RunAtLoad</key><{{bool .RunAtLoad}}/>
<key>Disabled</key><false/>
</dict>
</plist>
`

