package executor

import (
	"fmt"
	"strconv"
	"time"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/constants"
	"kubesphere.io/alert/pkg/etcd"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
)

type AlertReceiver struct {
	alertQueue      []*etcd.Queue
	runningAlertIds chan string
	executor        *Executor
}

func getQueueNum() int {
	cfg := config.GetInstance()
	queueNum, err := strconv.ParseInt(cfg.Etcd.QueueNum, 10, 0)
	if err != nil {
		queueNum = 1000
	}

	if queueNum < 10 {
		queueNum = 10
	}

	if queueNum > 5000 {
		queueNum = 5000
	}

	return int(queueNum)
}

func NewAlertReceiver() *AlertReceiver {
	queueNum := getQueueNum()

	alertQueue := make([]*etcd.Queue, 0)

	for i := 0; i < queueNum; i++ {
		alertQueue = append(alertQueue, global.GetInstance().GetEtcd().NewQueue(fmt.Sprintf("%s-%d", constants.AlertTopicPrefix, i)))
	}

	return &AlertReceiver{
		alertQueue:      alertQueue,
		runningAlertIds: make(chan string, 1000),
	}
}

func (ar *AlertReceiver) SetExecutor(executor *Executor) {
	ar.executor = executor
}

func (ar *AlertReceiver) Serve() {
	for i := 0; i < len(ar.alertQueue); i++ {
		go ar.ExtractAlerts(i)
	}

	for i := 0; i < constants.MaxWorkingAlerts; i++ {
		go ar.HandleAlert(strconv.Itoa(i))
	}
}

func (ar *AlertReceiver) ExtractAlerts(index int) error {
	for {
		alertId, err := ar.alertQueue[index].Dequeue()
		if err != nil {
			logger.Error(nil, "AlertReceiver failed to dequeue alert from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		logger.Debug(nil, "AlertReceiver dequeue alert [%s] from etcd queue succeed", alertId)
		ar.runningAlertIds <- alertId
	}
}

func (ar *AlertReceiver) HandleAlert(handlerNum string) {
	for {
		alertId := <-ar.runningAlertIds
		logger.Debug(nil, "AlertReceiver handle alert [%s] from etcd queue", alertId)
		ar.executor.AddAlert(alertId)
	}
}
