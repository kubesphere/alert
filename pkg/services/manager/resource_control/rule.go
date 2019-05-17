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

func CreateRule(ctx context.Context, rule *models.Rule) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&rule).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "CreateRule Insert Rule failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribeRules(ctx context.Context, req *pb.DescribeRulesRequest) ([]*models.Rule, uint64, error) {
	req.RuleId = stringutil.SimplifyStringList(req.RuleId)
	req.RuleName = stringutil.SimplifyStringList(req.RuleName)
	req.Severity = stringutil.SimplifyStringList(req.Severity)
	req.MetricsType = stringutil.SimplifyStringList(req.MetricsType)
	req.ConditionType = stringutil.SimplifyStringList(req.ConditionType)
	req.Thresholds = stringutil.SimplifyStringList(req.Thresholds)
	req.Unit = stringutil.SimplifyStringList(req.Unit)
	req.PolicyId = stringutil.SimplifyStringList(req.PolicyId)
	req.MetricId = stringutil.SimplifyStringList(req.MetricId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var rss []*models.Rule
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableRule)).
		AddQueryOrderDir(req, models.RlColCreateTime).
		BuildFilterConditions(req, models.TableRule).
		Offset(offset).
		Limit(limit).
		Find(&rss).Error; err != nil {
		logger.Error(ctx, "DescribeRules Describe Rules failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableRule)).
		BuildFilterConditions(req, models.TableRule).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "DescribeRules Describe Rules count failed: %+v", err)
		return nil, 0, err
	}

	return rss, count, nil
}

func ModifyRule(ctx context.Context, req *pb.ModifyRuleRequest) (string, error) {
	ruleId := req.RuleId

	attributes := make(map[string]interface{})

	if req.RuleName != "" {
		attributes[models.RlColName] = req.RuleName
	}
	attributes[models.RlColDisabled] = req.Disabled
	attributes[models.RlColMonitorPeriods] = req.MonitorPeriods
	if req.Severity != "" {
		attributes[models.RlColSeverity] = req.Severity
	}
	if req.MetricsType != "" {
		attributes[models.RlColMetricsType] = req.MetricsType
	}
	if req.ConditionType != "" {
		attributes[models.RlColConditionType] = req.ConditionType
	}
	if req.Thresholds != "" {
		attributes[models.RlColThresholds] = req.Thresholds
	}
	if req.Unit != "" {
		attributes[models.RlColUnit] = req.Unit
	}
	attributes[models.RlColConsecutiveCount] = req.ConsecutiveCount
	attributes[models.RlColInhibit] = req.Inhibit

	attributes[models.RlColUpdateTime] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var rule models.Rule
	err := tx.Model(&rule).Where(models.RlColId+" = ?", ruleId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "ModifyRule Update Rule [%s] failed: %+v", ruleId, err.Error)
		return "", err.Error
	}

	tx.Commit()
	return ruleId, nil
}

func DeleteRules(ctx context.Context, ruleIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	var rule models.Rule
	err := tx.Model(&rule).Where(models.RlColId+" in (?)", ruleIds).Delete(models.Rule{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "DeleteRules Delete Rules failed: %+v", err.Error)
		return nil, err.Error
	}
	tx.Commit()
	return ruleIds, nil
}
