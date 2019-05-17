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

func CreateMetric(ctx context.Context, metric *models.Metric) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&metric).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert Metric failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribeMetrics(ctx context.Context, req *pb.DescribeMetricsRequest) ([]*models.Metric, uint64, error) {
	req.MetricId = stringutil.SimplifyStringList(req.MetricId)
	req.MetricName = stringutil.SimplifyStringList(req.MetricName)
	req.Status = stringutil.SimplifyStringList(req.Status)
	req.RsTypeId = stringutil.SimplifyStringList(req.RsTypeId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var rss []*models.Metric
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableMetric)).
		AddQueryOrderDir(req, models.MtColCreateTime).
		BuildFilterConditions(req, models.TableMetric).
		Offset(offset).
		Limit(limit).
		Find(&rss).Error; err != nil {
		logger.Error(ctx, "Describe Metrics failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableMetric)).
		BuildFilterConditions(req, models.TableMetric).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe Metrics count failed: %+v", err)
		return nil, 0, err
	}

	return rss, count, nil
}

func ModifyMetric(ctx context.Context, req *pb.ModifyMetricRequest) (string, error) {
	metricId := req.MetricId

	attributes := make(map[string]interface{})

	if req.MetricName != "" {
		attributes[models.MtColName] = req.MetricName
	}
	if req.MetricParam != "" {
		attributes[models.MtColParam] = req.MetricParam
	}
	if req.Status != "" {
		attributes[models.MtColStatus] = req.Status
	}
	if req.RsTypeId != "" {
		attributes[models.MtColTypeId] = req.RsTypeId
	}

	attributes[models.MtColUpdateTime] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var metric models.Metric
	err := tx.Model(&metric).Where(models.MtColId+" = ?", metricId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update Metric [%s] failed: %+v", metricId, err.Error)
		return "", err.Error
	}

	tx.Commit()
	return metricId, nil
}

func DeleteMetrics(ctx context.Context, metricIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	var metric models.Metric
	err := tx.Model(&metric).Where(models.MtColId+" in (?)", metricIds).Delete(models.Metric{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Delete Metrics failed: %+v", err.Error)
		return nil, err.Error
	}
	tx.Commit()
	return metricIds, nil
}
