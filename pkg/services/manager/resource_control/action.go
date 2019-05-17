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

func CreateAction(ctx context.Context, action *models.Action) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&action).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert Action failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribeActions(ctx context.Context, req *pb.DescribeActionsRequest) ([]*models.Action, uint64, error) {
	req.ActionId = stringutil.SimplifyStringList(req.ActionId)
	req.ActionName = stringutil.SimplifyStringList(req.ActionName)
	req.TriggerStatus = stringutil.SimplifyStringList(req.TriggerStatus)
	req.TriggerAction = stringutil.SimplifyStringList(req.TriggerStatus)
	req.PolicyId = stringutil.SimplifyStringList(req.PolicyId)
	req.NfAddressListId = stringutil.SimplifyStringList(req.NfAddressListId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var rss []*models.Action
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableAction)).
		AddQueryOrderDir(req, models.AcColCreateTime).
		BuildFilterConditions(req, models.TableAction).
		Offset(offset).
		Limit(limit).
		Find(&rss).Error; err != nil {
		logger.Error(ctx, "Describe Actions failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableAction)).
		BuildFilterConditions(req, models.TableAction).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe Actions count failed: %+v", err)
		return nil, 0, err
	}

	return rss, count, nil
}

func ModifyAction(ctx context.Context, req *pb.ModifyActionRequest) (string, error) {
	actionId := req.ActionId

	attributes := make(map[string]interface{})

	if req.ActionName != "" {
		attributes[models.AcColName] = req.ActionName
	}
	if req.TriggerStatus != "" {
		attributes[models.AcColTriggerStatus] = req.TriggerStatus
	}
	if req.TriggerAction != "" {
		attributes[models.AcColTriggerAction] = req.TriggerAction
	}
	if req.PolicyId != "" {
		attributes[models.AcColPolicyId] = req.PolicyId
	}
	if req.NfAddressListId != "" {
		attributes[models.AcColNfAddressListId] = req.NfAddressListId
	}

	attributes[models.AcColUpdateTime] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var action models.Action
	err := tx.Model(&action).Where(models.AcColId+" = ?", actionId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update Action [%s] failed: %+v", actionId, err.Error)
		return "", err.Error
	}

	tx.Commit()
	return actionId, nil
}

func DeleteActions(ctx context.Context, actionIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	var action models.Action
	err := tx.Model(&action).Where("action_id in (?) and policy_id not in (select policy_id from policy)", actionIds).Delete(models.Action{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Delete Actions failed: %+v", err.Error)
		return nil, err.Error
	}
	tx.Commit()
	return actionIds, nil
}
