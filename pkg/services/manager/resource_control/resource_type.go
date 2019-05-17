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

func CreateResourceType(ctx context.Context, resourceType *models.ResourceType) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&resourceType).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert ResourceType failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribeResourceTypes(ctx context.Context, req *pb.DescribeResourceTypesRequest) ([]*models.ResourceType, uint64, error) {
	req.RsTypeId = stringutil.SimplifyStringList(req.RsTypeId)
	req.RsTypeName = stringutil.SimplifyStringList(req.RsTypeName)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var rts []*models.ResourceType
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableResourceType)).
		AddQueryOrderDir(req, models.RtColCreateTime).
		BuildFilterConditions(req, models.TableResourceType).
		Offset(offset).
		Limit(limit).
		Find(&rts).Error; err != nil {
		logger.Error(ctx, "Describe ResourceTypes failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableResourceType)).
		BuildFilterConditions(req, models.TableResourceType).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe ResourceTypes count failed: %+v", err)
		return nil, 0, err
	}

	return rts, count, nil
}

func ModifyResourceType(ctx context.Context, req *pb.ModifyResourceTypeRequest) (string, error) {
	rsTypeId := req.RsTypeId

	attributes := make(map[string]interface{})

	if req.RsTypeName != "" {
		attributes[models.RtColName] = req.RsTypeName
	}
	if req.RsTypeParam != "" {
		attributes[models.RtColParam] = req.RsTypeParam
	}

	attributes[models.RtColUpdateTime] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var resourceType models.ResourceType
	err := tx.Model(&resourceType).Where(models.RtColId+" = ?", rsTypeId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update ReourceType [%s] failed: %+v", rsTypeId, err.Error)
		return "", err.Error
	}

	tx.Commit()
	return rsTypeId, nil
}

func DeleteResourceTypes(ctx context.Context, rsTypeIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	var resourceType models.ResourceType
	err := tx.Model(&resourceType).Where(models.RtColId+" in (?)", rsTypeIds).Delete(models.ResourceType{})
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Delete ResourceTypes failed: %+v", err.Error)
		return nil, err.Error
	}
	tx.Commit()
	return rsTypeIds, nil
}
