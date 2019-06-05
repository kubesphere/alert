package resource_control

import (
	"context"

	aldb "kubesphere.io/alert/pkg/db"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	"kubesphere.io/alert/pkg/pb"
	"kubesphere.io/alert/pkg/util/pbutil"
	"kubesphere.io/alert/pkg/util/stringutil"
)

func DescribeHistories(ctx context.Context, req *pb.DescribeHistoriesRequest) ([]*models.History, uint64, error) {
	req.HistoryId = stringutil.SimplifyStringList(req.HistoryId)
	req.HistoryName = stringutil.SimplifyStringList(req.HistoryName)
	req.Event = stringutil.SimplifyStringList(req.Event)
	req.Content = stringutil.SimplifyStringList(req.Content)
	req.AlertId = stringutil.SimplifyStringList(req.AlertId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var hss []*models.History
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableHistory)).
		AddQueryOrderDir(req, models.HsColCreateTime).
		BuildFilterConditions(req, models.TableHistory).
		Offset(offset).
		Limit(limit).
		Find(&hss).Error; err != nil {
		logger.Error(ctx, "Describe Histories failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableHistory)).
		BuildFilterConditions(req, models.TableHistory).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe Histories count failed: %+v", err)
		return nil, 0, err
	}

	return hss, count, nil
}
