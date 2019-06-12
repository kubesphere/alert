// Copyright 2018 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package watcher

import (
	"context"
	"strings"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"

	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
	rs "kubesphere.io/alert/pkg/services/watcher/resource_control"
	"kubesphere.io/alert/pkg/util/jsonutil"
)

type Executor struct {
	Name      string
	TaskCount int
}

// ExecutorInfo is the service register information to etcd
type ExecutorInfo struct {
	Name      string
	TaskCount int
}

type ExecutorWatcher struct {
	sync.RWMutex
	members       map[string]*Member
	alertQueue    *AlertQueue
	healthChecker *HealthChecker
}

// Member is a client machine
type Member struct {
	Name      string
	TaskCount int
}

func NewExecutorWatcher() *ExecutorWatcher {
	ew := &ExecutorWatcher{
		members:       make(map[string]*Member),
		alertQueue:    NewAlertQueue(),
		healthChecker: NewHealthChecker(),
	}

	ew.healthChecker.SetExecutorWatcher(ew)

	return ew
}

func (ew *ExecutorWatcher) addExecutor(info *ExecutorInfo) {
	member := &Member{
		Name:      info.Name,
		TaskCount: info.TaskCount,
	}
	ew.Lock()
	ew.members[member.Name] = member
	ew.Unlock()
}

func (ew *ExecutorWatcher) updateExecutor(info *ExecutorInfo) {
	member := ew.members[info.Name]
	ew.Lock()
	member.TaskCount = info.TaskCount
	ew.Unlock()
}

func (ew *ExecutorWatcher) deleteExecutor(name string) {
	_, ok := ew.members[name]
	if ok {
		ew.Lock()
		delete(ew.members, name)
		ew.Unlock()
		ew.migrateAlerts(name)
	}
}

func (ew *ExecutorWatcher) GetExecutors() []ExecutorInfo {
	var executors []ExecutorInfo

	ew.Lock()
	for name := range ew.members {
		member := ew.members[name]
		executors = append(executors, ExecutorInfo{Name: member.Name, TaskCount: member.TaskCount})
	}
	ew.Unlock()

	return executors
}

func (ew *ExecutorWatcher) GetExecutorCount() int {
	return len(ew.members)
}

func (ew *ExecutorWatcher) initExecutors() error {
	ctx := context.Background()
	e := global.GetInstance().GetEtcd()
	key := "alert-executors/"
	resp, err := e.Get(ctx, key, clientv3.WithPrefix())

	if err != nil {
		logger.Error(nil, "ExecutorWatcher Get executors error %+v", err)
		return err
	}

	for _, kv := range resp.Kvs {
		var info ExecutorInfo
		err := jsonutil.Decode(kv.Value, &info)
		if err == nil {
			if _, ok := ew.members[info.Name]; ok {
				logger.Info(nil, "ExecutorWatcher Init update executor %+v", info)
				ew.updateExecutor(&info)
			} else {
				logger.Info(nil, "ExecutorWatcher Init add executor %+v", info)
				ew.addExecutor(&info)
			}
		} else {
			logger.Error(nil, "ExecutorWatcher Decode ExecutorInfo error %+v", err)
		}
	}

	return nil
}

func (ew *ExecutorWatcher) watchExecutors() {
	e := global.GetInstance().GetEtcd()
	watchRes := e.Watch(context.Background(), "alert-executors/", clientv3.WithPrefix())

	for res := range watchRes {
		for _, ev := range res.Events {
			if ev.Type == mvccpb.PUT {
				var info ExecutorInfo
				err := jsonutil.Decode(ev.Kv.Value, &info)
				if err != nil {
					logger.Error(nil, "watchExecutors decode event [%s] [%s] failed: %+v", string(ev.Kv.Key), string(ev.Kv.Value), err)
				} else {
					if _, ok := ew.members[info.Name]; ok {
						logger.Info(nil, "ExecutorWatcher update executor [%s] [%s]", string(ev.Kv.Key), string(ev.Kv.Value))
						ew.updateExecutor(&info)
					} else {
						logger.Info(nil, "ExecutorWatcher add executor [%s] [%s]", string(ev.Kv.Key), string(ev.Kv.Value))
						ew.addExecutor(&info)
					}
				}
			} else if ev.Type == mvccpb.DELETE {
				name := strings.TrimPrefix(string(ev.Kv.Key), "alert-executors/")
				logger.Info(nil, "ExecutorWatcher delete executor [%s]", name)
				ew.deleteExecutor(name)
			}
		}
	}
}

func (ew *ExecutorWatcher) migrateAlerts(executorId string) {
	alertsToMigrate := rs.GetMigrationAlerts(executorId)
	rs.UpdateAlertByExecutorId(executorId, "running", "migrating")

	for _, alert := range alertsToMigrate {
		ew.alertQueue.WriteBackAlert(alert.AlertId)
	}
}

func (ew *ExecutorWatcher) Serve() {
	ew.initExecutors()

	go ew.healthChecker.HealthCheck()
	ew.watchExecutors()
}
