// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/koding/multiconfig"

	"kubesphere.io/alert/pkg/constants"
	"kubesphere.io/alert/pkg/logger"
)

type Config struct {
	Log  LogConfig
	Grpc GrpcConfig

	Mysql struct {
		Host     string `default:"139.198.121.96"`
		Port     string `default:"3306"`
		User     string `default:"root"`
		Password string `default:"password"`
		Database string `default:"alert"`
		Disable  bool   `default:"false"`
		LogMode  bool   `default:"true"`
	}

	Etcd struct {
		Endpoints string `default:"139.198.121.96:2379"`
	}

	Queue struct {
		Type string `default:"redis"`
		Addr string `default:"redis://redis.kubesphere-system.svc:6379"`
	}

	App struct {
		Host string `default:"localhost"`
		Port string `default:"9201"`

		ApiHost string `default:"localhost"`
		ApiPort string `default:"9200"`

		NotificationHost string `default:"notification.kubesphere-notification-system.svc:9201"`

		RunMode string `default:"none"`

		AdapterPort string `default:"8080"`
	}
}

var instance *Config

var once sync.Once

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

type GrpcConfig struct {
	ShowErrorCause bool `default:"false"` // show grpc error cause to frontend
}

type LogConfig struct {
	Level string `default:"debug"` // debug, info, warn, error, fatal
}

func (c *Config) PrintUsage() {
	fmt.Fprintf(os.Stdout, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stdout, "\nSupported environment variables:\n")
	e := newLoader(constants.ServiceName)
	e.PrintEnvs(new(Config))
	fmt.Println("")
}

func (c *Config) GetFlagSet() *flag.FlagSet {
	flag.CommandLine.Usage = c.PrintUsage
	return flag.CommandLine
}

func (c *Config) ParseFlag() {
	c.GetFlagSet().Parse(os.Args[1:])
}

func (c *Config) LoadConf() *Config {
	c.ParseFlag()
	config := instance

	m := &multiconfig.DefaultLoader{}
	m.Loader = multiconfig.MultiLoader(newLoader(constants.ServiceName))
	m.Validator = multiconfig.MultiValidator(
		&multiconfig.RequiredValidator{},
	)
	err := m.Load(config)
	if err != nil {
		panic(err)
	}

	logger.Info(nil, "LoadConf: %+v", config)

	return config
}
