package manager

import (
	"context"

	"kubesphere.io/alert/pkg/gerr"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	. "kubesphere.io/alert/pkg/pb"
	rs "kubesphere.io/alert/pkg/services/manager/resource_control"
	"kubesphere.io/alert/pkg/util/stringutil"
)

func (s *Server) createAlert(ctx context.Context, req *CreateAlertRequest) (string, error) {
	alert := models.NewAlert(
		req.GetAlertName(),
		req.GetDisabled(),
		"adding",
		"",
		req.GetPolicyId(),
		req.GetRsFilterId(),
	)

	err := rs.CreateAlert(ctx, alert)
	if err != nil {
		logger.Error(ctx, "Manager create Alert[%s] in DB error, [%+v].", req.GetAlertName(), err)
		return "", err
	}
	logger.Debug(ctx, "Manager create Alert[%s][%s] in DB successfully.", alert.AlertId, req.GetAlertName())

	// Enqueue alert after create tasks.
	err = s.alertQueue.Enqueue(alert.AlertId)
	if err != nil {
		logger.Error(ctx, "Manager push alert[%s] into etcd failed, [%+v].", alert.AlertId, err)
		return alert.AlertId, err
	}
	logger.Debug(ctx, "Manager push alert[%s] into etcd successfully.", alert.AlertId)

	return alert.AlertId, nil
}

func (s *Server) operateAlert(ctx context.Context, req *ModifyAlertRequest, alertId string, operation string) (string, error) {
	alertId, err := rs.UpdateAlert(ctx, req, alertId, operation)
	if err != nil {
		logger.Error(ctx, "Manager update Alert[%s] in DB error, [%+v].", alertId, err)
		return "", gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, alertId)
	}

	// Broadcast alert operation after update DB.
	err = s.alertBroadcast.Broadcast(alertId, operation, 10)
	if err != nil {
		logger.Error(ctx, "Manager broadast alert %s[%s] into etcd failed, [%+v].", operation, alertId, err)
		return "", err
	}
	logger.Debug(ctx, "Manager broadast alert %s[%s] into etcd successfully.", operation, alertId)

	return alertId, nil
}

//0.Executor
//********************************************************************************************************
func (s *Server) CreateExecutor(context.Context, *CreateExecutorRequest) (*CreateExecutorResponse, error) {
	return &CreateExecutorResponse{}, nil
}

func (s *Server) DescribeExecutors(context.Context, *DescribeExecutorsRequest) (*DescribeExecutorsResponse, error) {
	return &DescribeExecutorsResponse{}, nil
}

func (s *Server) ModifyExecutor(context.Context, *ModifyExecutorRequest) (*ModifyExecutorResponse, error) {
	return &ModifyExecutorResponse{}, nil
}

func (s *Server) DeleteExecutors(context.Context, *DeleteExecutorsRequest) (*DeleteExecutorsResponse, error) {
	return &DeleteExecutorsResponse{}, nil
}

//1.ResourceType
//********************************************************************************************************
func (s *Server) CreateResourceType(ctx context.Context, req *CreateResourceTypeRequest) (*CreateResourceTypeResponse, error) {
	err := ValidateCreateResourceTypeParams(ctx, req)
	if err != nil {
		return nil, err
	}

	resourceType := models.NewResourceType(
		req.GetRsTypeName(),
		req.GetRsTypeParam(),
	)

	err = rs.CreateResourceType(ctx, resourceType)
	if err != nil {
		return nil, err
	}
	logger.Debug(ctx, "Create ResourceType[%s] in DB successfully.", resourceType.RsTypeId)

	return &CreateResourceTypeResponse{RsTypeId: resourceType.RsTypeId}, nil
}

func (s *Server) DescribeResourceTypes(ctx context.Context, req *DescribeResourceTypesRequest) (*DescribeResourceTypesResponse, error) {
	rts, rtCnt, err := rs.DescribeResourceTypes(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe ResourceTypes, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	rtPbSet := models.ParseRtSet2PbSet(rts)
	res := &DescribeResourceTypesResponse{
		Total:           uint32(rtCnt),
		ResourceTypeSet: rtPbSet,
	}

	logger.Debug(ctx, "Describe ResourceTypes successfully, ResourceTypes=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyResourceType(ctx context.Context, req *ModifyResourceTypeRequest) (*ModifyResourceTypeResponse, error) {
	err := ValidateModifyResourceTypeParams(ctx, req)
	if err != nil {
		return nil, err
	}

	rsTypeId, err := rs.ModifyResourceType(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Modify ResourceType[%s], [%+v].", rsTypeId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, rsTypeId)
	}
	logger.Debug(ctx, "Modify ResourceType[%s] successfully.", rsTypeId)
	return &ModifyResourceTypeResponse{
		RsTypeId: rsTypeId,
	}, nil
}

func (s *Server) DeleteResourceTypes(ctx context.Context, req *DeleteResourceTypesRequest) (*DeleteResourceTypesResponse, error) {
	rsTypeIds, err := rs.DeleteResourceTypes(ctx, stringutil.SimplifyStringList(req.RsTypeId))
	if err != nil {
		logger.Error(ctx, "Failed to Delete ResourceTypes[%+v], [%+v].", rsTypeIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, rsTypeIds)
	}
	logger.Debug(ctx, "Delete ResourceTypes[%+v] successfully.", rsTypeIds)
	return &DeleteResourceTypesResponse{
		RsTypeId: rsTypeIds,
	}, nil
}

//2.Resource Filter
//********************************************************************************************************
func (s *Server) CreateResourceFilter(ctx context.Context, req *CreateResourceFilterRequest) (*CreateResourceFilterResponse, error) {
	err := ValidateCreateResourceFilterParams(ctx, req)
	if err != nil {
		return nil, err
	}

	rsFilter := models.NewResourceFilter(
		req.GetRsFilterName(),
		req.GetRsFilterParam(),
		req.GetRsTypeId(),
	)

	err = rs.CreateResourceFilter(ctx, rsFilter)
	if err != nil {
		return nil, err
	}
	logger.Debug(ctx, "Create ResourceFilter[%s] in DB successfully.", rsFilter.RsFilterId)

	return &CreateResourceFilterResponse{RsFilterId: rsFilter.RsFilterId}, nil
}

func (s *Server) DescribeResourceFilters(ctx context.Context, req *DescribeResourceFiltersRequest) (*DescribeResourceFiltersResponse, error) {
	rfs, rfCnt, err := rs.DescribeResourceFilters(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Resource Filters, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	rfPbSet := models.ParseRfSet2PbSet(rfs)
	res := &DescribeResourceFiltersResponse{
		Total:             uint32(rfCnt),
		ResourcefilterSet: rfPbSet,
	}

	logger.Debug(ctx, "Describe Resource Filters successfully, Resource Filters=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyResourceFilter(ctx context.Context, req *ModifyResourceFilterRequest) (*ModifyResourceFilterResponse, error) {
	err := ValidateModifyResourceFilterParams(ctx, req)
	if err != nil {
		return nil, err
	}

	rsFilterId, err := rs.ModifyResourceFilter(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Resource Filter[%s], [%+v].", rsFilterId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, rsFilterId)
	}
	logger.Debug(ctx, "Modify Resource Filter[%s] successfully.", rsFilterId)
	return &ModifyResourceFilterResponse{
		RsFilterId: rsFilterId,
	}, nil
}

func (s *Server) DeleteResourceFilters(ctx context.Context, req *DeleteResourceFiltersRequest) (*DeleteResourceFiltersResponse, error) {
	rsFilterIds, err := rs.DeleteResourceFilters(ctx, stringutil.SimplifyStringList(req.RsFilterId))
	if err != nil {
		logger.Error(ctx, "Failed to Delete Resource Filters[%+v], [%+v].", rsFilterIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, rsFilterIds)
	}
	logger.Debug(ctx, "Delete Resource Filters[%+v] successfully.", rsFilterIds)
	return &DeleteResourceFiltersResponse{
		RsFilterId: rsFilterIds,
	}, nil
}

//3.Metric
//********************************************************************************************************
func (s *Server) CreateMetric(ctx context.Context, req *CreateMetricRequest) (*CreateMetricResponse, error) {
	err := ValidateCreateMetricParams(ctx, req)
	if err != nil {
		return nil, err
	}

	metric := models.NewMetric(
		req.GetMetricName(),
		req.GetMetricParam(),
		req.GetRsTypeId(),
	)

	err = rs.CreateMetric(ctx, metric)
	if err != nil {
		return nil, err
	}
	logger.Debug(ctx, "Create Metric[%s] in DB successfully.", metric.MetricId)

	return &CreateMetricResponse{MetricId: metric.MetricId}, nil
}

func (s *Server) DescribeMetrics(ctx context.Context, req *DescribeMetricsRequest) (*DescribeMetricsResponse, error) {
	mts, mtCnt, err := rs.DescribeMetrics(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Metrics, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	mtPbSet := models.ParseMtSet2PbSet(mts)
	res := &DescribeMetricsResponse{
		Total:     uint32(mtCnt),
		MetricSet: mtPbSet,
	}

	logger.Debug(ctx, "Describe Metrics successfully, Metrics=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyMetric(ctx context.Context, req *ModifyMetricRequest) (*ModifyMetricResponse, error) {
	err := ValidateModifyMetricParams(ctx, req)
	if err != nil {
		return nil, err
	}

	metricId, err := rs.ModifyMetric(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Metric[%s], [%+v].", metricId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, metricId)
	}
	logger.Debug(ctx, "Modify Metric[%s] successfully.", metricId)
	return &ModifyMetricResponse{
		MetricId: metricId,
	}, nil
}

func (s *Server) DeleteMetrics(ctx context.Context, req *DeleteMetricsRequest) (*DeleteMetricsResponse, error) {
	metricIds, err := rs.DeleteMetrics(ctx, stringutil.SimplifyStringList(req.MetricId))
	if err != nil {
		logger.Error(ctx, "Failed to Delete Metrics[%+v], [%+v].", metricIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, metricIds)
	}
	logger.Debug(ctx, "Delete Metrics[%+v] successfully.", metricIds)
	return &DeleteMetricsResponse{
		MetricId: metricIds,
	}, nil
}

//4.Policy
//********************************************************************************************************
func (s *Server) CreatePolicy(ctx context.Context, req *CreatePolicyRequest) (*CreatePolicyResponse, error) {
	err := ValidateCreatePolicyParams(ctx, req)
	if err != nil {
		return nil, err
	}

	policy := models.NewPolicy(
		req.GetPolicyName(),
		req.GetPolicyDescription(),
		req.GetPolicyConfig(),
		req.GetCreator(),
		req.GetAvailableStartTime(),
		req.GetAvailableEndTime(),
		req.GetRsTypeId(),
		req.GetLanguage(),
	)

	err = rs.CreatePolicy(ctx, policy)
	if err != nil {
		return nil, err
	}
	logger.Debug(ctx, "Create Policy[%s] in DB successfully.", policy.PolicyId)

	return &CreatePolicyResponse{PolicyId: policy.PolicyId}, nil
}

func (s *Server) DescribePolicies(ctx context.Context, req *DescribePoliciesRequest) (*DescribePoliciesResponse, error) {
	pls, plCnt, err := rs.DescribePolicies(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Policies, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	plPbSet := models.ParsePlSet2PbSet(pls)
	res := &DescribePoliciesResponse{
		Total:     uint32(plCnt),
		PolicySet: plPbSet,
	}

	logger.Debug(ctx, "Describe Policies successfully, Policies=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyPolicy(ctx context.Context, req *ModifyPolicyRequest) (*ModifyPolicyResponse, error) {
	err := ValidateModifyPolicyParams(ctx, req)
	if err != nil {
		return nil, err
	}

	policyId, err := rs.ModifyPolicy(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Policy[%s], [%+v].", policyId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, policyId)
	}
	logger.Debug(ctx, "Modify Policy[%s] successfully.", policyId)
	return &ModifyPolicyResponse{
		PolicyId: policyId,
	}, nil
}

func (s *Server) DeletePolicies(ctx context.Context, req *DeletePoliciesRequest) (*DeletePoliciesResponse, error) {
	policyIds, err := rs.DeletePolicies(ctx, stringutil.SimplifyStringList(req.PolicyId))
	if err != nil {
		logger.Error(ctx, "Failed to Delete Policies[%+v], [%+v].", policyIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, policyIds)
	}
	logger.Debug(ctx, "Delete Policies[%+v] successfully.", policyIds)
	return &DeletePoliciesResponse{
		PolicyId: policyIds,
	}, nil
}

//5.Rule
//********************************************************************************************************
func (s *Server) CreateRule(ctx context.Context, req *CreateRuleRequest) (*CreateRuleResponse, error) {
	err := ValidateCreateRuleParams(ctx, req)
	if err != nil {
		return nil, err
	}

	rule := models.NewRule(
		req.GetRuleName(),
		req.GetDisabled(),
		req.GetMonitorPeriods(),
		req.GetSeverity(),
		req.GetMetricsType(),
		req.GetConditionType(),
		req.GetThresholds(),
		req.GetUnit(),
		req.GetConsecutiveCount(),
		req.GetInhibit(),
		req.GetPolicyId(),
		req.GetMetricId(),
	)

	err = rs.CreateRule(ctx, rule)
	if err != nil {
		return nil, err
	}
	logger.Debug(ctx, "Create Rule[%s] in DB successfully.", rule.RuleId)

	return &CreateRuleResponse{RuleId: rule.RuleId}, nil
}

func (s *Server) DescribeRules(ctx context.Context, req *DescribeRulesRequest) (*DescribeRulesResponse, error) {
	rls, rlCnt, err := rs.DescribeRules(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Rules, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	rlPbSet := models.ParseRlSet2PbSet(rls)
	res := &DescribeRulesResponse{
		Total:   uint32(rlCnt),
		RuleSet: rlPbSet,
	}

	logger.Debug(ctx, "Describe Rules successfully, Rules=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyRule(ctx context.Context, req *ModifyRuleRequest) (*ModifyRuleResponse, error) {
	err := ValidateModifyRuleParams(ctx, req)
	if err != nil {
		return nil, err
	}

	ruleId, err := rs.ModifyRule(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Rule[%s], [%+v].", ruleId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, ruleId)
	}
	logger.Debug(ctx, "Modify Rule[%s] successfully.", ruleId)
	return &ModifyRuleResponse{
		RuleId: ruleId,
	}, nil
}

func (s *Server) DeleteRules(ctx context.Context, req *DeleteRulesRequest) (*DeleteRulesResponse, error) {
	ruleIds, err := rs.DeleteRules(ctx, stringutil.SimplifyStringList(req.RuleId))
	if err != nil {
		logger.Error(ctx, "Failed to Delete Rules[%+v], [%+v].", ruleIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, ruleIds)
	}
	logger.Debug(ctx, "Delete Rules[%+v] successfully.", ruleIds)
	return &DeleteRulesResponse{
		RuleId: ruleIds,
	}, nil
}

//6.Alert
//********************************************************************************************************
func (s *Server) CreateAlert(ctx context.Context, req *CreateAlertRequest) (*CreateAlertResponse, error) {
	err := ValidateCreateAlertParams(ctx, req)
	if err != nil {
		return nil, err
	}

	alertId, err := s.createAlert(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Create Alert, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	return &CreateAlertResponse{AlertId: alertId}, nil
}

func (s *Server) DescribeAlerts(ctx context.Context, req *DescribeAlertsRequest) (*DescribeAlertsResponse, error) {
	als, alCnt, err := rs.DescribeAlerts(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Alerts, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	alPbSet := models.ParseAlSet2PbSet(als)
	res := &DescribeAlertsResponse{
		Total:    uint32(alCnt),
		AlertSet: alPbSet,
	}

	logger.Debug(ctx, "Describe Alerts successfully, Alerts=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyAlert(ctx context.Context, req *ModifyAlertRequest) (*ModifyAlertResponse, error) {
	err := ValidateModifyAlertParams(ctx, req)
	if err != nil {
		return nil, err
	}

	alertId, err := s.operateAlert(ctx, req, req.AlertId, "updating")
	if err != nil {
		logger.Error(ctx, "Failed to Modify Alert, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	logger.Debug(ctx, "Modify Alert[%s] successfully.", alertId)
	return &ModifyAlertResponse{
		AlertId: alertId,
	}, nil
}

func (s *Server) DeleteAlerts(ctx context.Context, req *DeleteAlertsRequest) (*DeleteAlertsResponse, error) {
	alertIds := stringutil.SimplifyStringList(req.AlertId)
	alertIdsSuccess := []string{}
	for _, alertId := range alertIds {
		alertIdSuccess, err := s.operateAlert(ctx, nil, alertId, "deleting")
		if err != nil {
			logger.Error(ctx, "Failed to Delete Alert[%s], [%+v].", alertId, err)
		} else {
			alertIdsSuccess = append(alertIdsSuccess, alertIdSuccess)
		}
	}

	logger.Debug(ctx, "Delete Alerts[%+v] successfully.", alertIdsSuccess)
	return &DeleteAlertsResponse{
		AlertId: alertIdsSuccess,
	}, nil
}

//7.History
//********************************************************************************************************
func (s *Server) CreateHistory(context.Context, *CreateHistoryRequest) (*CreateHistoryResponse, error) {
	return &CreateHistoryResponse{}, nil
}

func (s *Server) DescribeHistories(ctx context.Context, req *DescribeHistoriesRequest) (*DescribeHistoriesResponse, error) {
	hss, hsCnt, err := rs.DescribeHistories(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Histories, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	hsPbSet := models.ParseHsSet2PbSet(hss)
	res := &DescribeHistoriesResponse{
		Total:      uint32(hsCnt),
		HistorySet: hsPbSet,
	}

	logger.Debug(ctx, "Describe Histories successfully, Histories=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyHistory(context.Context, *ModifyHistoryRequest) (*ModifyHistoryResponse, error) {
	return &ModifyHistoryResponse{}, nil
}

func (s *Server) DeleteHistories(context.Context, *DeleteHistoriesRequest) (*DeleteHistoriesResponse, error) {
	return &DeleteHistoriesResponse{}, nil
}

//8.Comment
//********************************************************************************************************
func (s *Server) CreateComment(ctx context.Context, req *CreateCommentRequest) (*CreateCommentResponse, error) {
	err := ValidateCreateCommentParams(ctx, req)
	if err != nil {
		return nil, err
	}

	alert := rs.GetAlertByHistoryId(req.GetHistoryId())

	if alert.AlertId == "" {
		logger.Error(ctx, "Create Comment history_id [%s] does not exist.", req.GetHistoryId())
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, nil, gerr.ErrorCreateResourcesFailed)
	}

	comment := models.NewComment(
		req.GetAddresser(),
		req.GetContent(),
		req.GetHistoryId(),
	)

	err = rs.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	logger.Debug(ctx, "Create Comment[%s] in DB successfully.", comment.CommentId)

	// Broadcast alert comment after update DB.
	alertId := alert.AlertId
	operation := "commenting " + req.GetHistoryId()
	err = s.alertBroadcast.Broadcast(alertId, operation, 10)
	if err != nil {
		logger.Error(ctx, "Manager broadast alert %s[%s] into etcd failed, [%+v].", operation, alertId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, nil, gerr.ErrorCreateResourcesFailed)
	}
	logger.Debug(ctx, "Manager broadast alert %s[%s] into etcd successfully.", operation, alertId)

	return &CreateCommentResponse{CommentId: comment.CommentId}, nil
}

func (s *Server) DescribeComments(ctx context.Context, req *DescribeCommentsRequest) (*DescribeCommentsResponse, error) {
	cms, cmCnt, err := rs.DescribeComments(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Comments, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	cmPbSet := models.ParseCmtSet2PbSet(cms)
	res := &DescribeCommentsResponse{
		Total:      uint32(cmCnt),
		CommentSet: cmPbSet,
	}

	logger.Debug(ctx, "Describe Comments successfully, Comments=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyComment(context.Context, *ModifyCommentRequest) (*ModifyCommentResponse, error) {
	return &ModifyCommentResponse{}, nil
}

func (s *Server) DeleteComments(context.Context, *DeleteCommentsRequest) (*DeleteCommentsResponse, error) {
	return &DeleteCommentsResponse{}, nil
}

//9.Action
//********************************************************************************************************
func (s *Server) CreateAction(ctx context.Context, req *CreateActionRequest) (*CreateActionResponse, error) {
	err := ValidateCreateActionParams(ctx, req)
	if err != nil {
		return nil, err
	}

	action := models.NewAction(
		req.GetActionName(),
		req.GetTriggerStatus(),
		req.GetTriggerAction(),
		req.GetPolicyId(),
		req.GetNfAddressListId(),
	)

	err = rs.CreateAction(ctx, action)
	if err != nil {
		return nil, err
	}
	logger.Debug(ctx, "Create Action[%s] in DB successfully.", action.ActionId)

	return &CreateActionResponse{ActionId: action.ActionId}, nil
}

func (s *Server) DescribeActions(ctx context.Context, req *DescribeActionsRequest) (*DescribeActionsResponse, error) {
	acs, acCnt, err := rs.DescribeActions(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Actions, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	acPbSet := models.ParseAcSet2PbSet(acs)
	res := &DescribeActionsResponse{
		Total:     uint32(acCnt),
		ActionSet: acPbSet,
	}

	logger.Debug(ctx, "Describe Actions successfully, Actions=[%+v].", res)
	return res, nil
}

func (s *Server) ModifyAction(ctx context.Context, req *ModifyActionRequest) (*ModifyActionResponse, error) {
	err := ValidateModifyActionParams(ctx, req)
	if err != nil {
		return nil, err
	}

	actionId, err := rs.ModifyAction(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Action[%s], [%+v].", actionId, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, actionId)
	}
	logger.Debug(ctx, "Modify Action[%s] successfully.", actionId)
	return &ModifyActionResponse{
		ActionId: actionId,
	}, nil
}

func (s *Server) DeleteActions(ctx context.Context, req *DeleteActionsRequest) (*DeleteActionsResponse, error) {
	actionIds, err := rs.DeleteActions(ctx, stringutil.SimplifyStringList(req.ActionId))
	if err != nil {
		logger.Error(ctx, "Failed to Delete Actions[%+v], [%+v].", actionIds, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDeleteResourceFailed, actionIds)
	}
	logger.Debug(ctx, "Delete Actions[%+v] successfully.", actionIds)
	return &DeleteActionsResponse{
		ActionId: actionIds,
	}, nil
}
