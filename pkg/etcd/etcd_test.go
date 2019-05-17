// Copyright 2018 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build etcd

package etcd

import (
	"log"
	"testing"

	pkg "kubesphere.io/alert/pkg"
)

func TestConnect(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("LocalDevEnv disabled")
	}
	//e := new(Etcd)
	endpoints := []string{"192.168.0.7:2379"}
	//endpoints:=[]string{"192.168.0.3:2379"}
	prefix := "test"
	e, err := Connect(endpoints, prefix)
	log.Println(e)
	if err != nil {
		t.Fatal(err)
	}

}
