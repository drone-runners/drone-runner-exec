// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kardianos/service"
	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/joho/godotenv"
)

// required for windows installation.
var username, password string

func start(*kingpin.ParseContext) error {
	s, err := create()
	if err != nil {
		return err
	}
	return s.Start()
}

func stop(*kingpin.ParseContext) error {
	s, err := create()
	if err != nil {
		return err
	}
	return s.Stop()
}

func install(*kingpin.ParseContext) error {
	f := configPath()
	_, err := os.Stat(f)
	if err != nil {
		return fmt.Errorf("cannot load configuration: %s", f)
	}
	s, err := create()
	if err != nil {
		return err
	}
	return s.Install()
}

func uninstall(*kingpin.ParseContext) error {
	s, err := create()
	if err != nil {
		return err
	}
	return s.Uninstall()
}

func run(*kingpin.ParseContext) error {
	godotenv.Load(envfile)
	s, err := create()
	if err != nil {
		return err
	}
	return s.Run()
}

func createService(name, load string) (service.Service, error) {
	config := &service.Config{
		Name:        name,
		DisplayName: "drone-runner-exec",
		Description: "drone exec runner",
		Arguments:   []string{"service", "run", load},
	}

	switch runtime.GOOS {
	case "darwin":
		config.Option = service.KeyValue{
			"KeepAlive":   true,
			"RunAtLoad":   true,
			"UserService": os.Getuid() != 0,
		}
	case "windows":
		if username != "" {
			config.UserName = username
			config.Option = service.KeyValue{
				"Password": password,
			}
		}
	}

	m := new(manager)
	return service.New(m, config)
}

func create() (service.Service, error) {
	config := &service.Config{
		Name:        "drone-runner-exec",
		DisplayName: "drone-runner-exec",
		Description: "drone exec runner",
		Arguments:   []string{"service", "run", configPath()},
	}

	switch runtime.GOOS {
	case "darwin":
		config.Option = service.KeyValue{
			"KeepAlive":   true,
			"RunAtLoad":   true,
			"UserService": os.Getuid() != 0,
		}
	case "windows":
		if username != "" {
			config.UserName = username
			config.Option = service.KeyValue{
				"Password": password,
			}
		}
	}

	m := new(manager)
	return service.New(m, config)
}

// Register registers the command.
func Register(app *kingpin.Application) {
	cmd := app.Command("service", "manages the runner service")

	sub := cmd.Command("install", "install the service").Action(install)
	sub.Flag("username", "windows account username").Default("").StringVar(&username)
	sub.Flag("password", "windows account password").Default("").StringVar(&password)

	cmd.Command("start", "start the server").Action(start)
	cmd.Command("stop", "stop the service").Action(stop)
	cmd.Command("uninstall", "uninstall the service").Action(uninstall)
	run := cmd.Command("run", "run the service").Hidden().Action(run)
	run.Arg("envfile", "load the environment variable file").Default("").StringVar(&envfile)
}

var envfile string
