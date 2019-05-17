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

func CreateResourceFilter(ctx context.Context, rsFilter *models.ResourceFilter) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&rsFilter).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert Resource Filter failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribeResourceFilters(ctx context.Context, req *pb.DescribeResourceFiltersRequest) ([]*models.ResourceFilter, uint64, error) {
	req.RsFilterId = stringutil.SimplifyStringList(req.RsFilterId)
	req.RsFilterName = stringutil.SimplifyStringList(req.RsFilterName)
	req.Status = stringutil.SimplifyStringList(req.Status)
	req.RsTypeId = stringutil.SimplifyStringList(req.RsTypeId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var rfs []*models.ResourceFilter
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableResourceFilter)).
		AddQueryOrderDir(req, models.RfColCreateTime).
		BuildFilterConditions(req, models.TableResourceFilter).
		Offset(offset).
		Limit(limit).
		Find(&rfs).Error; err != nil {
		logger.Error(ctx, "Describe Resource Filters failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableResourceFilter)).
		BuildFilterConditions(req, models.TableResourceFilter).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe Resource Filters count failed: %+v", err)
		return nil, 0, err
	}

	return rfs, count, nil
}

func ModifyResourceFilter(ctx context.Context, req *pb.ModifyResourceFilterRequest) (string, error) {
	rsFilterId := req.RsFilterId

	attributes := make(map[string]interface{})

	if req.RsFilterName != "" {
		attributes[models.RfColName] = req.RsFilterName
	}
	if req.RsFilterParam != "" {
		attributes[models.RfColParam] = req.RsFilterParam
	}
	if req.Status != "" {
		attributes[models.RfColStatus] = req.Status
	}
	if req.RsTypeId != "" {
		attributes[models.RfColTypeId] = req.RsTypeId
	}

	attributes[models.RfColUpdateTime] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var resourceFilter models.ResourceFilter
	err := tx.Model(&resourceFilter).Where(models.RfColId+" = ?", rsFilterId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update Resource Filter [%s] failed: %+v", rsFilterId, err.Error)
		return "", err.Error
	}

	tx.Commit()
	return rsFilterId, nil
}

func DeleteResourceFilters(ctx context.Context, rsFilterIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	var resourceFilter models.ResourceFilter
	err := tx.Model(&resourceFilter).Where("rs_filter_id in (?) and rs_filter_id not in (select rs_filter_id from alert)", rsFilterIds).Delete(models.ResourceFilter{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Delete Resource Filters failed: %+v", err.Error)
		return nil, err.Error
	}
	tx.Commit()
	return rsFilterIds, nil
}
