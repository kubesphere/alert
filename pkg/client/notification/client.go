package notification

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"

	"kubesphere.io/alert/pkg/client/notification/pb"
	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/logger"
)

var nfClient *grpc.ClientConn

func getNotificationConn(svcAddress string) (*grpc.ClientConn, error) {
	if nfClient != nil {
		return nfClient, nil
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	keepAlive := keepalive.ClientParameters{
		30 * time.Second,
		10 * time.Second,
		true,
	}

	var err error

	nfClient, err = grpc.DialContext(ctx, svcAddress, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepAlive))

	if err != nil {
		return nil, err
	}

	return nfClient, nil
}

func CheckTimeAvailable(availableStartTimeStr string, availableEndTimeStr string) bool {
	timeFmt := "15:04:05"
	currentTime := time.Now().Format(timeFmt)
	currentTime1, _ := time.Parse(timeFmt, currentTime)

	availableStartTime, _ := time.Parse(timeFmt, availableStartTimeStr)
	availableEndTime, _ := time.Parse(timeFmt, availableEndTimeStr)
	return availableStartTime.Before(currentTime1) && availableEndTime.After(currentTime1)
}

func SendNotification(method string, receiver string, title string, content string) (bool, string) {
	cfg := config.GetInstance()
	conn, err := getNotificationConn(cfg.App.NotificationHost)
	if err != nil {
		logger.Error(nil, "SendNotification getNotificationConn failed %v", err)
		return false, ""
	}

	// sleep a few millsecond for grpc dial etcd
	time.Sleep(time.Millisecond * 500)
	clientX := pb.NewNotificationClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	resp, err := clientX.CreateNotification(ctx, &pb.CreateNotificationRequest{ContentType: &wrappers.StringValue{Value: method}, Title: &wrappers.StringValue{Value: title}, Content: &wrappers.StringValue{Value: content}, ExpiredDays: &wrappers.UInt32Value{Value: 0}, Owner: &wrappers.StringValue{Value: "KubeSphere"}, AddressInfo: &wrappers.StringValue{Value: receiver}})
	if err != nil {
		logger.Error(nil, "SendNotification CreateNotification failed %v", err)
		return false, ""
	}

	return true, resp.GetNotificationId().GetValue()
}

func GetNotificationStatus(notificationIds []string) map[string][]string {
	cfg := config.GetInstance()
	conn, err := getNotificationConn(cfg.App.NotificationHost)
	if err != nil {
		logger.Error(nil, "GetNotificationStatus getNotificationConn failed %v", err)
		return nil
	}

	// sleep a few millsecond for grpc dial etcd
	time.Sleep(time.Millisecond * 500)
	clientX := pb.NewNotificationClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	resp, err := clientX.DescribeTasks(ctx, &pb.DescribeTasksRequest{NotificationId: notificationIds})
	if err != nil {
		logger.Error(nil, "GetNotificationStatus DescribeTasks failed %v", err)
		return nil
	}

	notificationStatusMap := make(map[string][]string)

	for _, task := range resp.TaskSet {
		if notificationStatusMap[task.NotificationId.GetValue()] == nil {
			notificationStatusMap[task.NotificationId.GetValue()] = []string{}
		}
		notificationStatusMap[task.NotificationId.GetValue()] = append(notificationStatusMap[task.NotificationId.GetValue()], task.Directive.GetValue())
		notificationStatusMap[task.NotificationId.GetValue()] = append(notificationStatusMap[task.NotificationId.GetValue()], task.Status.GetValue())
		notificationStatusMap[task.NotificationId.GetValue()] = append(notificationStatusMap[task.NotificationId.GetValue()], fmt.Sprintf("%d", task.StatusTime.GetSeconds()))
	}

	return notificationStatusMap
}
