// Copyright 2018 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package manager

import (
	"strconv"

	"google.golang.org/grpc"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/manager"
	"kubesphere.io/alert/pkg/pb"
)

type Server struct {
	alertQueue     *AlertQueue
	alertBroadcast *AlertBroadcast
}

func Serve() {
	alertQueue := NewAlertQueue()
	alertBroadcast := NewAlertBroadcast()

	s := &Server{
		alertQueue:     alertQueue,
		alertBroadcast: alertBroadcast,
	}

	cfg := config.GetInstance()

	managerHost := cfg.App.Host
	managerPort, _ := strconv.Atoi(cfg.App.Port)

	go ServeApiGateway()

	manager.NewGrpcServer(managerHost, managerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		Serve(func(server *grpc.Server) {
			pb.RegisterAlertManagerServer(server, s)
		})
}
