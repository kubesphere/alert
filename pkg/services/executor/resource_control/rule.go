package resource_control

import (
	aldb "kubesphere.io/alert/pkg/db"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
)

type RuleDetail struct {
	RuleId           string `gorm:"column:rule_id" json:"rule_id"`
	RuleName         string `gorm:"column:rule_name" json:"rule_name"`
	Disabled         bool   `gorm:"column:disabled" json:"disabled"`
	MonitorPeriods   uint32 `gorm:"column:monitor_periods" json:"monitor_periods"`
	Severity         string `gorm:"column:severity" json:"severity"`
	MetricsType      string `gorm:"column:metrics_type" json:"metrics_type"`
	ConditionType    string `gorm:"column:condition_type" json:"condition_type"`
	Thresholds       string `gorm:"column:thresholds" json:"thresholds"`
	Unit             string `gorm:"column:unit" json:"unit"`
	ConsecutiveCount uint32 `gorm:"column:consecutive_count" json:"consecutive_count"`
	Inhibit          bool   `gorm:"column:inhibit" json:"inhibit"`
	MetricName       string `gorm:"column:metric_name" json:"metric_name"`
	MetricParam      string `gorm:"column:metric_param" json:"metric_param"`
}

func QueryRuleDetails(alertId string) []RuleDetail {
	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("rule t1").
		Select("t1.rule_id,t1.rule_name,t1.disabled,t1.monitor_periods,t1.severity,t1.metrics_type,t1.condition_type,t1.thresholds,t1.unit,t1.consecutive_count,t1.inhibit,t1.policy_id,t2.metric_name,t2.metric_param").
		Joins("left join metric t2 on t2.metric_id=t1.metric_id"))

	dbChain.DB = dbChain.DB.Where("t1.policy_id in (select policy_id from alert where alert_id = ?)", alertId)

	var rds []RuleDetail

	err := dbChain.
		Scan(&rds).
		Error
	if err != nil {
		logger.Error(nil, "Failed to QueryRuleDetails [%v], error: %+v.", alertId, err)
		return nil
	}

	return rds
}
