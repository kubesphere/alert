package resource_control

import (
	"errors"
	"fmt"
	"time"

	aldb "kubesphere.io/alert/pkg/db"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
)

type AlertDetail struct {
	AlertId            string `gorm:"column:alert_id" json:"alert_id"`
	AlertName          string `gorm:"column:alert_name" json:"alert_name"`
	Disabled           bool   `gorm:"column:disabled" json:"disabled"`
	AlertStatus        string `gorm:"column:alert_status" json:"alert_status"`
	RsTypeName         string `gorm:"column:rs_type_name" json:"rs_type_name"`
	RsTypeParam        string `gorm:"column:rs_type_param" json:"rs_type_param"`
	RsFilterName       string `gorm:"column:rs_filter_name" json:"rs_filter_name"`
	RsFilterParam      string `gorm:"column:rs_filter_param" json:"rs_filter_param"`
	PolicyConfig       string `gorm:"column:policy_config" json:"policy_config"`
	AvailableStartTime string `gorm:"column:available_start_time" json:"available_start_time"`
	AvailableEndTime   string `gorm:"column:available_end_time" json:"available_end_time"`
	NfAddressListId    string `gorm:"column:nf_address_list_id" json:"nf_address_list_id"`
}

type RunnerInfo struct {
	AlertId     string
	AlertStatus string
	UpdateTime  time.Time
}

func QueryAlertDetail(alertId string) (AlertDetail, error) {
	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("alert t1").
		Select("t1.alert_id, t1.alert_name, t1.disabled, t1.alert_status, t3.rs_type_name, t3.rs_type_param, t2.rs_filter_name, t2.rs_filter_param, t4.policy_config, t4.available_start_time, t4.available_end_time, t5.nf_address_list_id").
		Joins("left join resource_filter t2 on t2.rs_filter_id=t1.rs_filter_id").
		Joins("left join resource_type t3 on t3.rs_type_id=t2.rs_type_id").
		Joins("left join policy t4 on t4.policy_id=t1.policy_id").
		Joins("left join action t5 on t5.policy_id=t1.policy_id"))

	dbChain.DB = dbChain.DB.Where("t1.alert_id in (?)", alertId)

	ad := AlertDetail{}

	var ads []*AlertDetail

	err := dbChain.
		Offset(0).
		Limit(100).
		Scan(&ads).
		Error
	if err != nil {
		logger.Error(nil, "Failed to QueryAlertDetail [%v], error: %+v.", alertId, err)
		return ad, err
	}

	if len(ads) <= 0 {
		logger.Error(nil, "QueryAlertDetail [%v] has no match alert.", alertId)
		return ad, errors.New("QueryAlertDetail has no match alert")
	}

	ad = *ads[0]

	return ad, nil
}

func GetAlertInfo(alertId string) models.Alert {
	db := global.GetInstance().GetDB()
	var alert models.Alert
	db.First(&alert, "alert_id = ?", alertId)
	return alert
}

func UpdateAlertInfo(alertId string, executorId string, runningStatus string) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	var alert models.Alert
	err := tx.Model(&alert).Where("alert_id = ?", alertId).Updates(map[string]interface{}{"executor_id": executorId, "running_status": runningStatus, "update_time": time.Now()})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "Update alert failed, [%+v]\n", err.Error)
		return err.Error
	}
	tx.Commit()
	return nil
}

func DeleteAlert(alertId string) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()

	//1. Delete ResourceFilter
	sql := fmt.Sprintf("DELETE rf FROM resource_filter rf, alert al WHERE rf.rs_filter_id=al.rs_filter_id AND al.alert_id IN ('%s')", alertId)
	err := tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlert Delete ResourceFilter failed, [%+v]\n", err.Error)
		return err.Error
	}

	//2. Delete Rule
	sql = fmt.Sprintf("DELETE rl FROM rule rl, alert al WHERE rl.policy_id=al.policy_id AND al.alert_id IN ('%s')", alertId)
	err = tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlert Delete Rules failed, [%+v]\n", err.Error)
		return err.Error
	}

	//3. Delete Action
	sql = fmt.Sprintf("DELETE ac FROM action ac, alert al WHERE ac.policy_id=al.policy_id AND al.alert_id IN ('%s')", alertId)
	err = tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlert Delete Action failed, [%+v]\n", err.Error)
		return err.Error
	}

	//4. Delete Policy
	sql = fmt.Sprintf("DELETE pl FROM policy pl, alert al WHERE pl.policy_id=al.policy_id AND al.alert_id IN ('%s')", alertId)
	err = tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlert Delete Policy failed, [%+v]\n", err.Error)
		return err.Error
	}

	//5. Delete Alert
	var alert models.Alert
	err = tx.Model(&alert).Where("alert_id in (?)", alertId).Delete(models.Alert{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "DeleteAlert Delete Alert failed, [%+v]\n", err.Error)
		return err.Error
	}

	tx.Commit()
	return nil
}

func QueryAlerts(executorId string, runningStatus string) []models.Alert {
	var alerts []models.Alert

	db := global.GetInstance().GetDB()
	db.Where("executor_id = ? AND running_status = ?", executorId, runningStatus).Find(&alerts)

	return alerts
}

func UpdateAlertStatus(runners []RunnerInfo, executorId string) error {
	if len(runners) == 0 {
		return nil
	}

	db := global.GetInstance().GetDB()
	coreAlertStatus := ""
	coreUpdateTime := ""
	where := ""
	for i, runner := range runners {
		coreAlertStatus += fmt.Sprintf("WHEN '%s' THEN '%s' ", runner.AlertId, runner.AlertStatus)
		coreUpdateTime += fmt.Sprintf("WHEN '%s' THEN '%s' ", runner.AlertId, runner.UpdateTime.Format("2006-01-02 15:04:05"))
		if i == 0 {
			where += fmt.Sprintf("'%s'", runner.AlertId)
		} else {
			where += fmt.Sprintf(",'%s'", runner.AlertId)
		}
	}
	sql := fmt.Sprintf("UPDATE `alert` SET `alert_status`= (CASE `alert_id` %s END), `update_time`= (CASE `alert_id` %s END) WHERE `alert_id` IN (%s) AND `executor_id` = '%s' AND `running_status` = 'running'", coreAlertStatus, coreUpdateTime, where, executorId)
	tx := db.Begin()
	err := tx.Exec(sql)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(nil, "UpdateAlertStatus failed, [%+v]\n", err.Error)
		return err.Error
	}
	tx.Commit()
	return nil
}
