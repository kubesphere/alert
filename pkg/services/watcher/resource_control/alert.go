package resource_control

import (
	"time"

	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
)

func GetMigrationAlerts(executorId string) []models.Alert {
	db := global.GetInstance().GetDB()
	var alerts []models.Alert
	db.Where("executor_id = ? AND running_status = 'running'", executorId).Find(&alerts)
	return alerts
}

func UpdateAlertByExecutorId(executorId string, oldRunningStatus string, newRunningStatus string) error {
	attributes := make(map[string]interface{})

	attributes["executor_id"] = ""
	attributes["running_status"] = newRunningStatus
	attributes["update_time"] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()
	var alert models.Alert
	err := tx.Model(&alert).Where("executor_id = ? AND running_status = ?", executorId, oldRunningStatus).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "UpdateAlertByExecutorId failed, [%+v]\n", err.Error)
		return err.Error
	}
	tx.Commit()
	return nil
}

func GetTimeoutAlerts(runningStatus string, timeOut time.Duration) []models.Alert {
	db := global.GetInstance().GetDB()
	var alerts []models.Alert
	db.Where("running_status = ? AND update_time < ?", runningStatus, time.Now().Add(-timeOut)).Find(&alerts)
	return alerts
}

func UpdateAlertByAlertId(alertId string, oldStatus string, newStatus string) error {
	attributes := make(map[string]interface{})

	attributes["executor_id"] = ""
	attributes["running_status"] = newStatus
	attributes["update_time"] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()
	var alert models.Alert
	err := tx.Model(&alert).Where("alert_id = ? and running_status = ?", alertId, oldStatus).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "UpdateAlertByAlertId failed, [%+v]\n", err.Error)
		return err.Error
	}
	tx.Commit()
	return nil
}

func DeleteAlertByAlertId(alertId string) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()

	//1. Delete ResourceFilter
	var resourceFilter models.ResourceFilter
	err := tx.Model(&resourceFilter).Where("rs_filter_id in (select rs_filter_id from alert where alert_id in (?))", alertId).Delete(models.ResourceFilter{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete ResourceFilter failed, [%+v]\n", err.Error)
		return err.Error
	}

	//2. Delete Rules
	var rule models.Rule
	err = tx.Model(&rule).Where("policy_id in (select policy_id from policy where policy_id in (select policy_id from alert where alert_id in (?)))", alertId).Delete(models.Rule{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete Rules failed, [%+v]\n", err.Error)
		return err.Error
	}

	//3. Delete Action
	var action models.Action
	err = tx.Model(&action).Where("policy_id in (select policy_id from policy where policy_id in (select policy_id from alert where alert_id in (?)))", alertId).Delete(models.Action{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete Action failed, [%+v]\n", err.Error)
		return err.Error
	}

	//4. Delete Policy
	var policy models.Policy
	err = tx.Model(&policy).Where("policy_id in (select policy_id from alert where alert_id in (?))", alertId).Delete(models.Policy{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete Policy failed, [%+v]\n", err.Error)
		return err.Error
	}

	//5. Delete Alert
	var alert models.Alert
	err = tx.Model(&alert).Where("alert_id in (?)", alertId).Delete(models.Alert{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete Alert failed, [%+v]\n", err.Error)
		return err.Error
	}

	tx.Commit()
	return nil
}
