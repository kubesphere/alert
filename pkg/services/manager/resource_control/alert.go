package resource_control

import (
	"context"
	"time"

	aldb "kubesphere.io/alert/pkg/db"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	"kubesphere.io/alert/pkg/pb"
	"kubesphere.io/alert/pkg/util/pbutil"
	"kubesphere.io/alert/pkg/util/stringutil"
)

func UpdateAlertInfo(ctx context.Context, alertId string, runningStatus string) error {
	attributes := make(map[string]interface{})

	attributes["running_status"] = runningStatus
	attributes["update_time"] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var alert models.Alert
	err := tx.Model(&alert).Where(models.AlColId+" = ?", alertId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update alert failed, [%+v]\n", err.Error)
		return err.Error
	}

	tx.Commit()
	return nil
}

func CreateAlert(ctx context.Context, alert *models.Alert) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&alert).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert Alert failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribeAlerts(ctx context.Context, req *pb.DescribeAlertsRequest) ([]*models.Alert, uint64, error) {
	req.AlertId = stringutil.SimplifyStringList(req.AlertId)
	req.AlertName = stringutil.SimplifyStringList(req.AlertName)
	req.RunningStatus = stringutil.SimplifyStringList(req.RunningStatus)
	req.PolicyId = stringutil.SimplifyStringList(req.PolicyId)
	req.RsFilterId = stringutil.SimplifyStringList(req.RsFilterId)
	req.ExecutorId = stringutil.SimplifyStringList(req.ExecutorId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var als []*models.Alert
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableAlert)).
		AddQueryOrderDir(req, models.AlColCreateTime).
		BuildFilterConditions(req, models.TableAlert).
		Offset(offset).
		Limit(limit).
		Find(&als).Error; err != nil {
		logger.Error(ctx, "Describe Alerts failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableAlert)).
		BuildFilterConditions(req, models.TableAlert).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe Alerts count failed: %+v", err)
		return nil, 0, err
	}

	return als, count, nil
}

func UpdateAlert(ctx context.Context, req *pb.ModifyAlertRequest, alertId string, operation string) (string, error) {
	attributes := make(map[string]interface{})

	switch operation {
	case "updating":
		attributes[models.AlColId] = alertId
		if req.AlertName != "" {
			attributes[models.AlColName] = req.AlertName
		}
		attributes[models.AlColDisabled] = req.Disabled
		if req.PolicyId != "" {
			attributes[models.AlColPolicyId] = req.PolicyId
		}
		if req.RsFilterId != "" {
			attributes[models.AlColRsFilterId] = req.RsFilterId
		}

		attributes[models.AlColRunningStatus] = operation
		attributes[models.AlColUpdateTime] = time.Now()
	case "deleting":
		attributes[models.AlColId] = alertId

		attributes[models.AlColRunningStatus] = operation
		attributes[models.AlColUpdateTime] = time.Now()
	}

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var alert models.Alert
	err := tx.Model(&alert).Where(models.AlColId+" = ?", alertId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update Alert [%s] failed: %+v", alertId, err.Error)
		return alertId, err.Error
	}
	tx.Commit()

	return alertId, nil
}
