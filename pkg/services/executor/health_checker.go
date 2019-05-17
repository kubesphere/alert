package executor

import (
	"time"

	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	rs "kubesphere.io/alert/pkg/services/executor/resource_control"
)

type HealthChecker struct {
	executor *Executor
	UpdateCh chan string
}

func NewHealthChecker() *HealthChecker {
	hc := &HealthChecker{
		UpdateCh: make(chan string, 10000),
	}

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
			hc.update()
		}
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
