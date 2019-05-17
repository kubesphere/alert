// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package executor

import (
	"strings"
	"sync"

	"kubesphere.io/alert/pkg/logger"
	rs "kubesphere.io/alert/pkg/services/executor/resource_control"
)

type Executor struct {
	name              string
	alertReceiver     *AlertReceiver
	runner            *Runner
	aliveReporter     *AliveReporter
	broadcastReceiver *BroadcastReceiver
	healthChecker     *HealthChecker
}

type Runner struct {
	sync.RWMutex
	Map map[string]*AlertRunner
}

func NewExecutor(name string, alertReceiver *AlertReceiver, aliveReporter *AliveReporter, broadcastReceiver *BroadcastReceiver, healthChecker *HealthChecker) *Executor {
	e := &Executor{
		name:              name,
		alertReceiver:     alertReceiver,
		runner:            &Runner{Map: make(map[string]*AlertRunner)},
		aliveReporter:     aliveReporter,
		broadcastReceiver: broadcastReceiver,
		healthChecker:     healthChecker,
	}
	return e
}

func (e *Executor) startRunner(alertId string) bool {
	e.runner.Lock()
	_, ok := e.runner.Map[alertId]
	e.runner.Unlock()
	if ok {
		logger.Error(nil, "Executor startRunner error: alert %s already exists", alertId)
		return false
	}

	//Query DB, check if this alert is in adding or migrating state
	alert := rs.GetAlertInfo(alertId)
	if !((alert.RunningStatus == "adding" || alert.RunningStatus == "migrating") && alert.ExecutorId == "") {
		logger.Error(nil, "Executor startRunner error: alert %s should not be dispatched", alertId)
		return false
	}

	initStatus := alert.RunningStatus

	//Update DB, set this alert to running
	err := rs.UpdateAlertInfo(alertId, e.name, "running")

	if err != nil {
		logger.Error(nil, "Executor startRunner error: update alert "+alertId+" error")
		return false
	}

	var runner = NewAlertRunner(alertId, e.healthChecker.UpdateCh)

	e.runner.Lock()
	e.runner.Map[alertId] = runner
	e.runner.Unlock()

	go runner.Run(initStatus)

	logger.Debug(nil, "Executor startRunner "+alertId+" success")

	return true
}

func (e *Executor) stopRunner(alertId string) bool {
	e.runner.Lock()
	runner, ok := e.runner.Map[alertId]
	e.runner.Unlock()
	if !ok {
		logger.Error(nil, "Executor stopRunner error: runner does not exist")
		return false
	}

	//Query DB, check if this alert is in deleting state
	alert := rs.GetAlertInfo(alertId)
	if !(alert.RunningStatus == "deleting" && alert.ExecutorId == e.name) {
		logger.Error(nil, "Executor stopRunner error: runner alert %s should not be deleted", alertId)
		return false
	}

	//Update DB, delete this alert
	err := rs.DeleteAlert(alertId)

	if err != nil {
		logger.Error(nil, "Executor stopRunner error: delete alert "+alertId+" error")
		return false
	}

	e.runner.Lock()
	delete(e.runner.Map, alertId)
	e.runner.Unlock()

	runner.SignalCh <- "Stop"

	logger.Debug(nil, "Executor stopRunner "+alertId+" success")

	return true
}

func (e *Executor) updateRunner(alertId string) bool {
	e.runner.Lock()
	runner, ok := e.runner.Map[alertId]
	e.runner.Unlock()
	if !ok {
		logger.Error(nil, "Executor updateRunner error: runner does not exist")
		return false
	}

	//Query DB, check if this alert is in updating state
	alert := rs.GetAlertInfo(alertId)
	if !(alert.RunningStatus == "updating" && alert.ExecutorId == e.name) {
		logger.Error(nil, "Executor updateRunner error: runner alert %s should not be updated", alertId)
		return false
	}

	//Update DB, update this alert
	err := rs.UpdateAlertInfo(alertId, e.name, "running")

	if err != nil {
		logger.Error(nil, "Executor updateRunner error: update alert "+alertId+" error")
		return false
	}

	runner.SignalCh <- "Update"

	logger.Debug(nil, "Executor updateRunner "+alertId+" success")

	return true
}

func (e *Executor) commentRunner(alertId string, historyId string) bool {
	e.runner.Lock()
	runner, ok := e.runner.Map[alertId]
	e.runner.Unlock()
	if !ok {
		logger.Error(nil, "Executor commentRunner error: runner does not exist")
		return false
	}

	//Query DB, check if this alert is in running state
	alert := rs.GetAlertInfo(alertId)
	if !(alert.RunningStatus == "running" && alert.ExecutorId == e.name) {
		logger.Error(nil, "Executor commentRunner error: runner alert %s should not be commented", alertId)
		return false
	}

	runner.SignalCh <- "Comment " + historyId

	logger.Debug(nil, "Executor commentRunner "+alertId+" "+historyId+" success")

	return true
}

func (e *Executor) stopAllRunners() {
	e.runner.Lock()
	for alertId, _ := range e.runner.Map {
		e.runner.Map[alertId].SignalCh <- "Stop"
		delete(e.runner.Map, alertId)
	}
	e.runner.Unlock()
}

func (e *Executor) GetName() string {
	return e.name
}

func (e *Executor) GetTaskCount() int {
	e.runner.Lock()
	count := len(e.runner.Map)
	e.runner.Unlock()

	return count
}

func (e *Executor) getRunnerInfo(alertId string) rs.RunnerInfo {
	runnerInfo := rs.RunnerInfo{}

	runner := e.runner.Map[alertId]

	runnerInfo.AlertId = alertId
	runnerInfo.AlertStatus, runnerInfo.UpdateTime = runner.GetAlertStatus()

	return runnerInfo
}

func (e *Executor) GetRunners() []rs.RunnerInfo {
	runners := []rs.RunnerInfo{}
	e.runner.Lock()
	for alertId, _ := range e.runner.Map {
		runners = append(runners, e.getRunnerInfo(alertId))
	}
	e.runner.Unlock()

	return runners
}

func (e *Executor) AddAlert(alertId string) {
	e.startRunner(alertId)
}

func (e *Executor) ProcessAlert(alertId string, operation string) {
	switch operation {
	case "deleting":
		e.stopRunner(alertId)
	case "updating":
		e.updateRunner(alertId)
	default:
		param := strings.Split(operation, " ")
		if len(param) != 2 {
			break
		}
		switch param[0] {
		case "commenting":
			e.commentRunner(alertId, param[1])
		}
	}
}

func (e *Executor) DeleteAllAlerts() {
	e.stopAllRunners()
}

func (e *Executor) TerminateRunner(alertId string) {
	e.runner.Lock()
	runner, ok := e.runner.Map[alertId]
	e.runner.Unlock()
	if !ok {
		logger.Error(nil, "Executor TerminateRunner error: runner does not exist")
		return
	}

	e.runner.Lock()
	delete(e.runner.Map, alertId)
	e.runner.Unlock()

	runner.SignalCh <- "Stop"

	logger.Debug(nil, "Executor TerminateRunner "+alertId+" success")
}

func (e *Executor) Serve() {
	go e.alertReceiver.Serve()
	go e.broadcastReceiver.WatchBroadcast()
	go e.healthChecker.HealthCheck()
	e.aliveReporter.HeartBeat()
}

func (e *Executor) CheckExist() bool {
	return e.aliveReporter.CheckExist()
}

func (e *Executor) HeartBoot() {
	e.aliveReporter.HeartBoot()
}

func Init(name string) *Executor {
	alertReceiver := NewAlertReceiver()
	aliveReporter := NewAliveReporter()
	broadcastReceiver := NewBroadcastReceiver()
	healthChecker := NewHealthChecker()
	executor := NewExecutor(name, alertReceiver, aliveReporter, broadcastReceiver, healthChecker)

	alertReceiver.SetExecutor(executor)
	aliveReporter.SetExecutor(executor)
	broadcastReceiver.SetExecutor(executor)
	healthChecker.SetExecutor(executor)

	return executor
}
