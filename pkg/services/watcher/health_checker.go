package watcher

import (
	rs "kubesphere.io/alert/pkg/services/watcher/resource_control"
	"time"
)

type HealthChecker struct {
	executorWatcher *ExecutorWatcher
}

const (
	DefaultTimeOut = 60 * time.Second
)

func NewHealthChecker() *HealthChecker {
	hc := &HealthChecker{}

	return hc
}

func (hc *HealthChecker) SetExecutorWatcher(executorWatcher *ExecutorWatcher) {
	hc.executorWatcher = executorWatcher
}

func (hc *HealthChecker) doHealthCheck() {
	hc.processTimeoutRunningAlerts()
	hc.processTimeoutAddingAlerts()
	hc.processTimeoutUpdatingAlerts()
	hc.processTimeoutMigratingAlerts()
	hc.processTimeoutDeletingAlerts()
}

func (hc *HealthChecker) HealthCheck() {
	timer := time.NewTicker(time.Second * 30)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			hc.doHealthCheck()
		}
	}
}

func (hc *HealthChecker) processTimeoutRunningAlerts() {
	alerts := rs.GetTimeoutAlerts("running", DefaultTimeOut)
	if len(alerts) == 0 {
		return
	}

	executorCount := hc.executorWatcher.GetExecutorCount()

	for _, alert := range alerts {
		rowsAffected := rs.UpdateAlertByAlertIdWithTimeOut(alert.AlertId, "running", "migrating", DefaultTimeOut)
		if rowsAffected > 0 && executorCount > 0 {
			hc.executorWatcher.alertQueue.WriteBackAlert(alert.AlertId)
		}
	}
}

func (hc *HealthChecker) processTimeoutAddingAlerts() {
	alerts := rs.GetTimeoutAlerts("adding", DefaultTimeOut)
	if len(alerts) == 0 {
		return
	}

	executorCount := hc.executorWatcher.GetExecutorCount()

	for _, alert := range alerts {
		err := rs.UpdateAlertByAlertId(alert.AlertId, "adding", "adding")
		if err == nil && executorCount > 0 {
			hc.executorWatcher.alertQueue.WriteBackAlert(alert.AlertId)
		}
	}
}

func (hc *HealthChecker) processTimeoutUpdatingAlerts() {
	alerts := rs.GetTimeoutAlerts("updating", DefaultTimeOut)
	if len(alerts) == 0 {
		return
	}

	executorCount := hc.executorWatcher.GetExecutorCount()

	for _, alert := range alerts {
		err := rs.UpdateAlertByAlertId(alert.AlertId, "updating", "adding")
		if err == nil && executorCount > 0 {
			hc.executorWatcher.alertQueue.WriteBackAlert(alert.AlertId)
		}
	}
}

func (hc *HealthChecker) processTimeoutMigratingAlerts() {
	alerts := rs.GetTimeoutAlerts("migrating", DefaultTimeOut)
	if len(alerts) == 0 {
		return
	}

	executorCount := hc.executorWatcher.GetExecutorCount()

	for _, alert := range alerts {
		err := rs.UpdateAlertByAlertId(alert.AlertId, "migrating", "migrating")
		if err == nil && executorCount > 0 {
			hc.executorWatcher.alertQueue.WriteBackAlert(alert.AlertId)
		}
	}
}

func (hc *HealthChecker) processTimeoutDeletingAlerts() {
	alerts := rs.GetTimeoutAlerts("deleting", DefaultTimeOut)
	if len(alerts) == 0 {
		return
	}

	alertIds := []string{}
	for _, alert := range alerts {
		alertIds = append(alertIds, alert.AlertId)
	}

	rs.DeleteAlerts(alertIds)
}
