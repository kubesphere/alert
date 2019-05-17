// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package watcher

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/constants"
	"kubesphere.io/alert/pkg/etcd"
	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
)

type AlertQueue struct {
	alertQueue []*etcd.Queue
	QueueNum   int
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

func NewAlertQueue() *AlertQueue {
	queueNum := getQueueNum()

	alertQueue := make([]*etcd.Queue, 0)

	for i := 0; i < queueNum; i++ {
		alertQueue = append(alertQueue, global.GetInstance().GetEtcd().NewQueue(fmt.Sprintf("%s-%d", constants.AlertTopicPrefix, i)))
	}

	return &AlertQueue{
		alertQueue: alertQueue,
		QueueNum:   queueNum,
	}
}

func (aq *AlertQueue) WriteBackAlert(alertId string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	aq.createAlert(ctx, alertId)
}

func (aq *AlertQueue) createAlert(ctx context.Context, alertId string) bool {
	// Enqueue alert.
	err := aq.alertQueue[rand.Intn(aq.QueueNum)].Enqueue(alertId)
	if err != nil {
		logger.Error(ctx, "AlertQueue push alert[%s] into etcd failed, [%+v].", alertId, err)
		return false
	}
	logger.Debug(ctx, "AlertQueue push alert[%s] into etcd successfully.", alertId)
	logger.Debug(ctx, "AlertQueue create Alert[%s] successfully.", alertId)

	return true
}
