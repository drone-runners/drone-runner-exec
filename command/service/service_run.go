// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import (
	"github.com/drone-runners/drone-runner-exec/daemon/service"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/joho/godotenv"
)

type runCommand struct{
	config service.Config
}

func (c *runCommand) run(*kingpin.ParseContext) error {
	godotenv.Load(c.config.ConfigFile)
	s, err := service.New(c.config)
	if err != nil {
		return err
	}
	return s.Run()
}

func registerRun(cmd *kingpin.CmdClause) {
	c := new(runCommand)
	s := cmd.Command("run", "run the service").
		Action(c.run).
		Hidden()

	s.Flag("name", "service name").
		Default(service.DefaultName).
		StringVar(&c.config.Name)
	
	s.Flag("desc", "service description").
		Default(service.DefaultDesc).
		StringVar(&c.config.Desc)

	s.Flag("username", "windows account username").
		Default("").
		StringVar(&c.config.Username)

	s.Flag("password", "windows account password").
		Default("").
		StringVar(&c.config.Password)
	
	s.Flag("config", "service configuration file").
		Default(configPath()).
		StringVar(&c.config.ConfigFile)
}
