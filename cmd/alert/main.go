// Copyright 2018 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/services/client"
	"kubesphere.io/alert/pkg/services/executor"
	"kubesphere.io/alert/pkg/services/manager"
	"kubesphere.io/alert/pkg/services/watcher"
)

func exitHandler() {
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				ExitFunc()
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
			default:
			}
		}
	}()
}

func ExitFunc() {
	cfg := config.GetInstance()
	switch cfg.App.RunMode {
	case "executor":
		exitFuncExecutor()
	}
	os.Exit(0)
}

var e *executor.Executor

func exitFuncExecutor() {
	e.DeleteAllAlerts()
	time.Sleep(time.Second)
}

func mainFuncExecutor() {
	host, err := os.Hostname()
	if err != nil {
		logger.Error(nil, "Get Host Name error: %s", err)
		return
	}

	e = executor.Init(host)

	if e.CheckExist() {
		logger.Error(nil, "Executor ID already used, exiting...")
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
	}
	e.HeartBoot()
	e.Serve()
}

func mainFuncWatcher() {
	ew := watcher.NewExecutorWatcher()
	ew.Serve()
}

func mainFuncManager() {
	manager.Serve()
}

func mainFuncClient() {
	client.Run()
}

func main() {
	exitHandler()

	cfg := config.GetInstance().LoadConf()
	switch cfg.App.RunMode {
	case "executor":
		mainFuncExecutor()
	case "watcher":
		mainFuncWatcher()
	case "manager":
		mainFuncManager()
	case "client":
		mainFuncClient()
	default:
		logger.Error(nil, "Run mode error, exiting...")
	}
}
