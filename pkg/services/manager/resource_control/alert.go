package resource_control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"kubesphere.io/alert/pkg/constants"
	aldb "kubesphere.io/alert/pkg/db"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/metric"
	"kubesphere.io/alert/pkg/models"
	"kubesphere.io/alert/pkg/pb"
	"kubesphere.io/alert/pkg/util/pbutil"
	"kubesphere.io/alert/pkg/util/stringutil"
)

func RegisterAlert(ctx context.Context, alert *models.Alert) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&alert).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert alert failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func UpdateAlertInfo(ctx context.Context, alertId string, runningStatus string) error {
	attributes := make(map[string]interface{})

	attributes["running_status"] = runningStatus
	attributes["update_time"] = time.Now()

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var alert models.Alert
	err := tx.Model(&alert).Where(models.AlColId+" = ?", alertId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update alert failed, [%+v]\n", err.Error)
		return err.Error
	}

	tx.Commit()
	return nil
}

func CreateAlert(ctx context.Context, alert *models.Alert) error {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&alert).Error
	if err != nil {
		tx.Rollback()
		logger.Error(ctx, "Insert Alert failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func DescribeAlerts(ctx context.Context, req *pb.DescribeAlertsRequest) ([]*models.Alert, uint64, error) {
	req.AlertId = stringutil.SimplifyStringList(req.AlertId)
	req.AlertName = stringutil.SimplifyStringList(req.AlertName)
	req.RunningStatus = stringutil.SimplifyStringList(req.RunningStatus)
	req.PolicyId = stringutil.SimplifyStringList(req.PolicyId)
	req.RsFilterId = stringutil.SimplifyStringList(req.RsFilterId)
	req.ExecutorId = stringutil.SimplifyStringList(req.ExecutorId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var als []*models.Alert
	var count uint64

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableAlert)).
		AddQueryOrderDir(req, models.AlColCreateTime).
		BuildFilterConditions(req, models.TableAlert).
		Offset(offset).
		Limit(limit).
		Find(&als).Error; err != nil {
		logger.Error(ctx, "Describe Alerts failed: %+v", err)
		return nil, 0, err
	}

	if err := aldb.GetChain(global.GetInstance().GetDB().Table(models.TableAlert)).
		BuildFilterConditions(req, models.TableAlert).
		Count(&count).Error; err != nil {
		logger.Error(ctx, "Describe Alerts count failed: %+v", err)
		return nil, 0, err
	}

	return als, count, nil
}

func UpdateAlert(ctx context.Context, req *pb.ModifyAlertRequest, alertId string, operation string) (string, error) {
	attributes := make(map[string]interface{})

	switch operation {
	case "updating":
		attributes[models.AlColId] = alertId
		if req.AlertName != "" {
			attributes[models.AlColName] = req.AlertName
		}
		attributes[models.AlColDisabled] = req.Disabled
		if req.PolicyId != "" {
			attributes[models.AlColPolicyId] = req.PolicyId
		}
		if req.RsFilterId != "" {
			attributes[models.AlColRsFilterId] = req.RsFilterId
		}

		attributes[models.AlColRunningStatus] = operation
		attributes[models.AlColUpdateTime] = time.Now()
	case "deleting":
		attributes[models.AlColId] = alertId

		attributes[models.AlColRunningStatus] = operation
		attributes[models.AlColUpdateTime] = time.Now()
	}

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var alert models.Alert
	err := tx.Model(&alert).Where(models.AlColId+" = ?", alertId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update Alert [%s] failed: %+v", alertId, err.Error)
		return alertId, err.Error
	}
	tx.Commit()

	return alertId, nil
}

func GetAlertByName(resourceMap map[string]string, alertName string) ([]*models.Alert, uint64, error) {
	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("alert t1").
		Select("t1.alert_id,t1.alert_name").
		Joins("left join resource_filter t2 on t1.rs_filter_id=t2.rs_filter_id").
		Joins("left join resource_type t3 on t2.rs_type_id=t3.rs_type_id"))

	var als []*models.Alert
	var count uint64

	dbChain, orderByStr := buildDB4GetAlertByName(dbChain, resourceMap, alertName)

	err := dbChain.
		Order(orderByStr).
		Scan(&als).
		Error
	if err != nil {
		logger.Error(nil, "Failed to GetAlertByName [%v], error: %+v.", alertName, err)
		return nil, 0, err
	}

	err = dbChain.Count(&count).Error
	if err != nil {
		logger.Error(nil, "Failed to GetAlertByName [%v], error: %+v.", alertName, err)
		return nil, 0, err
	}

	return als, count, nil
}

func buildDB4GetAlertByName(dbChain *aldb.Chain, resourceMap map[string]string, alertName string) (*aldb.Chain, string) {
	alertNames := stringutil.SimplifyStringList(strings.Split(alertName, ","))

	dbChain.DB = dbChain.DB.Where(`t3.rs_type_name in (?)`, resourceMap["rs_type_name"])
	switch resourceMap["rs_type_name"] {
	case "cluster":
		break
	case "node":
		break
	case "workspace":
		if resourceMap["ws_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t2.rs_filter_param, "$.ws_name") in (?)`, resourceMap["ws_name"])
		}
	case "namespace":
		if resourceMap["ns_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t2.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
		}
	case "workload":
		if resourceMap["ns_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t2.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
		}
	case "pod":
		if resourceMap["ns_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t2.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
		}
		if resourceMap["node_id"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t2.rs_filter_param, "$.node_id") in (?)`, resourceMap["node_id"])
		}
	case "container":
		if resourceMap["ns_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t2.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
		}
		if resourceMap["node_id"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t2.rs_filter_param, "$.node_id") in (?)`, resourceMap["node_id"])
		}
		if resourceMap["pod_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t2.rs_filter_param, "$.pod_name") in (?)`, resourceMap["pod_name"])
		}
	}
	if len(alertNames) != 0 {
		dbChain.DB = dbChain.DB.Where(`t1.alert_name in (?)`, alertNames)
	}
	var sortKeyStr string = "t1.alert_id"
	var reverseStr string = constants.DESC
	orderByStr := sortKeyStr + " " + reverseStr

	return dbChain, orderByStr
}

func UpdateAlertByName(ctx context.Context, resourceMap map[string]string, alertName string, operation string, req *pb.ModifyAlertByNameRequest) (string, string, error) {
	alerts, count, _ := GetAlertByName(resourceMap, alertName)

	if count != 1 {
		return "", "", errors.New("alert_name not found or duplicate")
	}

	alertId := alerts[0].AlertId

	attributes := make(map[string]interface{})

	switch operation {
	case "updating":
		attributes[models.AlColDisabled] = req.Disabled
		if req.PolicyId != "" {
			attributes[models.AlColPolicyId] = req.PolicyId
		}
		if req.RsFilterId != "" {
			attributes[models.AlColRsFilterId] = req.RsFilterId
		}

		attributes[models.AlColRunningStatus] = operation
		attributes[models.AlColUpdateTime] = time.Now()
	case "deleting":
		attributes[models.AlColRunningStatus] = operation
		attributes[models.AlColUpdateTime] = time.Now()
	}

	db := global.GetInstance().GetDB()
	tx := db.Begin()

	var alert models.Alert
	err := tx.Model(&alert).Where(models.AlColId+" = ?", alertId).Updates(attributes)
	if err.Error != nil {
		tx.Rollback()
		logger.Error(ctx, "Update Alert By Name[%s] failed: %+v", alertName, err.Error)
		return "", "", err.Error
	}

	tx.Commit()
	return alertId, alertName, nil
}

type MetricInfo struct {
	MetricName string `gorm:"column:metric_name" json:"metric_name"`
}

func getMetricsByAlertId(alertId string) []string {
	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("alert t1").
		Select("t3.metric_name").
		Joins("left join rule t2 on t2.policy_id=t1.policy_id").
		Joins("left join metric t3 on t3.metric_id=t2.metric_id"))

	dbChain.DB = dbChain.DB.Where("t1.alert_id in (?)", alertId)

	var mts []*MetricInfo

	err := dbChain.
		Offset(0).
		Limit(100).
		Order("t3.metric_name asc").
		Scan(&mts).
		Error
	if err != nil {
		logger.Error(nil, "Failed to getMetricsByAlertId [%v], error: %+v.", alertId, err)
		return nil
	}

	metrics := []string{}

	for _, mt := range mts {
		metrics = append(metrics, mt.MetricName)
	}

	return metrics
}

type StatusAlert struct {
	ResourceStatus map[string]StatusResource `json:resource_status`
	UpdateTime     time.Time
}

type StatusResource struct {
	CurrentLevel       string          `json:current_level`
	PositiveCount      uint32          `json:positive_count`
	CumulatedSendCount uint32          `json:cumulated_send_count`
	NextResendInterval uint32          `json:next_resend_interval`
	NextSendableTime   time.Time       `json:next_sendable_time`
	AggregatedAlerts   AggregatedAlert `json:aggregated_alerts`
}

type AggregatedAlert struct {
	CumulatedCount  uint32           `json:"cumulated_count"`
	FirstAlertTime  string           `json:"first_alert_time"`
	LastAlertTime   string           `json:"last_alert_time"`
	LastAlertValues []RecordedMetric `json:"last_alert_values"`
}

type RecordedMetric struct {
	RuleName     string
	ResourceName string
	tvs          []metric.TV
}

type MessageInfo struct {
	CreateAt string `gorm:"column:create_time" json:"create_time"`
}

func getMostRecentAlertTimeByAlertId(alertId string) string {
	//TODO: Optimize
	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("alert t1").
		Select("t2.create_time").
		Joins("left join history t2 on t2.alert_id=t1.alert_id"))

	dbChain.DB = dbChain.DB.Where(`t1.alert_id in (?) and t2.event in ("triggered", "sent_success", "sent_failed")`, alertId)

	var mis []*MessageInfo

	err := dbChain.
		Offset(0).
		Limit(1).
		Order("t2.create_time desc").
		Scan(&mis).
		Error
	if err != nil {
		logger.Error(nil, "Failed to getMostRecentAlertTimeByAlertId [%v], error: %+v.", alertId, err)
		return ""
	}

	mostRecentAlertTime := ""

	for _, mi := range mis {
		mostRecentAlertTime = mi.CreateAt
	}

	return mostRecentAlertTime
}

func DescribeAlertDetails(ctx context.Context, req *pb.DescribeAlertDetailsRequest) ([]*models.AlertDetail, uint64, error) {
	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("alert t1").
		Select("t1.alert_id,t1.alert_name,t1.disabled,t1.create_time,t1.running_status,t1.alert_status,t1.policy_id,t3.rs_filter_name,t3.rs_filter_param,t4.rs_type_name,t1.executor_id,t2.policy_name,t2.policy_description,t2.policy_config,t2.creator,t2.available_start_time,t2.available_end_time,t5.nf_address_list_id").
		Joins("left join policy t2 on t1.policy_id=t2.policy_id").
		Joins("left join resource_filter t3 on t1.rs_filter_id=t3.rs_filter_id").
		Joins("left join resource_type t4 on t3.rs_type_id=t4.rs_type_id").
		Joins("left join action t5 on t5.policy_id=t1.policy_id"))

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var alds []*models.AlertDetail
	var count uint64

	dbChain, orderByStr := buildDB4DescAlertDetails(dbChain, req)

	err := dbChain.
		Offset(offset).
		Limit(limit).
		Order(orderByStr).
		Scan(&alds).
		Error
	if err != nil {
		logger.Error(nil, "Failed to Describe Alert Details [%v], error: %+v.", req, err)
		return nil, 0, err
	}

	err = dbChain.Count(&count).Error
	if err != nil {
		logger.Error(nil, "Failed to Describe Alert Details [%v], error: %+v.", req, err)
		return nil, 0, err
	}

	for _, ald := range alds {
		ald.Metrics = getMetricsByAlertId(ald.AlertId)
		alertStatus := StatusAlert{}
		err := json.Unmarshal([]byte(ald.AlertStatus), &alertStatus)
		if err == nil {
			ald.RulesCount = uint32(0)
			ald.PositivesCount = uint32(0)
			for _, v := range alertStatus.ResourceStatus {
				ald.RulesCount = ald.RulesCount + 1
				if v.CurrentLevel != "cleared" {
					ald.PositivesCount = ald.PositivesCount + 1
				}
			}
		}
		if ald.RulesCount > 0 {
			ald.MostRecentAlertTime = getMostRecentAlertTimeByAlertId(ald.AlertId)
		}
	}

	return alds, count, nil
}

func buildDB4DescAlertDetails(dbChain *aldb.Chain, req *pb.DescribeAlertDetailsRequest) (*aldb.Chain, string) {
	resourceSearch := req.ResourceSearch
	alertId := stringutil.SimplifyStringList(req.AlertId)
	alertName := stringutil.SimplifyStringList(req.AlertName)
	runningStatus := stringutil.SimplifyStringList(req.RunningStatus)
	policyId := stringutil.SimplifyStringList(req.PolicyId)
	creator := stringutil.SimplifyStringList(req.Creator)
	rsFilterId := stringutil.SimplifyStringList(req.RsFilterId)
	executorId := stringutil.SimplifyStringList(req.ExecutorId)

	if resourceSearch != "" {
		resourceMap := map[string]string{}
		err := json.Unmarshal([]byte(resourceSearch), &resourceMap)
		if err == nil {
			dbChain.DB = dbChain.DB.Where(`t4.rs_type_name in (?)`, resourceMap["rs_type_name"])
			switch resourceMap["rs_type_name"] {
			case "cluster":
				break
			case "node":
				break
			case "workspace":
				if resourceMap["ws_name"] != "" {
					dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t3.rs_filter_param, "$.ws_name") in (?)`, resourceMap["ws_name"])
				}
			case "namespace":
				if resourceMap["ns_name"] != "" {
					dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t3.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
				}
			case "workload":
				if resourceMap["ns_name"] != "" {
					dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t3.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
				}
			case "pod":
				if resourceMap["ns_name"] != "" {
					dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t3.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
				}
				if resourceMap["node_id"] != "" {
					dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t3.rs_filter_param, "$.node_id") in (?)`, resourceMap["node_id"])
				}
			case "container":
				if resourceMap["ns_name"] != "" {
					dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t3.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
				}
				if resourceMap["node_id"] != "" {
					dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t3.rs_filter_param, "$.node_id") in (?)`, resourceMap["node_id"])
				}
				if resourceMap["pod_name"] != "" {
					dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t3.rs_filter_param, "$.pod_name") in (?)`, resourceMap["pod_name"])
				}
			}
		} else {
			logger.Error(nil, "Unmarshal resourceSearch error %v", err)
		}
	}
	if len(alertId) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.alert_id in (?)", alertId)
	}
	if len(alertName) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.alert_name in (?)", alertName)
	}
	/*if len(req.Disabled) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.disabled in (?)", req.Disabled)
	}*/
	if len(runningStatus) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.running_status in (?)", runningStatus)
	}
	if len(policyId) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.policy_id in (?)", policyId)
	}
	if len(creator) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.creator in (?)", creator)
	}
	if len(rsFilterId) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.rs_filter_id in (?)", rsFilterId)
	}
	if len(executorId) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.executor_id in (?)", executorId)
	}
	//Step2：get SearchWord
	if req.SearchWord != "" {
		dbChain.DB = dbChain.DB.Where("alert_name LIKE ?", "%"+req.SearchWord+"%")
	}
	//Step3：get OrderByStr
	var sortKeyStr string = "t1.alert_id"
	var reverseStr string = constants.DESC
	orderByStr := sortKeyStr + " " + reverseStr

	if req.SortKey != "" {
		sortKeyStr = req.SortKey
		if req.Reverse {
			reverseStr = constants.DESC
		} else {
			reverseStr = constants.ASC
		}
		orderByStr = sortKeyStr + " " + reverseStr
	}

	return dbChain, orderByStr
}

func DescribeAlertStatus(ctx context.Context, req *pb.DescribeAlertStatusRequest) ([]models.AlertStatus, uint64, error) {
	resourceMap := map[string]string{}
	err := json.Unmarshal([]byte(req.ResourceSearch), &resourceMap)
	if err != nil {
		logger.Error(nil, "Failed to Describe Alert Status [%v], error: %+v.", req, err)
		return nil, 0, err
	}

	alertId := stringutil.SimplifyStringList(req.AlertId)
	alertName := stringutil.SimplifyStringList(req.AlertName)
	runningStatus := stringutil.SimplifyStringList(req.RunningStatus)
	policyId := stringutil.SimplifyStringList(req.PolicyId)
	creator := stringutil.SimplifyStringList(req.Creator)
	rsFilterId := stringutil.SimplifyStringList(req.RsFilterId)
	executorId := stringutil.SimplifyStringList(req.ExecutorId)
	ruleId := stringutil.SimplifyStringList(req.RuleId)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("alert t1").
		Select("t2.rule_id,t2.rule_name,t2.disabled,t2.monitor_periods,t2.severity,t2.metrics_type,t2.condition_type,t2.thresholds,t2.unit,t2.consecutive_count,t2.inhibit,t3.metric_name,t2.create_time,t2.update_time,t1.alert_status").
		Joins("left join rule t2 on t2.policy_id=t1.policy_id").
		Joins("left join metric t3 on t3.metric_id=t2.metric_id").
		Joins("left join resource_filter t4 on t4.rs_filter_id=t1.rs_filter_id").
		Joins("left join resource_type t5 on t5.rs_type_id=t3.rs_type_id"))

	dbChain.DB = dbChain.DB.Where(`t5.rs_type_name in (?)`, resourceMap["rs_type_name"])
	switch resourceMap["rs_type_name"] {
	case "cluster":
		break
	case "node":
		break
	case "workspace":
		if resourceMap["ws_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t4.rs_filter_param, "$.ws_name") in (?)`, resourceMap["ws_name"])
		}
	case "namespace":
		if resourceMap["ns_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
		}
	case "workload":
		if resourceMap["ns_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
		}
	case "pod":
		if resourceMap["ns_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
		}
		if resourceMap["node_id"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t4.rs_filter_param, "$.node_id") in (?)`, resourceMap["node_id"])
		}
	case "container":
		if resourceMap["ns_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in (?)`, resourceMap["ns_name"])
		}
		if resourceMap["node_id"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t4.rs_filter_param, "$.node_id") in (?)`, resourceMap["node_id"])
		}
		if resourceMap["pod_name"] != "" {
			dbChain.DB = dbChain.DB.Where(`JSON_EXTRACT(t4.rs_filter_param, "$.pod_name") in (?)`, resourceMap["pod_name"])
		}
	}

	if len(alertId) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.alert_id in (?)", alertId)
	}
	if len(alertName) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.alert_name in (?)", alertName)
	}
	if len(req.Disabled) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.disabled in (?)", req.Disabled)
	}
	if len(runningStatus) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.running_status in (?)", runningStatus)
	}
	if len(policyId) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.policy_id in (?)", policyId)
	}
	if len(creator) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.creator in (?)", creator)
	}
	if len(rsFilterId) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.rs_filter_id in (?)", rsFilterId)
	}
	if len(executorId) != 0 {
		dbChain.DB = dbChain.DB.Where("t1.executor_id in (?)", executorId)
	}
	if len(ruleId) != 0 {
		dbChain.DB = dbChain.DB.Where("t2.rule_id in (?)", ruleId)
	}
	//Step2：get SearchWord
	if req.SearchWord != "" {
		dbChain.DB = dbChain.DB.Where("alert_name LIKE ?", "%"+req.SearchWord+"%")
	}
	//Step3：get OrderByStr
	var sortKeyStr string = "t2.create_time"
	var reverseStr string = constants.DESC
	orderByStr := sortKeyStr + " " + reverseStr

	if req.SortKey != "" {
		sortKeyStr = req.SortKey
		if req.Reverse {
			reverseStr = constants.DESC
		} else {
			reverseStr = constants.ASC
		}
		orderByStr = sortKeyStr + " " + reverseStr
	}

	var alss []*models.AlertStatus
	var count uint64

	err = dbChain.
		Offset(offset).
		Limit(limit).
		Order(orderByStr).
		Scan(&alss).
		Error
	if err != nil {
		logger.Error(nil, "Failed to Describe Alert Status [%v], error: %+v.", req, err)
		return nil, 0, err
	}

	err = dbChain.Count(&count).Error
	if err != nil {
		logger.Error(nil, "Failed to Describe Alert Status [%v], error: %+v.", req, err)
		return nil, 0, err
	}

	als_resources := []models.AlertStatus{}
	count_resources := 0

	for _, als := range alss {
		alertStatus := StatusAlert{}
		err := json.Unmarshal([]byte(als.AlertStatus), &alertStatus)
		if err == nil {
			als_resource := *als
			for k, v := range alertStatus.ResourceStatus {
				alertid_resourcename := strings.Split(k, " ")
				if alertid_resourcename[0] == als.RuleId {
					resourceStatus := models.ResourceStatus{}
					resourceStatus.ResourceName = alertid_resourcename[1]
					resourceStatus.CurrentLevel = v.CurrentLevel
					resourceStatus.PositiveCount = v.PositiveCount
					resourceStatus.CumulatedSendCount = v.CumulatedSendCount
					resourceStatus.NextResendInterval = v.NextResendInterval
					resourceStatus.NextSendableTime = v.NextSendableTime.Format("2006-01-02 15:04:05.99999")
					resourceStatus.AggregatedAlerts = fmt.Sprintf("%v", v.AggregatedAlerts)
					als_resource.Resources = append(als_resource.Resources, resourceStatus)
				}
			}
			als_resources = append(als_resources, als_resource)
			count_resources = count_resources + 1
		} else {
			als_resource := *als
			als_resources = append(als_resources, als_resource)
			count_resources = count_resources + 1
		}
	}

	return als_resources, uint64(count_resources), nil
}
