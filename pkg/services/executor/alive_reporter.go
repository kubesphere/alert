package executor

import (
	"context"
	"encoding/json"
	"os"
	"syscall"
	"time"

	"go.etcd.io/etcd/clientv3"

	"kubesphere.io/alert/pkg/global"
	"kubesphere.io/alert/pkg/logger"
)

// ExecutorInfo is the service register information to etcd
type ExecutorInfo struct {
	Name      string
	TaskCount int
}

type AliveReporter struct {
	executor *Executor
}

func NewAliveReporter() *AliveReporter {
	ar := &AliveReporter{}

	return ar
}

func (ar *AliveReporter) SetExecutor(executor *Executor) {
	ar.executor = executor
}

func (ar *AliveReporter) CheckExist() bool {
	ctx := context.Background()
	e := global.GetInstance().GetEtcd()
	key := "alert-executors/" + ar.executor.GetName()
	resp, err := e.Get(ctx, key)
	if err != nil {
		return false
	}

	if resp.Count == 0 {
		return false
	}

	return true
}

func (ar *AliveReporter) putKey(expireTime int64) error {
	ctx := context.Background()
	e := global.GetInstance().GetEtcd()

	key := "alert-executors/" + ar.executor.GetName()

	info := &ExecutorInfo{
		Name:      ar.executor.GetName(),
		TaskCount: ar.executor.GetTaskCount(),
	}

	value, _ := json.Marshal(info)

	resp, err := e.Grant(ctx, expireTime)
	if err != nil {
		logger.Error(nil, "Grant TTL from etcd failed: %+v", err)
		return err
	}

	_, err = e.Put(ctx, key, string(value), clientv3.WithLease(resp.ID))

	if err != nil {
		logger.Error(nil, "AliveReporter putKey [%s] [%s] to etcd failed: %+v", key, string(value), err)
		return err
	}

	return err
}

func (ar *AliveReporter) HeartBoot() {
	err := ar.putKey(30)

	if err != nil {
		logger.Error(nil, "AliveReporter HeartBoot failed: %+v", err)
	}
}

func (ar *AliveReporter) doHeartBeat() {
	if !ar.CheckExist() {
		logger.Info(nil, "AliveReporter detected kickout in etcd")
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
	}

	err := ar.putKey(30)
	if err != nil {
		logger.Error(nil, "AliveReporter HeartBeat error: %+v", err)
	}
}

func (ar *AliveReporter) HeartBeat() {
	timer := time.NewTicker(time.Second * 20)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			ar.doHeartBeat()
		}
	}
}
