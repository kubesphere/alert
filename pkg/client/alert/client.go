// Copyright 2018 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file

package alert

import (
	"strconv"

	"kubesphere.io/alert/pkg/config"
	"kubesphere.io/alert/pkg/manager"
	"kubesphere.io/alert/pkg/pb"
)

type Client struct {
	pb.AlertManagerClient
}

func NewClient() (*Client, error) {
	cfg := config.GetInstance()
	managerHost := cfg.App.Host
	managerPort, _ := strconv.Atoi(cfg.App.Port)

	conn, err := manager.NewClient(managerHost, managerPort)
	if err != nil {
		return nil, err
	}
	return &Client{
		AlertManagerClient: pb.NewAlertManagerClient(conn),
	}, nil
}
