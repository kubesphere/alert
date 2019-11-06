package executor

import (
	"sync"
	"time"

	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	rs "kubesphere.io/alert/pkg/services/executor/resource_control"
)

type UpdateRequest struct {
	sync.RWMutex
	RequestCount int64
}

type HealthChecker struct {
	executor      *Executor
	UpdateCh      chan string
	RequestStatus UpdateRequest
}

func NewHealthChecker() *HealthChecker {
	hc := &HealthChecker{
		UpdateCh: make(chan string, 10000),
	}

	hc.RequestStatus.RequestCount = 0

	return hc
}

func (hc *HealthChecker) SetExecutor(executor *Executor) {
	hc.executor = executor
}

func (hc *HealthChecker) update() {
	runners := hc.executor.GetRunners()

	//Update runner status
	rs.UpdateAlertStatus(runners, hc.executor.GetName())
}

func (hc *HealthChecker) checkAndUpdate() {
	runners := hc.executor.GetRunners()

	//Update runner status
	rs.UpdateAlertStatus(runners, hc.executor.GetName())

	//Check wild runners
	alerts := rs.QueryAlerts(hc.executor.GetName(), "running")
	wildRunners := difference(runners, alerts)
	if len(wildRunners) > 0 {
		logger.Error(nil, "Wild Runners %v", wildRunners)
	}

	for _, alertId := range wildRunners {
		hc.executor.TerminateRunner(alertId)
	}
}

func (hc *HealthChecker) HealthCheck() {
	timer := time.NewTicker(time.Second * 30)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			hc.checkAndUpdate()
		case <-hc.UpdateCh:
			hc.RequestStatus.Lock()
			hc.RequestStatus.RequestCount = hc.RequestStatus.RequestCount + 1
			hc.RequestStatus.Unlock()
		}
	}
}

func (hc *HealthChecker) UpdateLoop() {
	for {
		hc.RequestStatus.Lock()
		count := hc.RequestStatus.RequestCount
		hc.RequestStatus.Unlock()

		if count > 0 {
			hc.update()
		}

		hc.RequestStatus.Lock()
		hc.RequestStatus.RequestCount = 0
		hc.RequestStatus.Unlock()

		time.Sleep(time.Second)
	}
}

func difference(runners []rs.RunnerInfo, alerts []models.Alert) []string {
	alertmap := map[string]bool{}
	for _, x := range alerts {
		alertmap[x.AlertId] = true
	}
	diff := []string{}
	for _, x := range runners {
		if _, ok := alertmap[x.AlertId]; !ok {
			diff = append(diff, x.AlertId)
		}
	}
	return diff
}
