// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"errors"

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

func (aq *AlertQueue) Enqueue(alertId string) error {
	if aq.alertQueue == nil {
		return errors.New("AlertQueue not initialized")
	}

	return aq.alertQueue.Enqueue(alertId)
}
