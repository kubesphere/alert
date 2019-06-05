package resource_control

import (
	"encoding/json"
	"time"

	"kubesphere.io/alert/pkg/constants"
	aldb "kubesphere.io/alert/pkg/db"
	"kubesphere.io/alert/pkg/gerr"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/metric"
	"kubesphere.io/alert/pkg/util/stringutil"
)

func getOffset(offset uint32) uint32 {
	if offset == 0 {
		return constants.DefaultOffset
	}

	return offset
}

func getLimit(limit uint32) uint32 {
	if limit == 0 {
		return constants.DefaultLimit
	}
	return limit
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

type AlertDetail struct {
	AlertId             string    `gorm:"column:alert_id" json:"alert_id"`
	AlertName           string    `gorm:"column:alert_name" json:"alert_name"`
	Disabled            bool      `gorm:"column:disabled" json:"disabled"`
	CreateTime          time.Time `gorm:"column:create_time" json:"create_time"`
	RunningStatus       string    `gorm:"column:running_status" json:"running_status"`
	AlertStatus         string    `gorm:"column:alert_status" json:"alert_status"`
	PolicyId            string    `gorm:"column:policy_id" json:"policy_id"`
	RsFilterName        string    `gorm:"column:rs_filter_name" json:"rs_filter_name"`
	RsFilterParam       string    `gorm:"column:rs_filter_param" json:"rs_filter_param"`
	RsTypeName          string    `gorm:"column:rs_type_name" json:"rs_type_name"`
	ExecutorId          string    `gorm:"column:executor_id" json:"executor_id"`
	PolicyName          string    `gorm:"column:policy_name" json:"policy_name"`
	PolicyDescription   string    `gorm:"column:policy_description" json:"policy_description"`
	PolicyConfig        string    `gorm:"column:policy_config" json:"policy_config"`
	Creator             string    `gorm:"column:creator" json:"creator"`
	AvailableStartTime  string    `gorm:"column:available_start_time" json:"available_start_time"`
	AvailableEndTime    string    `gorm:"column:available_end_time" json:"available_end_time"`
	Metrics             []string  `gorm:"column:metrics" json:"metrics"`
	RulesCount          uint32    `json:"rules_count"`
	PositivesCount      uint32    `json:"positives_count"`
	MostRecentAlertTime string    `json:"most_recent_alert_time"`
	NfAddressListId     string    `gorm:"column:nf_address_list_id" json:"nf_address_list_id"`
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

type DescribeAlertDetailsRequest struct {
	SearchWord     string
	SortKey        string
	Reverse        bool
	Offset         uint32
	Limit          uint32
	ResourceSearch string
	AlertId        []string
	AlertName      []string
	Disabled       []bool
	RunningStatus  []string
	PolicyId       []string
	Creator        []string
	RsFilterId     []string
	ExecutorId     []string
}

type DescribeAlertDetailsResponse struct {
	Total          uint32         `json:"total"`
	AlertdetailSet []*AlertDetail `json:"alertdetail_set"`
}

func getAlertDetails(req *DescribeAlertDetailsRequest) ([]AlertDetail, uint64, error) {
	dbChain := aldb.GetChain(global.GetInstance().GetDB().Table("alert t1").
		Select("t1.alert_id,t1.alert_name,t1.disabled,t1.create_time,t1.running_status,t1.alert_status,t1.policy_id,t3.rs_filter_name,t3.rs_filter_param,t4.rs_type_name,t1.executor_id,t2.policy_name,t2.policy_description,t2.policy_config,t2.creator,t2.available_start_time,t2.available_end_time,t5.nf_address_list_id").
		Joins("left join policy t2 on t1.policy_id=t2.policy_id").
		Joins("left join resource_filter t3 on t1.rs_filter_id=t3.rs_filter_id").
		Joins("left join resource_type t4 on t3.rs_type_id=t4.rs_type_id").
		Joins("left join action t5 on t5.policy_id=t1.policy_id"))

	offset := getOffset(req.Offset)
	limit := getLimit(req.Limit)

	var alds []AlertDetail
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

func buildDB4DescAlertDetails(dbChain *aldb.Chain, req *DescribeAlertDetailsRequest) (*aldb.Chain, string) {
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

func parseAldSet(inAlds []AlertDetail) []*AlertDetail {
	var Alds []*AlertDetail
	for _, inAld := range inAlds {
		Alds = append(Alds, &inAld)
	}
	return Alds
}

func DescribeAlertDetails(req *DescribeAlertDetailsRequest) (*DescribeAlertDetailsResponse, error) {
	alds, aldCnt, err := getAlertDetails(req)
	if err != nil {
		logger.Error(nil, "Failed to Describe Alert Details, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(nil, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	aldSet := parseAldSet(alds)
	res := &DescribeAlertDetailsResponse{
		Total:          uint32(aldCnt),
		AlertdetailSet: aldSet,
	}

	logger.Debug(nil, "Describe Alert Details successfully, Alert Details=[%+v].", res)
	return res, nil
}
