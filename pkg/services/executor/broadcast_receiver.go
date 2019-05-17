package executor

import (
	"context"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"

	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	"kubesphere.io/alert/pkg/util/jsonutil"
)

// BroadcastInfo is the service register information to etcd
type BroadcastInfo struct {
	AlertId   string
	Operation string
}

type BroadcastReceiver struct {
	executor *Executor
}

func NewBroadcastReceiver() *BroadcastReceiver {
	master := &BroadcastReceiver{}

	return master
}

func (br *BroadcastReceiver) SetExecutor(executor *Executor) {
	br.executor = executor
}

func (br *BroadcastReceiver) WatchBroadcast() {
	e := global.GetInstance().GetEtcd()
	watchRes := e.Watch(context.Background(), "al-broadcast/", clientv3.WithPrefix())

	for res := range watchRes {
		for _, ev := range res.Events {
			if ev.Type == mvccpb.PUT {
				var info BroadcastInfo
				err := jsonutil.Decode(ev.Kv.Value, &info)
				if err != nil {
					logger.Error(nil, "WatchBroadcast decode event [%s] [%s] failed: %+v", string(ev.Kv.Key), string(ev.Kv.Value), err)
				} else {
					logger.Info(nil, "WatchBroadcast got put event [%s] [%s]", string(ev.Kv.Key), string(ev.Kv.Value))
					br.executor.ProcessAlert(info.AlertId, info.Operation)
				}
			} else if ev.Type == mvccpb.DELETE {
				logger.Info(nil, "WatchBroadcast got delete event [%s] [%s]", string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}
}
