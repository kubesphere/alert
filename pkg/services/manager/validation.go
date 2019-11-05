// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"context"
	"time"

	"kubesphere.io/alert/pkg/gerr"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/pb"
)

func checkStringLen(ctx context.Context, str string, length int) error {
	if len(str) <= length {
		return nil
	} else {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorStringLengthExceed, str, length)
	}
}

func checkTimeFormat(ctx context.Context, timeStr string) error {
	timeFmt := "15:04:05"

	_, err := time.Parse(timeFmt, timeStr)

	if err == nil {
		return nil
	} else {
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalTimeFormat, timeStr)
	}
}

func ValidateCreateResourceTypeParams(ctx context.Context, req *pb.CreateResourceTypeRequest) error {
	rsTypeName := req.GetRsTypeName()
	err := checkStringLen(ctx, rsTypeName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsTypeName [%s]: %+v", rsTypeName, err)
		return err
	}

	return nil
}

func ValidateModifyResourceTypeParams(ctx context.Context, req *pb.ModifyResourceTypeRequest) error {
	rsTypeId := req.GetRsTypeId()
	err := checkStringLen(ctx, rsTypeId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate rsTypeId [%s]: %+v", rsTypeId, err)
		return err
	}

	rsTypeName := req.GetRsTypeName()
	err = checkStringLen(ctx, rsTypeName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsTypeName [%s]: %+v", rsTypeName, err)
		return err
	}

	return nil
}

func ValidateCreateResourceFilterParams(ctx context.Context, req *pb.CreateResourceFilterRequest) error {
	rsFilterName := req.GetRsFilterName()
	err := checkStringLen(ctx, rsFilterName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsFilterName [%s]: %+v", rsFilterName, err)
		return err
	}

	return nil
}

func ValidateModifyResourceFilterParams(ctx context.Context, req *pb.ModifyResourceFilterRequest) error {
	rsFilterId := req.GetRsFilterId()
	err := checkStringLen(ctx, rsFilterId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsFilterId [%s]: %+v", rsFilterId, err)
		return err
	}

	rsFilterName := req.GetRsFilterName()
	err = checkStringLen(ctx, rsFilterName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsFilterName [%s]: %+v", rsFilterName, err)
		return err
	}

	return nil
}

func ValidateCreateMetricParams(ctx context.Context, req *pb.CreateMetricRequest) error {
	metricName := req.GetMetricName()
	err := checkStringLen(ctx, metricName, 100)
	if err != nil {
		logger.Error(ctx, "Failed to validate MetricName [%s]: %+v", metricName, err)
		return err
	}

	rsTypeId := req.GetRsTypeId()
	err = checkStringLen(ctx, rsTypeId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsTypeId [%s]: %+v", rsTypeId, err)
		return err
	}

	return nil
}

func ValidateModifyMetricParams(ctx context.Context, req *pb.ModifyMetricRequest) error {
	metricId := req.GetMetricId()
	err := checkStringLen(ctx, metricId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate MetricId [%s]: %+v", metricId, err)
		return err
	}

	metricName := req.GetMetricName()
	err = checkStringLen(ctx, metricName, 100)
	if err != nil {
		logger.Error(ctx, "Failed to validate MetricName [%s]: %+v", metricName, err)
		return err
	}

	rsTypeId := req.GetRsTypeId()
	err = checkStringLen(ctx, rsTypeId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsTypeId [%s]: %+v", rsTypeId, err)
		return err
	}

	return nil
}

func ValidateCreatePolicyParams(ctx context.Context, req *pb.CreatePolicyRequest) error {
	policyName := req.GetPolicyName()
	err := checkStringLen(ctx, policyName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyName [%s]: %+v", policyName, err)
		return err
	}

	policyDescription := req.GetPolicyDescription()
	err = checkStringLen(ctx, policyDescription, 255)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyDescription [%s]: %+v", policyDescription, err)
		return err
	}

	creator := req.GetCreator()
	err = checkStringLen(ctx, creator, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate Creator [%s]: %+v", creator, err)
		return err
	}

	availableStartTime := req.GetAvailableStartTime()
	err = checkTimeFormat(ctx, availableStartTime)
	if err != nil {
		logger.Error(ctx, "Failed to validate AvailableStartTime [%s]: %+v", availableStartTime, err)
		return err
	}

	availableEndTime := req.GetAvailableEndTime()
	err = checkTimeFormat(ctx, availableEndTime)
	if err != nil {
		logger.Error(ctx, "Failed to validate AvailableEndTime [%s]: %+v", availableEndTime, err)
		return err
	}

	rsTypeId := req.GetRsTypeId()
	err = checkStringLen(ctx, rsTypeId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsTypeId [%s]: %+v", rsTypeId, err)
		return err
	}

	return nil
}

func ValidateModifyPolicyParams(ctx context.Context, req *pb.ModifyPolicyRequest) error {
	policyId := req.GetPolicyId()
	err := checkStringLen(ctx, policyId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyId [%s]: %+v", policyId, err)
		return err
	}

	policyName := req.GetPolicyName()
	err = checkStringLen(ctx, policyName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyName [%s]: %+v", policyName, err)
		return err
	}

	policyDescription := req.GetPolicyDescription()
	err = checkStringLen(ctx, policyDescription, 255)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyDescription [%s]: %+v", policyDescription, err)
		return err
	}

	creator := req.GetCreator()
	err = checkStringLen(ctx, creator, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate Creator [%s]: %+v", creator, err)
		return err
	}

	/*availableStartTime := req.GetAvailableStartTime()
	err = checkTimeFormat(ctx, availableStartTime)
	if err != nil {
		logger.Error(ctx, "Failed to validate AvailableStartTime [%s]: %+v", availableStartTime, err)
		return err
	}

	availableEndTime := req.GetAvailableEndTime()
	err = checkTimeFormat(ctx, availableEndTime)
	if err != nil {
		logger.Error(ctx, "Failed to validate AvailableEndTime [%s]: %+v", availableEndTime, err)
		return err
	}*/

	rsTypeId := req.GetRsTypeId()
	err = checkStringLen(ctx, rsTypeId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsTypeId [%s]: %+v", rsTypeId, err)
		return err
	}

	return nil
}

func ValidateCreateRuleParams(ctx context.Context, req *pb.CreateRuleRequest) error {
	ruleName := req.GetRuleName()
	err := checkStringLen(ctx, ruleName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RuleName [%s]: %+v", ruleName, err)
		return err
	}

	severity := req.GetSeverity()
	err = checkStringLen(ctx, severity, 20)
	if err != nil {
		logger.Error(ctx, "Failed to validate Severity [%s]: %+v", severity, err)
		return err
	}

	metricsType := req.GetMetricsType()
	err = checkStringLen(ctx, metricsType, 10)
	if err != nil {
		logger.Error(ctx, "Failed to validate MetricsType [%s]: %+v", metricsType, err)
		return err
	}

	conditionType := req.GetConditionType()
	err = checkStringLen(ctx, conditionType, 10)
	if err != nil {
		logger.Error(ctx, "Failed to validate MetricsType [%s]: %+v", conditionType, err)
		return err
	}

	thresholds := req.GetThresholds()
	err = checkStringLen(ctx, thresholds, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate Thresholds [%s]: %+v", thresholds, err)
		return err
	}

	unit := req.GetUnit()
	err = checkStringLen(ctx, unit, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate Unit [%s]: %+v", unit, err)
		return err
	}

	policyId := req.GetPolicyId()
	err = checkStringLen(ctx, policyId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyId [%s]: %+v", policyId, err)
		return err
	}

	metricId := req.GetMetricId()
	err = checkStringLen(ctx, metricId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate MetricId [%s]: %+v", metricId, err)
		return err
	}

	return nil
}

func ValidateModifyRuleParams(ctx context.Context, req *pb.ModifyRuleRequest) error {
	ruleId := req.GetRuleId()
	err := checkStringLen(ctx, ruleId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RuleId [%s]: %+v", ruleId, err)
		return err
	}

	ruleName := req.GetRuleName()
	err = checkStringLen(ctx, ruleName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RuleName [%s]: %+v", ruleName, err)
		return err
	}

	severity := req.GetSeverity()
	err = checkStringLen(ctx, severity, 20)
	if err != nil {
		logger.Error(ctx, "Failed to validate Severity [%s]: %+v", severity, err)
		return err
	}

	metricsType := req.GetMetricsType()
	err = checkStringLen(ctx, metricsType, 10)
	if err != nil {
		logger.Error(ctx, "Failed to validate MetricsType [%s]: %+v", metricsType, err)
		return err
	}

	conditionType := req.GetConditionType()
	err = checkStringLen(ctx, conditionType, 10)
	if err != nil {
		logger.Error(ctx, "Failed to validate MetricsType [%s]: %+v", conditionType, err)
		return err
	}

	thresholds := req.GetThresholds()
	err = checkStringLen(ctx, thresholds, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate Thresholds [%s]: %+v", thresholds, err)
		return err
	}

	unit := req.GetUnit()
	err = checkStringLen(ctx, unit, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate Unit [%s]: %+v", unit, err)
		return err
	}

	return nil
}

func ValidateCreateAlertParams(ctx context.Context, req *pb.CreateAlertRequest) error {
	alertName := req.GetAlertName()
	err := checkStringLen(ctx, alertName, 100)
	if err != nil {
		logger.Error(ctx, "Failed to validate AlertName [%s]: %+v", alertName, err)
		return err
	}

	policyId := req.GetPolicyId()
	err = checkStringLen(ctx, policyId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyId [%s]: %+v", policyId, err)
		return err
	}

	rsFilterId := req.GetRsFilterId()
	err = checkStringLen(ctx, rsFilterId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsFilterId [%s]: %+v", rsFilterId, err)
		return err
	}

	return nil
}

func ValidateModifyAlertParams(ctx context.Context, req *pb.ModifyAlertRequest) error {
	alertId := req.GetAlertId()
	err := checkStringLen(ctx, alertId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate AlertId [%s]: %+v", alertId, err)
		return err
	}

	alertName := req.GetAlertName()
	err = checkStringLen(ctx, alertName, 100)
	if err != nil {
		logger.Error(ctx, "Failed to validate AlertName [%s]: %+v", alertName, err)
		return err
	}

	policyId := req.GetPolicyId()
	err = checkStringLen(ctx, policyId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyId [%s]: %+v", policyId, err)
		return err
	}

	rsFilterId := req.GetRsFilterId()
	err = checkStringLen(ctx, rsFilterId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate RsFilterId [%s]: %+v", rsFilterId, err)
		return err
	}

	return nil
}

func ValidateCreateCommentParams(ctx context.Context, req *pb.CreateCommentRequest) error {
	addresser := req.GetAddresser()
	err := checkStringLen(ctx, addresser, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate Addresser [%s]: %+v", addresser, err)
		return err
	}

	content := req.GetContent()
	err = checkStringLen(ctx, content, 255)
	if err != nil {
		logger.Error(ctx, "Failed to validate Content [%s]: %+v", content, err)
		return err
	}

	historyId := req.GetHistoryId()
	err = checkStringLen(ctx, historyId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate HistoryId [%s]: %+v", historyId, err)
		return err
	}

	return nil
}

func ValidateCreateActionParams(ctx context.Context, req *pb.CreateActionRequest) error {
	actionName := req.GetActionName()
	err := checkStringLen(ctx, actionName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate ActionName [%s]: %+v", actionName, err)
		return err
	}

	triggerStatus := req.GetTriggerStatus()
	err = checkStringLen(ctx, triggerStatus, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate TriggerStatus [%s]: %+v", triggerStatus, err)
		return err
	}

	triggerAction := req.GetTriggerAction()
	err = checkStringLen(ctx, triggerAction, 255)
	if err != nil {
		logger.Error(ctx, "Failed to validate TriggerAction [%s]: %+v", triggerAction, err)
		return err
	}

	policyId := req.GetPolicyId()
	err = checkStringLen(ctx, policyId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyId [%s]: %+v", policyId, err)
		return err
	}

	nfAddressListId := req.GetNfAddressListId()
	err = checkStringLen(ctx, nfAddressListId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate NfAddressListId [%s]: %+v", nfAddressListId, err)
		return err
	}

	return nil
}

func ValidateModifyActionParams(ctx context.Context, req *pb.ModifyActionRequest) error {
	actionId := req.GetActionId()
	err := checkStringLen(ctx, actionId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate ActionId [%s]: %+v", actionId, err)
		return err
	}

	actionName := req.GetActionName()
	err = checkStringLen(ctx, actionName, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate ActionName [%s]: %+v", actionName, err)
		return err
	}

	triggerStatus := req.GetTriggerStatus()
	err = checkStringLen(ctx, triggerStatus, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate TriggerStatus [%s]: %+v", triggerStatus, err)
		return err
	}

	triggerAction := req.GetTriggerAction()
	err = checkStringLen(ctx, triggerAction, 255)
	if err != nil {
		logger.Error(ctx, "Failed to validate TriggerAction [%s]: %+v", triggerAction, err)
		return err
	}

	policyId := req.GetPolicyId()
	err = checkStringLen(ctx, policyId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate PolicyId [%s]: %+v", policyId, err)
		return err
	}

	nfAddressListId := req.GetNfAddressListId()
	err = checkStringLen(ctx, nfAddressListId, 50)
	if err != nil {
		logger.Error(ctx, "Failed to validate NfAddressListId [%s]: %+v", nfAddressListId, err)
		return err
	}

	return nil
}
