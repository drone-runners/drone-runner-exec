// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package service

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type installer struct {
	Name   string // service name
	Config string // service configuration path
}

func (c *installer) start(*kingpin.ParseContext) error {
	s, err := create()
	if err != nil {
		return err
	}
	return s.Start()
}

func (c *installer) stop(*kingpin.ParseContext) error {
	s, err := create()
	if err != nil {
		return err
	}
	return s.Stop()
}

func (c *installer) run(*kingpin.ParseContext) error {
	godotenv.Load(c.Config)
	s, err := create()
	if err != nil {
		return err
	}
	return s.Run()
}

func (c *installer) install(*kingpin.ParseContext) error {
	// check configuration file exists
	// validate configuration file
	service, err := createService(c.Name, c.Config)
	if err != nil {
		return err
	}
	return service.Install()
}

func (c *installer) uninstall(*kingpin.ParseContext) error {
	service, err := createService(c.Name, c.Config)
	if err != nil {
		return err
	}
	return service.Uninstall()
}