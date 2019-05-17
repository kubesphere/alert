package resource_control

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"kubesphere.io/alert/pkg/constants"
	aldb "kubesphere.io/alert/pkg/db"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	"kubesphere.io/alert/pkg/pb"
	"kubesphere.io/alert/pkg/util/pbutil"
	"kubesphere.io/alert/pkg/util/stringutil"

	nf "kubesphere.io/alert/pkg/client/notification"
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

type EmailInfo struct {
	Address string `json:"address"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func parseInValues(values []string) string {
	output := ""

	for _, value := range values {
		output = output + fmt.Sprintf(`,'%s'`, value)
	}

	output = strings.TrimPrefix(output, ",")

	return output
}

func DescribeHistoryDetail(ctx context.Context, resourceMap map[string]string, req *pb.DescribeHistoryDetailRequest) ([]*models.HistoryDetail, uint64, error) {
	historyId := stringutil.SimplifyStringList(req.HistoryId)
	historyName := stringutil.SimplifyStringList(req.HistoryName)
	alertNames := stringutil.SimplifyStringList(req.AlertName)
	ruleName := stringutil.SimplifyStringList(req.RuleName)
	event := stringutil.SimplifyStringList(req.Event)
	ruleId := stringutil.SimplifyStringList(req.RuleId)
	resourceName := stringutil.SimplifyStringList(req.ResourceName)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	alertIds := []string{}
	for _, alertName := range alertNames {
		alerts, count, _ := GetAlertByName(resourceMap, alertName)

		if count == 1 {
			alertIds = append(alertIds, alerts[0].AlertId)
		}
	}

	if len(alertIds) == 0 && len(alertNames) != 0 {
		logger.Error(nil, "Describe History Detail has no match alert_name[%v]", req)
		return nil, 0, errors.New("alert_name not found")
	}

	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("history t1").
		Select("t1.history_id,t1.history_name,t1.event,t1.notification_id,t1.rule_id,t1.resource_name,t2.rule_name,t2.severity,t6.rs_type_name,t4.rs_filter_name,t5.metric_name,t2.condition_type,t2.thresholds,t2.unit,t3.alert_name,t4.rs_filter_param,t1.create_time,t1.update_time").
		Joins("left join rule t2 on t2.rule_id=t1.rule_id").
		Joins("left join alert t3 on t3.alert_id=t1.alert_id").
		Joins("left join resource_filter t4 on t4.rs_filter_id=t3.rs_filter_id").
		Joins("left join metric t5 on t5.metric_id=t2.metric_id").
		Joins("left join resource_type t6 on t6.rs_type_id=t4.rs_type_id"))

	whereResult := ""
	whereTriggered := ""

	whereResult = whereResult + fmt.Sprintf(`t6.rs_type_name in ("%s") and `, resourceMap["rs_type_name"])
	whereTriggered = whereTriggered + fmt.Sprintf(`t6.rs_type_name in ("%s") and `, resourceMap["rs_type_name"])

	switch resourceMap["rs_type_name"] {
	case "cluster":
		break
	case "node":
		break
	case "workspace":
		if resourceMap["ws_name"] != "" {
			whereResult = whereResult + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ws_name") in ("%s") and `, resourceMap["ws_name"])
			whereTriggered = whereTriggered + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ws_name") in ("%s") and `, resourceMap["ws_name"])
		}
	case "namespace":
		if resourceMap["ns_name"] != "" {
			whereResult = whereResult + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in ("%s") and `, resourceMap["ns_name"])
			whereTriggered = whereTriggered + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in ("%s") and `, resourceMap["ns_name"])
		}
	case "workload":
		if resourceMap["ns_name"] != "" {
			whereResult = whereResult + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in ("%s") and `, resourceMap["ns_name"])
			whereTriggered = whereTriggered + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in ("%s") and `, resourceMap["ns_name"])
		}
	case "pod":
		if resourceMap["ns_name"] != "" {
			whereResult = whereResult + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in ("%s") and `, resourceMap["ns_name"])
			whereTriggered = whereTriggered + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in ("%s") and `, resourceMap["ns_name"])
		}
		if resourceMap["node_id"] != "" {
			whereResult = whereResult + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.node_id") in ("%s") and `, resourceMap["node_id"])
			whereTriggered = whereTriggered + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.node_id") in ("%s") and `, resourceMap["node_id"])
		}
	case "container":
		if resourceMap["ns_name"] != "" {
			whereResult = whereResult + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in ("%s") and `, resourceMap["ns_name"])
			whereTriggered = whereTriggered + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.ns_name") in ("%s") and `, resourceMap["ns_name"])
		}
		if resourceMap["node_id"] != "" {
			whereResult = whereResult + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.node_id") in ("%s") and `, resourceMap["node_id"])
			whereTriggered = whereTriggered + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.node_id") in ("%s") and `, resourceMap["node_id"])
		}
		if resourceMap["pod_name"] != "" {
			whereResult = whereResult + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.pod_name") in ("%s") and `, resourceMap["pod_name"])
			whereTriggered = whereTriggered + fmt.Sprintf(`JSON_EXTRACT(t4.rs_filter_param, "$.pod_name") in ("%s") and `, resourceMap["pod_name"])
		}
	}

	if len(historyId) != 0 {
		whereResult = whereResult + fmt.Sprintf("t1.history_id in (%s) and ", parseInValues(historyId))
		whereTriggered = whereTriggered + fmt.Sprintf("t1.history_id in (%s) and ", parseInValues(historyId))
	}
	if len(historyName) != 0 {
		whereResult = whereResult + fmt.Sprintf("t1.history_name in (%s) and ", parseInValues(historyName))
		whereTriggered = whereTriggered + fmt.Sprintf("t1.history_name in (%s) and ", parseInValues(historyName))
	}
	if len(alertIds) != 0 {
		whereResult = whereResult + fmt.Sprintf("t3.alert_id in (%s) and ", parseInValues(alertIds))
		whereTriggered = whereTriggered + fmt.Sprintf("t3.alert_id in (%s) and ", parseInValues(alertIds))
	}
	if len(ruleName) != 0 {
		whereResult = whereResult + fmt.Sprintf("t2.rule_name in (%s) and ", parseInValues(ruleName))
		whereTriggered = whereTriggered + fmt.Sprintf("t2.rule_name in (%s) and ", parseInValues(ruleName))
	}
	if len(event) != 0 {
		whereResult = whereResult + fmt.Sprintf("t1.event in (%s) and ", parseInValues(event))
	}
	whereTriggered = whereTriggered + fmt.Sprintf("t1.event = 'triggered' and ")
	if len(ruleId) != 0 {
		whereResult = whereResult + fmt.Sprintf("t1.rule_id in (%s) and ", parseInValues(ruleId))
		whereTriggered = whereTriggered + fmt.Sprintf("t1.rule_id in (%s) and ", parseInValues(ruleId))
	}
	if len(resourceName) != 0 {
		whereResult = whereResult + fmt.Sprintf("t1.resource_name in (%s) and ", parseInValues(resourceName))
		whereTriggered = whereTriggered + fmt.Sprintf("t1.resource_name in (%s) and ", parseInValues(resourceName))
	}
	//Step2: get SearchWord
	if req.SearchWord != "" {
		whereResult = whereResult + fmt.Sprintf("t1.resource_name LIKE '%%%s%%' and ", req.SearchWord)
		whereTriggered = whereTriggered + fmt.Sprintf("t1.resource_name LIKE '%%%s%%' and ", req.SearchWord)
	}

	whereResult = strings.TrimSuffix(whereResult, " and ")
	whereTriggered = strings.TrimSuffix(whereTriggered, " and ")

	//Step3: get Recent
	if req.Recent {
		dbChain.DB = dbChain.DB.Where(fmt.Sprintf(`%s and t1.create_time >= (select max(t1.create_time) from history t1 left join rule t2 on t2.rule_id=t1.rule_id left join alert t3 on t3.alert_id=t1.alert_id left join resource_filter t4 on t4.rs_filter_id=t3.rs_filter_id left join metric t5 on t5.metric_id=t2.metric_id left join resource_type t6 on t6.rs_type_id=t4.rs_type_id where %s)`, whereResult, whereTriggered))
	} else {
		dbChain.DB = dbChain.DB.Where(whereResult)
	}
	//Step4: get OrderByStr
	var sortKeyStr string = "t1.create_time"
	var reverseStr string = constants.ASC
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

	var hsds []*models.HistoryDetail
	var count uint64

	err := dbChain.
		Offset(offset).
		Limit(limit).
		Order(orderByStr).
		Scan(&hsds).
		Error
	if err != nil {
		logger.Error(nil, "Failed to Describe History Detail [%v], error: %+v.", req, err)
		return nil, 0, err
	}

	notificationIds := []string{}
	for _, hsd := range hsds {
		if hsd.NotificationId == "" {
			continue
		}

		notificationIds = append(notificationIds, hsd.NotificationId)
	}

	if len(notificationIds) > 0 {
		notificationStatusMap := nf.GetNotificationStatus(notificationIds)

		if notificationStatusMap != nil {
			for _, hsd := range hsds {
				if hsd.NotificationId == "" {
					continue
				}

				emailInfos := []EmailInfo{}

				for i := 0; i < len(notificationStatusMap[hsd.NotificationId]); i += 3 {
					address := ""
					directiveMap := make(map[string]interface{})
					err := json.Unmarshal([]byte(notificationStatusMap[hsd.NotificationId][i]), &directiveMap)
					if err == nil {
						address = directiveMap["Address"].(string)
					} else {
						logger.Error(nil, "Unmarshal directive error: %+v.", err)
					}
					emailInfos = append(emailInfos, EmailInfo{address, notificationStatusMap[hsd.NotificationId][i+1], notificationStatusMap[hsd.NotificationId][i+2]})
				}

				notificationStatus, _ := json.Marshal(emailInfos)
				hsd.NotificationStatus = string(notificationStatus)
			}
		}
	}

	err = dbChain.Count(&count).Error
	if err != nil {
		logger.Error(nil, "Failed to Describe Alert Rule Status [%v], error: %+v.", req, err)
		return nil, 0, err
	}

	return hsds, count, nil
}
