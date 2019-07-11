// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package config

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kelseyhightower/envconfig"
)

// Config stores the system configuration.
type Config struct {
	Debug bool `envconfig:"DRONE_DEBUG"`
	Trace bool `envconfig:"DRONE_TRACE"`

	Logger struct {
		File       string `envconfig:"DRONE_LOG_FILE"`
		MaxAge     int    `envconfig:"DRONE_LOG_FILE_MAX_AGE"     default:"1"`
		MaxBackups int    `envconfig:"DRONE_LOG_FILE_MAX_BACKUPS" default:"1"`
		MaxSize    int    `envconfig:"DRONE_LOG_FILE_MAX_SIZE"    default:"100"`
	}

	Client struct {
		Address    string `ignored:"true"`
		Proto      string `envconfig:"DRONE_RPC_PROTO"  default:"http"`
		Host       string `envconfig:"DRONE_RPC_HOST"   required:"true"`
		Secret     string `envconfig:"DRONE_RPC_SECRET" required:"true"`
		SkipVerify bool   `envconfig:"DRONE_RPC_SKIP_VERIFY"`
		Dump       bool   `envconfig:"DRONE_RPC_DUMP_HTTP"`
		DumpBody   bool   `envconfig:"DRONE_RPC_DUMP_HTTP_BODY"`
	}

	Platform struct {
		OS      string `envconfig:"DRONE_PLATFORM_OS"`
		Arch    string `envconfig:"DRONE_PLATFORM_ARCH"`
		Kernel  string `envconfig:"DRONE_PLATFORM_KERNEL"`
		Variant string `envconfig:"DRONE_PLATFORM_VARIANT"`
	}

	Dashboard struct {
		Disabled bool   `envconfig:"DRONE_UI_DISABLE"`
		Username string `envconfig:"DRONE_UI_USERNAME" default:"admin"`
		Password string `envconfig:"DRONE_UI_PASSWORD" default:"admin"`
		Realm    string `envconfig:"DRONE_UI_REALM"    default:"MyRealm"`
	}

	Server struct {
		Proto string `envconfig:"DRONE_SERVER_PROTO"`
		Host  string `envconfig:"DRONE_SERVER_HOST"`
		Port  string `envconfig:"DRONE_SERVER_PORT" default:":3000"`
		Acme  bool   `envconfig:"DRONE_SERVER_ACME"`
	}

	Runner struct {
		Name     string            `envconfig:"DRONE_RUNNER_NAME"`
		Capacity int               `envconfig:"DRONE_RUNNER_CAPACITY" default:"2"`
		Procs    int64             `envconfig:"DRONE_RUNNER_MAX_PROCS"`
		Labels   map[string]string `envconfig:"DRONE_RUNNER_LABELS"`
		Environ  map[string]string `envconfig:"DRONE_RUNNER_ENVIRON"`
	}

	Limit struct {
		Repos   []string `envconfig:"DRONE_LIMIT_REPOS"`
		Events  []string `envconfig:"DRONE_LIMIT_EVENTS"`
		Trusted bool     `envconfig:"DRONE_LIMIT_TRUSTED"`
	}

	Secret struct {
		Endpoint   string `envconfig:"DRONE_SECRET_PLUGIN_ENDPOINT"`
		Token      string `envconfig:"DRONE_SECRET_PLUGIN_TOKEN"`
		SkipVerify bool   `envconfig:"DRONE_SECRET_PLUGIN_SKIP_VERIFY"`
	}
}

// Load loads the configuration from the environment.
func Load() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return config, err
	}
	if config.Runner.Name == "" {
		config.Runner.Name, _ = os.Hostname()
	}
	if config.Platform.OS == "" {
		config.Platform.OS = runtime.GOOS
	}
	if config.Platform.Arch == "" {
		config.Platform.Arch = runtime.GOARCH
	}
	config.Client.Address = fmt.Sprintf(
		"%s://%s",
		config.Client.Proto,
		config.Client.Host,
	)
	return config, nil
}
