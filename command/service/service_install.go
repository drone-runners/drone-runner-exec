// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import (
	"fmt"
	"os"

	"github.com/drone-runners/drone-runner-exec/daemon/service"

	"gopkg.in/alecthomas/kingpin.v2"
)

type installCommand struct {
	config service.Config
}

func (c *installCommand) run(*kingpin.ParseContext) error {
	fmt.Printf("read configuration %s\n", c.config.ConfigFile)
	fmt.Printf("installing service %s\n", c.config.Name)
	if _, err := os.Stat(c.config.ConfigFile); err != nil {
		return fmt.Errorf("cannot read configuration", c.config.ConfigFile)
	}
	s, err := service.New(c.config)
	if err != nil {
		return err
	}
	return s.Install()
}

func registerInstall(cmd *kingpin.CmdClause) {
	c := new(installCommand)
	s := cmd.Command("install", "install the service").
		Action(c.run)

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
