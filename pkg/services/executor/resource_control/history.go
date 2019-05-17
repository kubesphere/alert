package resource_control

import (
	"context"

	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
)

func CreateHistory(ctx context.Context, history *models.History) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&history).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert History failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}
