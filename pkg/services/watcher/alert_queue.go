// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package watcher

import (
	"context"
	"time"

	lib "openpitrix.io/libqueue"
	q "openpitrix.io/libqueue/queue"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/constants"
	"kubesphere.io/alert/pkg/logger"
)

type AlertQueue struct {
	alertQueue lib.Topic
}

func NewAlertQueue() *AlertQueue {
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

	return &AlertQueue{
		alertQueue: alertQueue,
	}
}

func (aq *AlertQueue) WriteBackAlert(alertId string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	aq.createAlert(ctx, alertId)
}

func (aq *AlertQueue) createAlert(ctx context.Context, alertId string) bool {
	// Enqueue alert.
	if aq.alertQueue == nil {
		logger.Error(ctx, "AlertQueue push alert[%s] into queue failed. AlertQueue not initialized", alertId)
		return false
	}
	err := aq.alertQueue.Enqueue(alertId)
	if err != nil {
		logger.Error(ctx, "AlertQueue push alert[%s] into queue failed, [%+v].", alertId, err)
		return false
	}
	logger.Debug(ctx, "AlertQueue push alert[%s] into queue successfully.", alertId)
	logger.Debug(ctx, "AlertQueue create Alert[%s] successfully.", alertId)

	return true
}
