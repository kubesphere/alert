package resource_control

import (
	"fmt"
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

func DeleteAlerts(alertIds []string) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()

	//1. Prepare where collect
	where := ""
	for i, alertId := range alertIds {
		if i == 0 {
			where += fmt.Sprintf("'%s'", alertId)
		} else {
			where += fmt.Sprintf(",'%s'", alertId)
		}
	}

	//2. Delete ResourceFilters
	sql := fmt.Sprintf("DELETE rf FROM resource_filter rf, alert al WHERE rf.rs_filter_id=al.rs_filter_id AND al.alert_id IN (%s)", where)
	err := tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete ResourceFilter failed, [%+v]\n", err.Error)
		return err.Error
	}

	//3. Delete Rules
	sql = fmt.Sprintf("DELETE rl FROM rule rl, alert al WHERE rl.policy_id=al.policy_id AND al.alert_id IN (%s)", where)
	err = tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete Rules failed, [%+v]\n", err.Error)
		return err.Error
	}

	//4. Delete Actions
	sql = fmt.Sprintf("DELETE ac FROM action ac, alert al WHERE ac.policy_id=al.policy_id AND al.alert_id IN (%s)", where)
	err = tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete Action failed, [%+v]\n", err.Error)
		return err.Error
	}

	//5. Delete Policies
	sql = fmt.Sprintf("DELETE pl FROM policy pl, alert al WHERE pl.policy_id=al.policy_id AND al.alert_id IN (%s)", where)
	err = tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete Policy failed, [%+v]\n", err.Error)
		return err.Error
	}

	//6. Delete Alerts
	var alert models.Alert
	err = tx.Model(&alert).Where("alert_id in (?)", alertIds).Delete(models.Alert{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlertByAlertId Delete Alert failed, [%+v]\n", err.Error)
		return err.Error
	}

	tx.Commit()
	return nil
}
