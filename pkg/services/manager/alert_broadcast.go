// Copyright 2018 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coreos/etcd/clientv3"

	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
)

type AlertBroadcast struct {
}

// BroadcastInfo is the service register information to etcd
type BroadcastInfo struct {
	AlertId   string
	Operation string
}

func NewAlertBroadcast() *AlertBroadcast {
	return &AlertBroadcast{}
}

func formatTopic(prefix string, alertId string) string {
	return fmt.Sprintf("%s/%s", prefix, alertId)
}

func (ab *AlertBroadcast) Broadcast(alertId string, operation string, expireTime int64) error {
	info := &BroadcastInfo{
		AlertId:   alertId,
		Operation: operation,
	}

	key := formatTopic("al-broadcast", alertId)
	value, err := json.Marshal(info)
	if err != nil {
		logger.Error(nil, "Marshal BroadcastInfo [%+v] to json failed", info)
		return err
	}

	ctx := context.Background()
	e := global.GetInstance().GetEtcd()

	resp, err := e.Grant(ctx, expireTime)
	if err != nil {
		logger.Error(ctx, "Grant TTL from etcd failed: %+v", err)
		return err
	}

	_, err = e.Put(ctx, key, string(value), clientv3.WithLease(resp.ID))

	if err != nil {
		logger.Error(ctx, "AlertBroadcast [%s] [%s] to etcd failed: %+v", key, string(value), err)
		return err
	}

	return err
}
