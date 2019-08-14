package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"kubesphere.io/alert/pkg/client/adapter"
	nf "kubesphere.io/alert/pkg/client/notification"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/metric"
	"kubesphere.io/alert/pkg/models"
	"kubesphere.io/alert/pkg/notification"
	rs "kubesphere.io/alert/pkg/services/executor/resource_control"
)

type AlertRunner struct {
	AlertConfig ConfigAlert
	AlertStatus StatusAlert
	SignalCh    chan string
	UpdateCh    chan string
}

type ConfigAlert struct {
	AlertId            string
	LoadSuccess        bool
	Disabled           bool
	RsTypeName         string
	RsTypeParam        string
	RsFilterName       string
	RsFilterParam      string
	PolicyConfig       map[string]ConfigPolicy `json:"policy_config"`
	AvailableStartTime string
	AvailableEndTime   string
	Rules              map[string]RuleInfo
	Requests           MonitoringRequest
	NfAddressListId    string
}

type ConfigPolicy struct {
	RepeatType              string `json:"repeat_type"`
	RepeatIntervalInitvalue uint32 `json:"repeat_interval_initvalue"`
	MaxSendCount            uint32 `json:"max_send_count"`
}

type RuleInfo struct {
	RuleName         string
	Disabled         bool
	MonitorPeriods   uint32
	Severity         string
	MetricsType      string
	ConditionType    string
	Thresholds       float64
	Scale            float64
	Unit             string
	ConsecutiveCount uint32
	Inhibit          bool
	MetricName       string
}

type StatusAlert struct {
	sync.RWMutex
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

type MonitoringRequest struct {
	RulesSamePeriod map[uint32][]string
	TickCount       map[uint32]uint32
	TickUpperBound  map[uint32]uint32
}

type RecordedMetric struct {
	RuleName     string
	ResourceName string
	tvs          []metric.TV
}

const (
	TickPeriodSecond = 10
)

func NewAlertRunner(alertId string, updateCh chan string) *AlertRunner {
	runner := &AlertRunner{}

	runner.AlertConfig.AlertId = alertId
	runner.AlertStatus.UpdateTime = time.Now()
	runner.SignalCh = make(chan string, 10)
	runner.UpdateCh = updateCh

	return runner
}

func (ar *AlertRunner) parseNotification(nfAddressListId string) {
	ar.AlertConfig.NfAddressListId = nfAddressListId
}

func (ar *AlertRunner) parsePolicyConfig(alertDetail rs.AlertDetail) {
	err := json.Unmarshal([]byte(alertDetail.PolicyConfig), &ar.AlertConfig.PolicyConfig)
	if err != nil {
		//Fill a default config for policy
		ar.AlertConfig.PolicyConfig = make(map[string]ConfigPolicy)

		ar.AlertConfig.PolicyConfig["minor"] = ConfigPolicy{"not-repeat", 3, 3}
		ar.AlertConfig.PolicyConfig["major"] = ConfigPolicy{"exp-minutes", 2, 5}
		ar.AlertConfig.PolicyConfig["critical"] = ConfigPolicy{"fixed-minutes", 1, 8}
	}

	ar.AlertConfig.AvailableStartTime = alertDetail.AvailableStartTime
	ar.AlertConfig.AvailableEndTime = alertDetail.AvailableEndTime
}

func (ar *AlertRunner) parseRules() {
	ruleDetails := rs.QueryRuleDetails(ar.AlertConfig.AlertId)
	//logger.Debug(nil, "rules: %v", rules)

	mapRules := make(map[string]RuleInfo)

	for _, ruleDetail := range ruleDetails {
		threshold, _ := strconv.ParseFloat(ruleDetail.Thresholds, 64)
		scale, _ := strconv.ParseFloat(ruleDetail.MetricParam, 64)
		ruleInfo := RuleInfo{
			RuleName:         ruleDetail.RuleName,
			Disabled:         ruleDetail.Disabled,
			MonitorPeriods:   ruleDetail.MonitorPeriods,
			Severity:         ruleDetail.Severity,
			MetricsType:      ruleDetail.MetricsType,
			ConditionType:    ruleDetail.ConditionType,
			Thresholds:       threshold,
			Scale:            scale,
			Unit:             ruleDetail.Unit,
			ConsecutiveCount: ruleDetail.ConsecutiveCount,
			Inhibit:          ruleDetail.Inhibit,
		}

		ruleInfo.MetricName = ruleDetail.MetricName
		mapRules[ruleDetail.RuleId] = ruleInfo
	}
	ar.AlertConfig.Rules = mapRules

	//Put rules with same period into same MonitoringRequest group
	rulesSamePeriod := make(map[uint32][]string)
	tickCount := make(map[uint32]uint32)
	tickUpperBound := make(map[uint32]uint32)

	for ruleId, ruleInfo := range ar.AlertConfig.Rules {
		if ruleInfo.Disabled {
			continue
		}
		rulesSamePeriod[ruleInfo.MonitorPeriods] = append(rulesSamePeriod[ruleInfo.MonitorPeriods], ruleId)
		tickUpperBound[ruleInfo.MonitorPeriods] = ruleInfo.MonitorPeriods * 60 / TickPeriodSecond
		tickCount[ruleInfo.MonitorPeriods] = tickUpperBound[ruleInfo.MonitorPeriods] - 1
	}

	ar.AlertConfig.Requests = MonitoringRequest{rulesSamePeriod, tickCount, tickUpperBound}
}

func (ar *AlertRunner) getResetResourceStatus(ruleId string) StatusResource {
	ruleConfig := ar.AlertConfig.Rules[ruleId]

	resourceStatus := StatusResource{
		CurrentLevel:       "cleared",
		PositiveCount:      0,
		CumulatedSendCount: 0,
		NextResendInterval: ar.AlertConfig.PolicyConfig[ruleConfig.Severity].RepeatIntervalInitvalue,
		NextSendableTime:   time.Now(),
	}

	return resourceStatus
}

func (ar *AlertRunner) resetAlertStatus() {
	ar.AlertStatus.ResourceStatus = make(map[string]StatusResource)
}

func (ar *AlertRunner) parseAlertConfigStatus(alertDetail rs.AlertDetail) {
	ar.AlertConfig.Disabled = alertDetail.Disabled

	ar.AlertStatus.Lock()
	err := json.Unmarshal([]byte(alertDetail.AlertStatus), &ar.AlertStatus)
	if err != nil {
		logger.Debug(nil, "Parse Alert Status error: %v", err)
		ar.resetAlertStatus()
	}
	ar.AlertStatus.Unlock()
}

func (ar *AlertRunner) loadAlertInfo() {
	alertDetail, err := rs.QueryAlertDetail(ar.AlertConfig.AlertId)

	if err != nil {
		logger.Error(nil, "loadAlertInfo error: %v", err)
		ar.AlertConfig.LoadSuccess = false
		return
	}

	ar.AlertConfig.LoadSuccess = true

	//1. Parse Resource
	ar.AlertConfig.RsTypeName = alertDetail.RsTypeName
	ar.AlertConfig.RsTypeParam = alertDetail.RsTypeParam
	ar.AlertConfig.RsFilterName = alertDetail.RsFilterName
	ar.AlertConfig.RsFilterParam = alertDetail.RsFilterParam

	//2. Parse Notification
	ar.parseNotification(alertDetail.NfAddressListId)

	//3. Parse policy config
	ar.parsePolicyConfig(alertDetail)

	//4. Parse Rules config
	ar.parseRules()

	//5. Parse Alert status
	ar.parseAlertConfigStatus(alertDetail)

	logger.Debug(nil, "loadAlertInfo alert: %v", ar)
}

func (ar *AlertRunner) getOneMetric(period uint32, ch chan metric.ResourceMetrics) {
	extraQueryParams := ""

	metrics := []string{}
	metricToRule := make(map[string][]string)

	for _, ruleId := range ar.AlertConfig.Requests.RulesSamePeriod[period] {
		metricName := ar.AlertConfig.Rules[ruleId].MetricName
		metrics = append(metrics, metricName)
		metricToRule[metricName] = append(metricToRule[metricName], ruleId)
	}

	metricParam := metric.MetricParam{
		RsTypeName:       ar.AlertConfig.RsTypeName,
		RsTypeParam:      ar.AlertConfig.RsTypeParam,
		RsFilterName:     ar.AlertConfig.RsFilterName,
		RsFilterParam:    ar.AlertConfig.RsFilterParam,
		ExtraQueryParams: extraQueryParams,
		Metrics:          metrics,
		MetricToRule:     metricToRule,
	}

	metricParamBytes, err := json.Marshal(metricParam)
	if err != nil {
		logger.Error(nil, "Marshal Metric Param error: %v", err)
		return
	}

	resourceMetricsStr := adapter.SendMetricRequest(string(metricParamBytes))

	resourceMetrics := []metric.ResourceMetrics{}

	err = json.Unmarshal([]byte(resourceMetricsStr), &resourceMetrics)

	if err != nil {
		logger.Debug(nil, "Unmarshal Metric Result error: %v", err)
		return
	}

	for _, rm := range resourceMetrics {
		ch <- rm
	}
}

func (ar *AlertRunner) getResourceMetrics(ch chan metric.ResourceMetrics) {
	wg := sync.WaitGroup{}
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second*3)

	for period, _ := range ar.AlertConfig.Requests.RulesSamePeriod {
		ar.AlertConfig.Requests.TickCount[period] = (ar.AlertConfig.Requests.TickCount[period] + 1) % ar.AlertConfig.Requests.TickUpperBound[period]

		if ar.AlertConfig.Requests.TickCount[period] == 0 {
			wg.Add(1)
			go func(period uint32) {
				for {
					select {
					case <-ctx.Done():
						logger.Debug(nil, "getResourceMetrics canceled")
						wg.Done()
						return
					default:
						ar.getOneMetric(period, ch)
						wg.Done()
						return
					}
				}
			}(period)
		}
	}

	wg.Wait()
}

func (ar *AlertRunner) readRuleResourceMetric(resourceMetrics metric.ResourceMetrics, triggeredMetrics *[]RecordedMetric, resumedMetrics *[]RecordedMetric) string {
	rule := ar.AlertConfig.Rules[resourceMetrics.RuleId]
	condition := rule.ConditionType
	threshold := rule.Thresholds
	scale := rule.Scale

	for resourceName, timeValue := range resourceMetrics.ResourceMetric {
		if len(timeValue) < int(1) {
			continue
		}
		//Fetch last time value
		v, err := strconv.ParseFloat(timeValue[(len(timeValue)-int(1))].V, 64)
		if err != nil {
			logger.Error(nil, "readRuleResourceMetric error %v, value will be ignored!", err)
			continue
		}
		resourceSet := false
		switch condition {
		case ">=":
			resourceSet = (v*scale >= threshold)
		case ">":
			resourceSet = (v*scale > threshold)
		case "<=":
			resourceSet = (v*scale <= threshold)
		case "<":
			resourceSet = (v*scale < threshold)
		}

		if resourceSet {
			*triggeredMetrics = append(*triggeredMetrics, RecordedMetric{rule.RuleName, resourceName, timeValue})
		} else {
			*resumedMetrics = append(*resumedMetrics, RecordedMetric{rule.RuleName, resourceName, timeValue})
		}
	}

	return resourceMetrics.RuleId
}

func getRuleResourceKey(ruleId string, resourceName string) string {
	return ruleId + " " + resourceName
}

func (ar *AlertRunner) signalUpdate() {
	ar.UpdateCh <- "update"
}

func (ar *AlertRunner) checkOneMetric(resourceMetrics metric.ResourceMetrics) bool {
	triggeredMetrics := []RecordedMetric{}
	resumedMetrics := []RecordedMetric{}

	ruleId := ar.readRuleResourceMetric(resourceMetrics, &triggeredMetrics, &resumedMetrics)

	oldResourceStatus := ar.AlertStatus.ResourceStatus
	newResourceStatus := make(map[string]StatusResource)

	needUpdate := false

	for _, triggeredMetric := range triggeredMetrics {
		resourceName := triggeredMetric.ResourceName
		ruleResourceKey := getRuleResourceKey(ruleId, resourceName)
		newStatus := StatusResource{}
		if _, ok := oldResourceStatus[ruleResourceKey]; ok {
			newStatus = oldResourceStatus[ruleResourceKey]
		} else {
			newStatus = ar.getResetResourceStatus(ruleId)
		}

		operation := ""
		resourceIsAlert := false
		newStatus.PositiveCount = newStatus.PositiveCount + 1
		if newStatus.PositiveCount >= ar.AlertConfig.Rules[ruleId].ConsecutiveCount {
			resourceIsAlert = true
			if newStatus.CurrentLevel == "cleared" {
				newStatus.CurrentLevel = ar.AlertConfig.Rules[ruleId].Severity
				newStatus.NextSendableTime = time.Now()
				operation = "trigger"
			}
		}

		if operation == "trigger" {
			logger.Debug(nil, "Rule[%v] Resource[%v] %v triggered, write to message", ruleId, resourceName, triggeredMetric)
			ar.writeHistory("", "triggered", fmt.Sprintf("%v", triggeredMetric), "", ruleId, resourceName)
			needUpdate = true
		}

		if resourceIsAlert {
			ar.sendActiveNotification(&newStatus, ruleId, resourceName, triggeredMetrics)
		}

		newResourceStatus[ruleResourceKey] = newStatus
	}

	for _, resumedMetric := range resumedMetrics {
		resourceName := resumedMetric.ResourceName
		ruleResourceKey := getRuleResourceKey(ruleId, resourceName)
		newStatus := StatusResource{}
		if _, ok := oldResourceStatus[ruleResourceKey]; ok {
			newStatus = oldResourceStatus[ruleResourceKey]
		} else {
			newStatus = ar.getResetResourceStatus(ruleId)
		}

		operation := ""
		newStatus.PositiveCount = 0
		if newStatus.CurrentLevel != "cleared" {
			newStatus = ar.getResetResourceStatus(ruleId)
			operation = "resume"
		}

		if operation == "resume" {
			logger.Debug(nil, "Rule[%v] Resource[%v] %v resumed, write to message", ruleId, resourceName, resumedMetric)
			ar.writeHistory("", "resumed", fmt.Sprintf("%v", resumedMetric), "", ruleId, resourceName)

			resumeStatus := StatusResource{}
			if _, ok := oldResourceStatus[ruleResourceKey]; ok {
				resumeStatus = oldResourceStatus[ruleResourceKey]
			} else {
				resumeStatus = ar.getResetResourceStatus(ruleId)
			}
			ar.sendResumeNotification(&resumeStatus, ruleId, resourceName, resumedMetric, resumedMetrics)
			needUpdate = true
		}

		newResourceStatus[ruleResourceKey] = newStatus
	}

	ar.AlertStatus.Lock()
	for k, v := range oldResourceStatus {
		oldRuleId := strings.Split(k, " ")[0]
		if oldRuleId != ruleId {
			newResourceStatus[k] = v
		}
	}
	ar.AlertStatus.ResourceStatus = newResourceStatus
	ar.AlertStatus.Unlock()

	return needUpdate
}

func (ar *AlertRunner) checkMetrics(ch chan metric.ResourceMetrics) {
	needUpdate := false

	for resourceMetrics := range ch {
		logger.Debug(nil, "resourceMetrics %v", resourceMetrics)

		needUpdate = needUpdate || ar.checkOneMetric(resourceMetrics)
	}

	if needUpdate {
		ar.signalUpdate()
	}
}

func (ar *AlertRunner) writeHistory(historyName string, status string, content string, notificatioId string, ruleId string, resourceName string) {
	history := models.NewHistory(
		"",
		status,
		content,
		notificatioId,
		ar.AlertConfig.AlertId,
		ruleId,
		resourceName,
	)

	err := rs.CreateHistory(nil, history)
	if err != nil {
		logger.Error(nil, "writeHistory Alert[%s] %s in DB error, [%+v].", ar.AlertConfig.AlertId, status, err)
		return
	}
	logger.Debug(nil, "writeHistory Alert[%s] %s in DB successfully.", ar.AlertConfig.AlertId, status)
}

func (ar *AlertRunner) commentAlert(historyId string) {
	//TODO: Find and record which rule id and resource name created message id
	ar.writeHistory("", "commented", historyId, "", "", "")
}

func (ar *AlertRunner) pushAggregatedAlerts(newStatus *StatusResource, ruleId string, resourceName string, triggeredRuleMetrics []RecordedMetric) {
	aggregatedAlerts := newStatus.AggregatedAlerts

	aggregatedAlerts.CumulatedCount = aggregatedAlerts.CumulatedCount + 1

	triggeredMetric := triggeredRuleMetrics[len(triggeredRuleMetrics)-1]
	alertTime := time.Unix(triggeredMetric.tvs[len(triggeredMetric.tvs)-1].T, 0).Format("2006-01-02 15:04:05.99999")
	if aggregatedAlerts.FirstAlertTime == "" {
		aggregatedAlerts.FirstAlertTime = alertTime
	}
	aggregatedAlerts.LastAlertTime = alertTime
	aggregatedAlerts.LastAlertValues = triggeredRuleMetrics

	newStatus.AggregatedAlerts = aggregatedAlerts
}

func (ar *AlertRunner) checkSendable(newStatus *StatusResource, ruleId string, resourceName string) bool {
	policyConfig := ar.AlertConfig.PolicyConfig[ar.AlertConfig.Rules[ruleId].Severity]
	switch policyConfig.RepeatType {
	case "normal":
		return true
	case "not-repeat":
		if newStatus.CumulatedSendCount > 0 {
			return false
		} else {
			return true
		}
	case "fixed-minutes":
		if newStatus.CumulatedSendCount >= policyConfig.MaxSendCount {
			return false
		}
		if !newStatus.NextSendableTime.After(time.Now().Add(time.Duration(TickPeriodSecond/3) * time.Second)) {
			return true
		} else {
			return false
		}
	case "exp-minutes":
		if newStatus.CumulatedSendCount >= policyConfig.MaxSendCount {
			return false
		}
		if !newStatus.NextSendableTime.After(time.Now().Add(time.Duration(TickPeriodSecond/3) * time.Second)) {
			return true
		} else {
			return false
		}
	}
	return false
}

func (ar *AlertRunner) clearAggregatedAlerts(newStatus *StatusResource, ruleId string, resourceName string) {
	newStatus.AggregatedAlerts = AggregatedAlert{}
}

func (ar *AlertRunner) refreshNextSendableTime(newStatus *StatusResource) {
	newStatus.NextSendableTime = time.Now()
}

func (ar *AlertRunner) processRepeat(newStatus *StatusResource, ruleId string, resourceName string) {
	newStatus.CumulatedSendCount = newStatus.CumulatedSendCount + 1

	//Update Next Sendable Time
	switch ar.AlertConfig.PolicyConfig[ar.AlertConfig.Rules[ruleId].Severity].RepeatType {
	case "fixed-minutes":
		newStatus.NextSendableTime = newStatus.NextSendableTime.Add(time.Duration(newStatus.NextResendInterval) * time.Minute)
	case "exp-minutes":
		newStatus.NextSendableTime = newStatus.NextSendableTime.Add(time.Duration(newStatus.NextResendInterval) * time.Minute)
		newStatus.NextResendInterval = newStatus.NextResendInterval * 2
	}
}

func processResourceName(resourceName string) string {
	if strings.Contains(resourceName, ":") {
		return strings.Split(resourceName, ":")[1]
	}

	return resourceName
}

func (ar *AlertRunner) formatActiveNotificationEmail(newStatus *StatusResource, ruleId string, resourceName string) *notification.Email {
	aggregatedAlerts := newStatus.AggregatedAlerts
	lastValue := ""
	for _, recordedRuleMetric := range aggregatedAlerts.LastAlertValues {
		tv := recordedRuleMetric.tvs[len(recordedRuleMetric.tvs)-1]
		if resourceName == recordedRuleMetric.ResourceName {
			v, _ := strconv.ParseFloat(tv.V, 64)
			lastValue = fmt.Sprintf("%.2f%s", v*ar.AlertConfig.Rules[ruleId].Scale, ar.AlertConfig.Rules[ruleId].Unit)
			break
		}
	}

	notificationParam := notification.NotificationParam{
		ResourceName:   processResourceName(resourceName),
		RuleName:       ar.AlertConfig.Rules[ruleId].RuleName,
		CumulatedCount: aggregatedAlerts.CumulatedCount,
		FirstTime:      aggregatedAlerts.FirstAlertTime,
		LastTime:       aggregatedAlerts.LastAlertTime,
		LastValue:      lastValue,
	}

	notificationParamBytes, err := json.Marshal(notificationParam)
	if err != nil {
		logger.Error(nil, "Marshal Notification Param error: %v", err)
		return nil
	}

	emailStr := adapter.SendEmailRequest(string(notificationParamBytes), "false")

	if emailStr == "" {
		return nil
	}

	email := notification.Email{}

	err = json.Unmarshal([]byte(emailStr), &email)
	if err != nil {
		logger.Error(nil, "Unmarshal Email error: %v", err)
		return nil
	}

	return &email
}

func (ar *AlertRunner) formatResumeNotificationEmail(resumeStatus *StatusResource, ruleId string, resourceName string, resumedMetric RecordedMetric) *notification.Email {
	aggregatedAlerts := resumeStatus.AggregatedAlerts
	lastValue := ""
	tv := resumedMetric.tvs[len(resumedMetric.tvs)-1]
	if resourceName == resumedMetric.ResourceName {
		v, _ := strconv.ParseFloat(tv.V, 64)
		lastValue = fmt.Sprintf("%.2f%s", v*ar.AlertConfig.Rules[ruleId].Scale, ar.AlertConfig.Rules[ruleId].Unit)
	}
	resumeTime := time.Unix(tv.T, 0).Format("2006-01-02 15:04:05.99999")

	notificationParam := notification.NotificationParam{
		ResourceName: processResourceName(resourceName),
		RuleName:     ar.AlertConfig.Rules[ruleId].RuleName,
		FirstTime:    aggregatedAlerts.FirstAlertTime,
		LastTime:     resumeTime,
		LastValue:    lastValue,
	}

	notificationParamBytes, err := json.Marshal(notificationParam)
	if err != nil {
		logger.Error(nil, "Marshal Notification Param error: %v", err)
		return nil
	}

	emailStr := adapter.SendEmailRequest(string(notificationParamBytes), "true")

	if emailStr == "" {
		return nil
	}

	email := notification.Email{}

	err = json.Unmarshal([]byte(emailStr), &email)
	if err != nil {
		logger.Error(nil, "Unmarshal Email error: %v", err)
		return nil
	}

	return &email
}

func (ar *AlertRunner) sendActiveNotification(newStatus *StatusResource, ruleId string, resourceName string, triggeredRuleMetrics []RecordedMetric) {
	ar.pushAggregatedAlerts(newStatus, ruleId, resourceName, triggeredRuleMetrics)

	//Check Notification Sendable
	if !nf.CheckTimeAvailable(ar.AlertConfig.AvailableStartTime, ar.AlertConfig.AvailableEndTime) {
		ar.refreshNextSendableTime(newStatus)
		logger.Debug(nil, "sendActiveNotification not in available time")
		return
	}

	//Check Policy Sendable
	if !ar.checkSendable(newStatus, ruleId, resourceName) {
		return
	}

	nfAddressListId := fmt.Sprintf(`["%s"]`, ar.AlertConfig.NfAddressListId)
	email := ar.formatActiveNotificationEmail(newStatus, ruleId, resourceName)
	if email == nil {
		logger.Error(nil, "formatActiveNotificationEmail failed")
	} else {
		sentSuccess, notificationId := nf.SendNotification("other", nfAddressListId, email.Title, email.Content)
		if sentSuccess {
			ar.writeHistory("", "sent_success", fmt.Sprintf("%v", triggeredRuleMetrics), notificationId, ruleId, resourceName)
			//ar.clearAggregatedAlerts(newStatus, ruleId, resourceName)
		} else {
			ar.writeHistory("", "sent_failed", fmt.Sprintf("%v", triggeredRuleMetrics), "", ruleId, resourceName)
			logger.Error(nil, "sendActiveNotification failed")
		}
	}

	ar.processRepeat(newStatus, ruleId, resourceName)
}

func (ar *AlertRunner) sendResumeNotification(resumeStatus *StatusResource, ruleId string, resourceName string, resumedMetric RecordedMetric, resumedMetrics []RecordedMetric) {
	//Check Policy Sendable
	if !ar.checkSendable(resumeStatus, ruleId, resourceName) {
		return
	}

	nfAddressListId := fmt.Sprintf(`["%s"]`, ar.AlertConfig.NfAddressListId)
	email := ar.formatResumeNotificationEmail(resumeStatus, ruleId, resourceName, resumedMetric)
	if email == nil {
		logger.Error(nil, "formatResumeNotificationEmail failed")
	} else {
		sentSuccess, notificationId := nf.SendNotification("other", nfAddressListId, email.Title, email.Content)
		if sentSuccess {
			ar.writeHistory("", "sent_success", fmt.Sprintf("%v", resumedMetrics), notificationId, ruleId, resourceName)
		} else {
			ar.writeHistory("", "sent_failed", fmt.Sprintf("%v", resumedMetrics), "", ruleId, resourceName)
			logger.Error(nil, "sendResumeNotification failed")
		}
	}
}

func (ar *AlertRunner) updateAlertUpdateTime() {
	ar.AlertStatus.Lock()
	ar.AlertStatus.UpdateTime = time.Now()
	ar.AlertStatus.Unlock()
}

func (ar *AlertRunner) GetAlertStatus() (string, time.Time) {
	alertStatus := ""
	updateTime := ar.AlertStatus.UpdateTime

	ar.AlertStatus.Lock()
	alertStatusBytes, err := json.Marshal(ar.AlertStatus)
	if err == nil {
		alertStatus = string(alertStatusBytes)
	}
	updateTime = ar.AlertStatus.UpdateTime
	ar.AlertStatus.Unlock()

	return alertStatus, updateTime
}

func (ar *AlertRunner) runAlertRules() {
	if !ar.AlertConfig.LoadSuccess {
		return
	}

	if ar.AlertConfig.Disabled {
		return
	}

	ch := make(chan metric.ResourceMetrics, 100)
	ar.getResourceMetrics(ch)
	close(ch)

	ar.checkMetrics(ch)
}

func (ar *AlertRunner) Run(initStatus string) {
	timer := time.NewTicker(time.Second * TickPeriodSecond)
	defer timer.Stop()

	ar.loadAlertInfo()
	ar.updateAlertUpdateTime()

	//If get alert from migrating, continue running with current status, only need to reset status when adding and updating
	if initStatus == "adding" {
		ar.AlertStatus.Lock()
		ar.resetAlertStatus()
		ar.AlertStatus.Unlock()
	}

	for {
		select {
		case <-timer.C:
			ar.runAlertRules()
			logger.Debug(nil, "AlertRunner alert %s run", ar.AlertConfig.AlertId)
			ar.updateAlertUpdateTime()
		case operation := <-ar.SignalCh:
			switch operation {
			case "Stop":
				//Drain SignalCh
				for len(ar.SignalCh) > 0 {
					<-ar.SignalCh
				}
				logger.Debug(nil, "AlertRunner alert %s stop", ar.AlertConfig.AlertId)
				return
			case "Update":
				ar.loadAlertInfo()
				ar.AlertStatus.Lock()
				ar.resetAlertStatus()
				ar.AlertStatus.Unlock()
				ar.updateAlertUpdateTime()
				logger.Debug(nil, "AlertRunner alert %s update", ar.AlertConfig.AlertId)
			default:
				param := strings.Split(operation, " ")
				if len(param) != 2 {
					break
				}
				switch param[0] {
				case "Comment":
					ar.commentAlert(param[1])
					ar.updateAlertUpdateTime()
					logger.Debug(nil, "AlertRunner alert %s comment", ar.AlertConfig.AlertId)
				}
			}
		}
	}
}
