// Copyright 2018 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcd

import (
	"context"
	"time"

	"kubesphere.io/alert/pkg/logger"
)

type callback func() error

func (etcd *Etcd) Dlock(ctx context.Context, key string, cb callback) error {
	logger.Debug(ctx, "Create dlock with key [%s]", key)
	mutex, err := etcd.NewMutex(key)
	if err != nil {
		logger.Critical(ctx, "Dlock lock error, failed to create mutex: %+v", err)
		return err
	}
	err = mutex.Lock(ctx)
	if err != nil {
		logger.Critical(ctx, "Dlock lock error, failed to lock mutex: %+v", err)
		return err
	}
	defer mutex.Unlock(ctx)
	err = cb()
	return err
}

func (etcd *Etcd) DlockWithTimeout(key string, timeout time.Duration, cb callback) error {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return etcd.Dlock(ctxWithTimeout, key, cb)
}
