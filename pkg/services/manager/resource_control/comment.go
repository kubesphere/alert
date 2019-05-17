package resource_control

import (
	"context"
	"errors"

	aldb "kubesphere.io/alert/pkg/db"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	"kubesphere.io/alert/pkg/pb"
	"kubesphere.io/alert/pkg/util/pbutil"
	"kubesphere.io/alert/pkg/util/stringutil"
)

func GetAlertByHistoryId(historyId string) models.Alert {
	db := global.GetInstance().GetDB()
	var alert models.Alert
	db.First(&alert, "alert_id in (select alert_id from history where history_id = ?)", historyId)
	return alert
}

func CreateComment(ctx context.Context, comment *models.Comment) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&comment).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert Comment failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribeComments(ctx context.Context, req *pb.DescribeCommentsRequest) ([]*models.Comment, uint64, error) {
	req.CommentId = stringutil.SimplifyStringList(req.CommentId)
	req.Addresser = stringutil.SimplifyStringList(req.Addresser)
	req.Content = stringutil.SimplifyStringList(req.Content)
	req.HistoryId = stringutil.SimplifyStringList(req.HistoryId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var cms []*models.Comment
	var count uint64

	if len(req.HistoryId) == 0 {
		logger.Debug(ctx, "Describe Comments history ids not specified")
		return nil, 0, errors.New("history ids not specified")
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableComment)).
		AddQueryOrderDir(req, models.CmColCreateTime).
		BuildFilterConditions(req, models.TableComment).
		Offset(offset).
		Limit(limit).
		Find(&cms).Error; err != nil {
		logger.Error(ctx, "Describe Comments failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableComment)).
		BuildFilterConditions(req, models.TableComment).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe Comments count failed: %+v", err)
		return nil, 0, err
	}

	return cms, count, nil
}
