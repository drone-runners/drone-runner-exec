// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package daemon

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kelseyhightower/envconfig"

	"github.com/joho/godotenv"
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
		Username string `envconfig:"DRONE_UI_USERNAME"`
		Password string `envconfig:"DRONE_UI_PASSWORD"`
		Realm    string `envconfig:"DRONE_UI_REALM" default:"MyRealm"`
	}

	Server struct {
		Proto string `envconfig:"DRONE_HTTP_PROTO"`
		Host  string `envconfig:"DRONE_HTTP_HOST"`
		Port  string `envconfig:"DRONE_HTTP_BIND" default:":3000"`
		Acme  bool   `envconfig:"DRONE_ACME_ENABLED"`
		Email string `envconfig:"DRONE_ACME_EMAIL"`
	}

	Runner struct {
		Name     string            `envconfig:"DRONE_RUNNER_NAME"`
		Capacity int               `envconfig:"DRONE_RUNNER_CAPACITY" default:"2"`
		Procs    int64             `envconfig:"DRONE_RUNNER_MAX_PROCS"`
		Labels   map[string]string `envconfig:"DRONE_RUNNER_LABELS"`
		Environ  map[string]string `envconfig:"DRONE_RUNNER_ENVIRON"`
		EnvFile  string            `envconfig:"DRONE_RUNNER_ENVFILE"`
		Path     string            `envconfig:"DRONE_RUNNER_PATH"`
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

// FromEnviron loads the configuration from the environment.
func FromEnviron() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return config, err
	}
	if config.Runner.Name == "" {
		config.Runner.Name, _ = os.Hostname()
	}
	if config.Runner.Environ == nil {
		config.Runner.Environ = map[string]string{}
	}
	if config.Platform.OS == "" {
		config.Platform.OS = runtime.GOOS
	}
	if config.Platform.Arch == "" {
		config.Platform.Arch = runtime.GOARCH
	}
	if config.Dashboard.Password == "" {
		config.Dashboard.Disabled = true
	}
	config.Client.Address = fmt.Sprintf(
		"%s://%s",
		config.Client.Proto,
		config.Client.Host,
	)
	if config.Runner.EnvFile != "" {
		envs, err := godotenv.Read(config.Runner.EnvFile)
		if err != nil {
			return config, err
		}
		for k, v := range envs {
			config.Runner.Environ[k] = v
		}
	}
	if path := config.Runner.Path; path != "" {
		config.Runner.Environ["PATH"] = path
	}
	return config, nil
}
