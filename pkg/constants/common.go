// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

const (
	DefaultSelectLimit = 200
)

const (
	DefaultOffset = uint32(0)
	DefaultLimit  = uint32(20)
)

const (
	EtcdPrefix       = "alert/"
	AlertTopicPrefix = "al-job"
	MaxWorkingAlerts = 5
)

const (
	DESC = "desc"
	ASC  = "asc"
)

const (
	StatusAdding    = "adding"
	StatusUpdating  = "updating"
	StatusRunning   = "running"
	StatusDeleting  = "deleting"
	StatusMigrating = "migrating"
)

const (
	ServiceName = "Alert"
)

const MIME_MERGEPATCH = "application/merge-patch+json"
