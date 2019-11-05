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

func CreatePolicy(ctx context.Context, policy *models.Policy) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&policy).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert Policy failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribePolicies(ctx context.Context, req *pb.DescribePoliciesRequest) ([]*models.Policy, uint64, error) {
	req.PolicyId = stringutil.SimplifyStringList(req.PolicyId)
	req.PolicyName = stringutil.SimplifyStringList(req.PolicyName)
	req.PolicyDescription = stringutil.SimplifyStringList(req.PolicyDescription)
	req.Creator = stringutil.SimplifyStringList(req.Creator)
	req.RsTypeId = stringutil.SimplifyStringList(req.RsTypeId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var rss []*models.Policy
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TablePolicy)).
		AddQueryOrderDir(req, models.PlColCreateTime).
		BuildFilterConditions(req, models.TablePolicy).
		Offset(offset).
		Limit(limit).
		Find(&rss).Error; err != nil {
		logger.Error(ctx, "Describe Policies failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TablePolicy)).
		BuildFilterConditions(req, models.TablePolicy).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe Policies count failed: %+v", err)
		return nil, 0, err
	}

	return rss, count, nil
}

func ModifyPolicy(ctx context.Context, req *pb.ModifyPolicyRequest) (string, error) {
	policyId := req.PolicyId

	attributes := make(map[string]interface{})

	if req.PolicyName != "" {
		attributes[models.PlColName] = req.PolicyName
	}
	if req.PolicyDescription != "" {
		attributes[models.PlColDescription] = req.PolicyDescription
	}
	if req.PolicyConfig != "" {
		attributes[models.PlColConfig] = req.PolicyConfig
	}
	if req.Creator != "" {
		attributes[models.PlColCreator] = req.Creator
	}
	/*if req.AvailableStartTime != "" {
		attributes[models.PlColAvailableStartTime] = req.AvailableStartTime
	}
	if req.AvailableEndTime != "" {
		attributes[models.PlColAvailableEndTime] = req.AvailableEndTime
	}*/
	if req.RsTypeId != "" {
		attributes[models.PlColTypeId] = req.RsTypeId
	}

	attributes[models.PlColUpdateTime] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var policy models.Policy
	err := tx.Model(&policy).Where(models.PlColId+" = ?", policyId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update Policy [%s] failed: %+v", policyId, err.Error)
		return "", err.Error
	}

	tx.Commit()
	return policyId, nil
}

func DeletePolicies(ctx context.Context, policyIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()

	//1. Delete Action
	var action models.Action
	err := tx.Model(&action).Where("policy_id in (?) and policy_id not in (select policy_id from alert)", policyIds).Delete(models.Action{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "DeletePolicies Delete Actions failed: %+v", err.Error)
		return nil, err.Error
	}

	//2. Delete Rules
	var rule models.Rule
	err = tx.Model(&rule).Where("policy_id in (?) and policy_id not in (select policy_id from alert)", policyIds).Delete(models.Rule{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "DeletePolicies Delete Rules failed: %+v", err.Error)
		return nil, err.Error
	}

	//3. Delete Policy
	var policy models.Policy
	err = tx.Model(&policy).Where("policy_id in (?) and policy_id not in (select policy_id from alert)", policyIds).Delete(models.Policy{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "DeletePolicies Delete Policies failed: %+v", err.Error)
		return nil, err.Error
	}

	tx.Commit()
	return policyIds, nil
}
