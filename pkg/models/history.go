package models

import (
	"time"

	"kubesphere.io/alert/pkg/pb"
	"kubesphere.io/alert/pkg/util/idutil"
	"kubesphere.io/alert/pkg/util/pbutil"
)

type History struct {
	HistoryId      string    `gorm:"column:history_id" json:"history_id"`
	HistoryName    string    `gorm:"column:history_name" json:"history_name"`
	Event          string    `gorm:"column:event" json:"event"`
	Content        string    `gorm:"column:content" json:"content"`
	NotificationId string    `gorm:"column:notification_id" json:"notification_id"`
	CreateTime     time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime     time.Time `gorm:"column:update_time" json:"update_time"`
	AlertId        string    `gorm:"column:alert_id" json:"alert_id"`
	RuleId         string    `gorm:"column:rule_id" json:"rule_id"`
	ResourceName   string    `gorm:"column:resource_name" json:"resource_name"`
}

//table name
const (
	TableHistory = "history"
)

const (
	HistoryIdPrefix = "hs-"
)

//field name
//Hs is short for history.
const (
	HsColId             = "history_id"
	HsColName           = "history_name"
	HsColEvent          = "event"
	HsColContent        = "content"
	HsColNotificationId = "notification_id"
	HsColCreateTime     = "create_time"
	HsColUpdateTime     = "update_time"
	HsColAlertId        = "alert_id"
	HsColRuleId         = "rule_id"
	HsColResourceName   = "resource_name"
)

func NewHistoryId() string {
	return idutil.GetUuid(HistoryIdPrefix)
}

func NewHistory(historyName string, event string, content string, notificationId string, alertId string, ruleId string, resourceName string) *History {
	history := &History{
		HistoryId:      NewHistoryId(),
		HistoryName:    historyName,
		Event:          event,
		Content:        content,
		NotificationId: notificationId,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		AlertId:        alertId,
		RuleId:         ruleId,
		ResourceName:   resourceName,
	}
	return history
}

func HistoryToPb(history *History) *pb.History {
	pbHistory := pb.History{}
	pbHistory.HistoryId = history.HistoryId
	pbHistory.HistoryName = history.HistoryName
	pbHistory.Event = history.Event
	pbHistory.Content = history.Content
	pbHistory.NotificationId = history.NotificationId
	pbHistory.CreateTime = pbutil.ToProtoTimestamp(history.CreateTime)
	pbHistory.UpdateTime = pbutil.ToProtoTimestamp(history.UpdateTime)
	pbHistory.AlertId = history.AlertId
	pbHistory.RuleId = history.RuleId
	pbHistory.ResourceName = history.ResourceName
	return &pbHistory
}

func ParseHsSet2PbSet(inHss []*History) []*pb.History {
	var pbHss []*pb.History
	for _, inHs := range inHss {
		pbHs := HistoryToPb(inHs)
		pbHss = append(pbHss, pbHs)
	}
	return pbHss
}

type HistoryDetail struct {
	HistoryId          string    `gorm:"column:history_id" json:"history_id"`
	HistoryName        string    `gorm:"column:history_name" json:"history_name"`
	RuleId             string    `gorm:"column:rule_id" json:"rule_id"`
	RuleName           string    `gorm:"column:rule_name" json:"rule_name"`
	Event              string    `gorm:"column:event" json:"event"`
	NotificationId     string    `gorm:"column:notification_id" json:"notification_id"`
	NotificationStatus string    `gorm:"column:notification_status" json:"notification_status"`
	Severity           string    `gorm:"column:severity" json:"severity"`
	RsTypeName         string    `gorm:"column:rs_type_name" json:"rs_type_name"`
	RsFilterName       string    `gorm:"column:rs_filter_name" json:"rs_filter_name"`
	MetricName         string    `gorm:"column:metric_name" json:"metric_name"`
	ConditionType      string    `gorm:"column:condition_type" json:"condition_type"`
	Thresholds         string    `gorm:"column:thresholds" json:"thresholds"`
	Unit               string    `gorm:"column:unit" json:"unit"`
	AlertName          string    `gorm:"column:alert_name" json:"alert_name"`
	RsFilterParam      string    `gorm:"column:rs_filter_param" json:"rs_filter_param"`
	ResourceName       string    `gorm:"column:resource_name" json:"resource_name"`
	CreateTime         time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time" json:"update_time"`
}

func HistoryDetailToPb(historyDetail *HistoryDetail) *pb.HistoryDetail {
	pbHistoryDetail := pb.HistoryDetail{}
	pbHistoryDetail.HistoryId = historyDetail.HistoryId
	pbHistoryDetail.HistoryName = historyDetail.HistoryName
	pbHistoryDetail.RuleId = historyDetail.RuleId
	pbHistoryDetail.RuleName = historyDetail.RuleName
	pbHistoryDetail.Event = historyDetail.Event
	pbHistoryDetail.NotificationId = historyDetail.NotificationId
	pbHistoryDetail.NotificationStatus = historyDetail.NotificationStatus
	pbHistoryDetail.Severity = historyDetail.Severity
	pbHistoryDetail.RsTypeName = historyDetail.RsTypeName
	pbHistoryDetail.RsFilterName = historyDetail.RsFilterName
	pbHistoryDetail.MetricName = historyDetail.MetricName
	pbHistoryDetail.ConditionType = historyDetail.ConditionType
	pbHistoryDetail.Thresholds = historyDetail.Thresholds
	pbHistoryDetail.Unit = historyDetail.Unit
	pbHistoryDetail.AlertName = historyDetail.AlertName
	pbHistoryDetail.RsFilterParam = historyDetail.RsFilterParam
	pbHistoryDetail.ResourceName = historyDetail.ResourceName
	pbHistoryDetail.CreateTime = pbutil.ToProtoTimestamp(historyDetail.CreateTime)
	pbHistoryDetail.UpdateTime = pbutil.ToProtoTimestamp(historyDetail.UpdateTime)
	return &pbHistoryDetail
}

func ParseHsdSet2PbSet(inHsds []*HistoryDetail) []*pb.HistoryDetail {
	var pbHsds []*pb.HistoryDetail
	for _, inHsd := range inHsds {
		pbHsd := HistoryDetailToPb(inHsd)
		pbHsds = append(pbHsds, pbHsd)
	}
	return pbHsds
}
