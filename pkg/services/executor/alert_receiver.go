package executor

import (
	"errors"
	"strconv"
	"time"

	lib "openpitrix.io/libqueue"
	q "openpitrix.io/libqueue/queue"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/constants"
	"kubesphere.io/alert/pkg/logger"
)

type AlertReceiver struct {
	alertQueue      lib.Topic
	runningAlertIds chan string
	executor        *Executor
}

func NewAlertReceiver() *AlertReceiver {
	cfg := config.GetInstance()
	queueConnStr := cfg.Queue.Addr
	queueType := cfg.Queue.Type

	queueConfigMap := map[string]interface{}{
		"connStr": queueConnStr,
	}

	var alertQueue lib.Topic

	c, err := q.New(queueType, queueConfigMap)
	if err != nil {
		logger.Error(nil, "Failed to connect redis queue: %+v.", err)
	}

	alertQueue, _ = c.SetTopic(constants.AlertTopicPrefix)

	return &AlertReceiver{
		alertQueue:      alertQueue,
		runningAlertIds: make(chan string, 1000),
	}
}

func (ar *AlertReceiver) SetExecutor(executor *Executor) {
	ar.executor = executor
}

func (ar *AlertReceiver) Serve() {
	go ar.ExtractAlerts()

	for i := 0; i < constants.MaxWorkingAlerts; i++ {
		go ar.HandleAlert(strconv.Itoa(i))
	}
}

func (ar *AlertReceiver) ExtractAlerts() error {
	if ar.alertQueue == nil {
		return errors.New("AlertQueue not initialized")
	}

	for {
		alertId, err := ar.alertQueue.Dequeue()
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
