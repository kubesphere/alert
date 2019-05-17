// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"context"

	"kubesphere.io/alert/pkg/manager"
	"kubesphere.io/alert/pkg/models"
	"kubesphere.io/alert/pkg/pb"
)

func (s *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateResourceTypeRequest:
		return manager.NewChecker(ctx, r).
			Required(models.RtColName, models.RtColParam).
			Exec()
	case *pb.ModifyResourceTypeRequest:
		return manager.NewChecker(ctx, r).
			Required(models.RtColId).
			Exec()
	case *pb.CreateResourceFilterRequest:
		return manager.NewChecker(ctx, r).
			Required(models.RfColParam).
			Exec()
	case *pb.ModifyResourceFilterRequest:
		return manager.NewChecker(ctx, r).
			Required(models.RfColId).
			Exec()
	case *pb.CreateMetricRequest:
		return manager.NewChecker(ctx, r).
			Required(models.MtColName, models.MtColParam, models.MtColTypeId).
			Exec()
	case *pb.ModifyMetricRequest:
		return manager.NewChecker(ctx, r).
			Required(models.MtColId).
			Exec()
	case *pb.CreatePolicyRequest:
		return manager.NewChecker(ctx, r).
			Required(models.PlColConfig, models.PlColCreator, models.PlColAvailableStartTime, models.PlColAvailableEndTime, models.PlColTypeId).
			Exec()
	case *pb.ModifyPolicyRequest:
		return manager.NewChecker(ctx, r).
			Required(models.MtColId).
			Exec()
	case *pb.CreateRuleRequest:
		return manager.NewChecker(ctx, r).
			Required(models.RlColSeverity, models.RlColConditionType, models.RlColThresholds, models.RlColPolicyId, models.RlColMetricId).
			Exec()
	case *pb.ModifyRuleRequest:
		return manager.NewChecker(ctx, r).
			Required(models.RlColId).
			Exec()
	case *pb.CreateAlertRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AlColName, models.AlColPolicyId, models.AlColRsFilterId).
			Exec()
	case *pb.ModifyAlertRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AlColId).
			Exec()
	case *pb.CreateCommentRequest:
		return manager.NewChecker(ctx, r).
			Required(models.CmColAddresser, models.CmColContent, models.CmColHistoryId).
			Exec()
	case *pb.CreateActionRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AcColPolicyId, models.AcColNfAddressListId).
			Exec()
	case *pb.ModifyActionRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AcColId).
			Exec()
	case *pb.ModifyPolicyByAlertRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AlColName).
			Exec()
	case *pb.ModifyAlertByNameRequest:
		return manager.NewChecker(ctx, r).
			Required(models.AlColName, models.AlColPolicyId, models.AlColRsFilterId).
			Exec()
	}

	return nil
}
