package manager

import (
	"context"
	"encoding/json"

	"kubesphere.io/alert/pkg/gerr"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	. "kubesphere.io/alert/pkg/pb"
	rs "kubesphere.io/alert/pkg/services/manager/resource_control"
	"kubesphere.io/alert/pkg/util/stringutil"
)

func (s *Server) operateAlertByName(ctx context.Context, resourceMap map[string]string, alertName string, operation string, req *ModifyAlertByNameRequest) (string, error) {
	alertId, alertName, err := rs.UpdateAlertByName(ctx, resourceMap, alertName, operation, req)
	if err != nil {
		logger.Error(ctx, "Manager update Alert ByName[%s] in DB error, [%+v].", alertId, err)
		return "", gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, alertId)
	}

	// Broadcast alert operation after update DB.
	err = s.alertBroadcast.Broadcast(alertId, operation, 10)
	if err != nil {
		logger.Error(ctx, "Manager broadast alert %s[%s] into etcd failed, [%+v].", operation, alertId, err)
		return "", err
	}
	logger.Debug(ctx, "Manager broadast alert %s[%s] into etcd successfully.", operation, alertId)

	return alertName, nil
}

//0.Policy
//********************************************************************************************************
func (s *Server) ModifyPolicyByAlert(ctx context.Context, req *ModifyPolicyByAlertRequest) (*ModifyPolicyByAlertResponse, error) {
	err := ValidateModifyPolicyByAlertParams(ctx, req)
	if err != nil {
		return nil, err
	}

	resourceMap := map[string]string{}
	err = json.Unmarshal([]byte(req.ResourceSearch), &resourceMap)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Policy By Alert, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	alertName, err := rs.ModifyPolicyByAlert(ctx, resourceMap, req)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Policy By Alert[%s], [%+v].", alertName, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed, alertName)
	}
	logger.Debug(ctx, "Modify Policy By Alert[%s] successfully.", alertName)
	return &ModifyPolicyByAlertResponse{
		AlertName: alertName,
	}, nil
}

//1.Alert
//********************************************************************************************************
func (s *Server) ModifyAlertByName(ctx context.Context, req *ModifyAlertByNameRequest) (*ModifyAlertByNameResponse, error) {
	err := ValidateModifyAlertByNameParams(ctx, req)
	if err != nil {
		return nil, err
	}

	resourceMap := map[string]string{}
	err = json.Unmarshal([]byte(req.ResourceSearch), &resourceMap)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Alert By Name, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	alertName, err := s.operateAlertByName(ctx, resourceMap, req.AlertName, "updating", req)
	if err != nil {
		logger.Error(ctx, "Failed to Modify Alert By Name, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	logger.Debug(ctx, "Modify Alert By Name [%s] successfully.", alertName)
	return &ModifyAlertByNameResponse{
		AlertName: alertName,
	}, nil
}

func (s *Server) DeleteAlertsByName(ctx context.Context, req *DeleteAlertsByNameRequest) (*DeleteAlertsByNameResponse, error) {
	resourceMap := map[string]string{}
	err := json.Unmarshal([]byte(req.ResourceSearch), &resourceMap)
	if err != nil {
		logger.Error(ctx, "Failed to Delete Alerts By Name, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	alertNames := stringutil.SimplifyStringList(req.AlertName)
	alertNamesSuccess := []string{}
	for _, alertName := range alertNames {
		alertNameSuccess, err := s.operateAlertByName(ctx, resourceMap, alertName, "deleting", nil)
		if err != nil {
			logger.Error(ctx, "Failed to Delete Alerts By Name[%s], [%+v].", alertName, err)
		} else {
			alertNamesSuccess = append(alertNamesSuccess, alertNameSuccess)
		}
	}

	logger.Debug(ctx, "Delete Alerts By Name [%+v] successfully.", alertNamesSuccess)
	return &DeleteAlertsByNameResponse{
		AlertName: alertNamesSuccess,
	}, nil
}

func (s *Server) DescribeAlertDetails(ctx context.Context, req *DescribeAlertDetailsRequest) (*DescribeAlertDetailsResponse, error) {
	alds, aldCnt, err := rs.DescribeAlertDetails(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Alert Details, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	aldPbSet := models.ParseAldSet2PbSet(alds)
	res := &DescribeAlertDetailsResponse{
		Total:          uint32(aldCnt),
		AlertdetailSet: aldPbSet,
	}

	logger.Debug(ctx, "Describe Alert Details successfully, Alert Details=[%+v].", res)
	return res, nil
}

func (s *Server) DescribeAlertStatus(ctx context.Context, req *DescribeAlertStatusRequest) (*DescribeAlertStatusResponse, error) {
	alss, alsCnt, err := rs.DescribeAlertStatus(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Alert Status, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	alsPbSet := models.ParseAlsSet2PbSet(alss)
	res := &DescribeAlertStatusResponse{
		Total:          uint32(alsCnt),
		AlertstatusSet: alsPbSet,
	}

	logger.Debug(ctx, "Describe Alert Status successfully, Alert Status=[%+v].", res)
	return res, nil
}

//2.History
//********************************************************************************************************
func (s *Server) DescribeHistoryDetail(ctx context.Context, req *DescribeHistoryDetailRequest) (*DescribeHistoryDetailResponse, error) {
	resourceMap := map[string]string{}
	err := json.Unmarshal([]byte(req.ResourceSearch), &resourceMap)
	if err != nil {
		logger.Error(ctx, "Failed to Describe History Detail, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorUpdateResourceFailed)
	}

	hsds, hsdCnt, err := rs.DescribeHistoryDetail(ctx, resourceMap, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe History Detail, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	hsdPbSet := models.ParseHsdSet2PbSet(hsds)
	res := &DescribeHistoryDetailResponse{
		Total:            uint32(hsdCnt),
		HistorydetailSet: hsdPbSet,
	}

	logger.Debug(ctx, "Describe History Detail successfully, Histories=[%+v].", res)
	return res, nil
}
