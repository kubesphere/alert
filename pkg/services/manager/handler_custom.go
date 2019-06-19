package manager

import (
	"context"

	"kubesphere.io/alert/pkg/gerr"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/models"
	. "kubesphere.io/alert/pkg/pb"
	rs "kubesphere.io/alert/pkg/services/manager/resource_control"
)

//0.Alert
//********************************************************************************************************
func (s *Server) DescribeAlertsWithResource(ctx context.Context, req *DescribeAlertsWithResourceRequest) (*DescribeAlertsWithResourceResponse, error) {
	als, alCnt, err := rs.DescribeAlertsWithResource(ctx, req)
	if err != nil {
		logger.Error(ctx, "Failed to Describe Alerts With Resource, [%+v], [%+v].", req, err)
		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	alPbSet := models.ParseAlSet2PbSet(als)
	res := &DescribeAlertsWithResourceResponse{
		Total:    uint32(alCnt),
		AlertSet: alPbSet,
	}

	logger.Debug(ctx, "Describe Alerts With Resource successfully, Alerts=[%+v].", res)
	return res, nil
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

//1.History
//********************************************************************************************************
func (s *Server) DescribeHistoryDetail(ctx context.Context, req *DescribeHistoryDetailRequest) (*DescribeHistoryDetailResponse, error) {
	hsds, hsdCnt, err := rs.DescribeHistoryDetail(ctx, req)
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
